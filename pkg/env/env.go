package env

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fusor/cpma/pkg/api"
	"github.com/fusor/cpma/pkg/env/clusterdiscovery"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	// AppName holds the name of this application
	AppName = "CPMA"
	logFile = "cpma.log"
)

var (
	// ConfigFile - keeps full path to the configuration file
	ConfigFile string
	// Login ssh login
	Login string
	// PrivateKey private key path
	PrivateKey string
	// Port ssh port
	Port string

	viperConfig *viper.Viper
)

func init() {
	viperConfig = viper.New()
}

// Config returns pointer to the viper config
func Config() *viper.Viper {
	return viperConfig
}

// InitConfig initializes application's configuration
func InitConfig() error {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		return errors.Wrap(err, "Can't detect home user directory")
	}
	viperConfig.Set("home", home)

	// Try to find config file if it wasn't provided as a flag
	if ConfigFile != "" {
		viperConfig.SetConfigFile(ConfigFile)
	} else {
		viperConfig.AddConfigPath(".")
		viperConfig.AddConfigPath(home)
		viperConfig.SetConfigName("cpma")
	}
	// Fill in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	readConfigErr := viperConfig.ReadInConfig()

	// read all nested values like ssh login, port and key
	getNestedArgValues()

	// Parse kubeconfig for creating api client later
	err = api.ParseKubeConfig()
	if err != nil {
		return errors.Wrap(err, "kubeconfig parsing failed")
	}

	// Ask for all values that are missing in flags or config yaml
	err = surveyMissingValues()
	if err != nil {
		return handleInterrupt(err)
	}

	// If no config was provided, ask to create one for future use
	if readConfigErr != nil {
		err = surveyCreateConfigFile()
		if err != nil {
			return handleInterrupt(err)
		}
		logrus.Debug("Can't read config file, all values were prompted and new config was asked to be created, err: ", readConfigErr)
	}

	return nil
}

func getNestedArgValues() {
	sshCreds := viperConfig.GetStringMapString("SSHCreds")
	if Login != "" {
		sshCreds["login"] = Login
	}

	if PrivateKey != "" {
		sshCreds["privatekey"] = PrivateKey
	}

	if Port != "" {
		sshCreds["port"] = Port
	}
	viperConfig.Set("SSHCreds", sshCreds)
}

func surveyMissingValues() error {
	err := surveyConfigSource()
	if err != nil {
		return err
	}

	switch viperConfig.GetString("ConfigSource") {
	case "remote":
		err := surveySSHConfigValues()
		if err != nil {
			return err
		}
		viperConfig.Set("FetchFromRemote", true)
	case "local":
		err := surveyConfigPaths()
		if err != nil {
			return err
		}
	default:
		return errors.New("Accepted values for config-source are: remote or local")
	}

	err = createAPIClients()
	if err != nil {
		return err
	}

	if viperConfig.GetString("OutputDir") == "" {
		outputDir := "."
		prompt := &survey.Input{
			Message: "Path to output, skip to use current directory",
			Default: ".",
		}
		err := survey.AskOne(prompt, &outputDir, nil)
		if err != nil {
			return err
		}

		viperConfig.Set("OutputDir", outputDir)
	}

	return nil
}

func surveyConfigSource() error {
	// Ask if source of config file should be a remote host or local
	configSource := ""
	if viperConfig.GetString("ConfigSource") == "" {
		prompt := &survey.Select{
			Message: "What will be the source for OCP3 config files?",
			Options: []string{"Remote host", "Local"},
		}
		err := survey.AskOne(prompt, &configSource, nil)
		if err != nil {
			return err
		}
		switch configSource {
		case "Remote host":
			viperConfig.Set("ConfigSource", "remote")
		case "Local":
			viperConfig.Set("ConfigSource", "local")
		}

	}
	return nil
}

func surveySSHConfigValues() error {
	if viperConfig.GetString("Hostname") == "" {
		discoverCluster := ""
		hostname := ""
		clusterName := ""
		var err error

		// Ask for source of master hostname, prompt or find it using KUBECONFIG
		prompt := &survey.Select{
			Message: "Do wish to find source cluster using KUBECONFIG or prompt it?",
			Options: []string{"KUBECONFIG", "prompt"},
		}
		err = survey.AskOne(prompt, &discoverCluster, nil)
		if err != nil {
			return err
		}

		if discoverCluster == "KUBECONFIG" {
			if hostname, clusterName, err = clusterdiscovery.DiscoverCluster(); err != nil {
				return err
			}
			// set cluster name in viper for dumping this value in reusable yaml config
			viperConfig.Set("ClusterName", clusterName)
		} else {
			prompt := &survey.Input{
				Message: "OCP3 Cluster hostname",
			}
			err = survey.AskOne(prompt, &hostname, survey.ComposeValidators(survey.Required))
			if err != nil {
				return err
			}
		}

		viperConfig.Set("Hostname", hostname)
	}

	sshCreds := viperConfig.GetStringMapString("SSHCreds")
	if sshCreds["login"] == "" {
		login := ""
		prompt := &survey.Input{
			Message: "SSH login",
			Default: "root",
		}
		err := survey.AskOne(prompt, &login, nil)
		if err != nil {
			return err
		}

		sshCreds["login"] = login
	}

	if sshCreds["port"] == "" {
		port := ""
		prompt := &survey.Input{
			Message: "SSH Port",
			Default: "22",
		}
		err := survey.AskOne(prompt, &port, nil)
		if err != nil {
			return err
		}

		sshCreds["port"] = port
	}

	if sshCreds["privatekey"] == "" {
		privatekey := ""
		prompt := &survey.Input{
			Message: "Path to private SSH key",
		}
		err := survey.AskOne(prompt, &privatekey, survey.ComposeValidators(survey.Required))
		if err != nil {
			return err
		}

		sshCreds["privatekey"] = privatekey
	}

	viperConfig.Set("SSHCreds", sshCreds)

	// set defaults for remote config files paths
	viperConfig.SetDefault("CrioConfigFile", "/etc/crio/crio.conf")
	viperConfig.SetDefault("ETCDConfigFile", "/etc/etcd/etcd.conf")
	viperConfig.SetDefault("MasterConfigFile", "/etc/origin/master/master-config.yaml")
	viperConfig.SetDefault("NodeConfigFile", "/etc/origin/node/node-config.yaml")
	viperConfig.SetDefault("RegistriesConfigFile", "/etc/containers/registries.conf")
	return nil
}

func surveyConfigPaths() error {
	config := ""
	if viperConfig.GetString("CrioConfigFile") == "" && !viperConfig.InConfig("crioconfigfile") {
		prompt := &survey.Input{
			Message: "Path to crio config file, example: /path/crio/crio.conf",
		}
		err := survey.AskOne(prompt, &config, nil)
		if err != nil {
			return err
		}
		viperConfig.Set("CrioConfigFile", config)
	}

	if viperConfig.GetString("ETCDConfigFile") == "" && !viperConfig.InConfig("etcdconfigfile") {
		prompt := &survey.Input{
			Message: "Path to etcd config file, example: /path/etcd/etcd.conf",
		}
		err := survey.AskOne(prompt, &config, nil)
		if err != nil {
			return err
		}
		viperConfig.Set("ETCDConfigFile", config)
	}

	if viperConfig.GetString("MasterConfigFile") == "" && !viperConfig.InConfig("masterconfigfile") {
		prompt := &survey.Input{
			Message: "Path to master config file, example: /path/etcd/master-config.yaml",
		}
		err := survey.AskOne(prompt, &config, nil)
		if err != nil {
			return err
		}
		viperConfig.Set("MasterConfigFile", config)
	}

	if viperConfig.GetString("NodeConfigFile") == "" && !viperConfig.InConfig("nodeconfigfile") {
		prompt := &survey.Input{
			Message: "Path to node config file, example: /path/node/node-config.yaml",
		}
		err := survey.AskOne(prompt, &config, nil)
		if err != nil {
			return err
		}
		viperConfig.Set("NodeConfigFile", config)
	}

	if viperConfig.GetString("RegistriesConfigFile") == "" && !viperConfig.InConfig("registriesconfigfile") {
		prompt := &survey.Input{
			Message: "Path to registries config file, example: /path/containers/registries.conf",
		}
		err := survey.AskOne(prompt, &config, nil)
		if err != nil {
			return err
		}
		viperConfig.Set("RegistriesConfigFile", config)
	}

	return nil
}

func createAPIClients() error {
	if api.O7tClient != nil && api.K8sClient != nil {
		return nil
	}

	// Ask for cluster name if not provided, can be either prompter or read from current context
	if viperConfig.GetString("ClusterName") == "" {
		contextSource := ""
		prompt := &survey.Select{
			Message: "What will be the source for cluster name used to connect to API?",
			Options: []string{"Current kubeconfig context", "Prompt"},
		}
		err := survey.AskOne(prompt, &contextSource, nil)
		if err != nil {
			return err
		}

		clusterName := ""
		if contextSource == "Prompt" {
			prompt := &survey.Input{
				Message: "Cluster name",
			}
			err := survey.AskOne(prompt, &clusterName, nil)
			if err != nil {
				return err
			}
			// set current context to cluster name for connecting to cluster using client-go
			api.KubeConfig.CurrentContext = api.ClusterNames[clusterName]
		} else {
			// get cluster name from current context for future use
			for key, value := range api.ClusterNames {
				if value == api.KubeConfig.CurrentContext {
					clusterName = key
				}
			}
		}

		viperConfig.Set("ClusterName", clusterName)
	}

	err := api.CreateK8sClient(viperConfig.GetString("ClusterName"))
	if err != nil {
		return errors.Wrap(err, "k8s api client failed to create")
	}

	err = api.CreateO7tClient(viperConfig.GetString("ClusterName"))
	if err != nil {
		return errors.Wrap(err, "OpenShift api client failed to create")
	}

	return nil
}

func surveyCreateConfigFile() error {
	createConfig := ""
	prompt := &survey.Select{
		Message: "No config file found, do you wish to create one for future use?",
		Options: []string{"yes", "no"},
	}
	err := survey.AskOne(prompt, &createConfig, nil)
	if err != nil {
		return err
	}

	if createConfig == "yes" {
		viperConfig.SetConfigFile("cpma.yaml")
		viperConfig.WriteConfig()
	}

	return nil
}

func handleInterrupt(err error) error {
	switch {
	case err.Error() == "interrupt":
		return errors.Wrap(err, "Exiting.")
	default:
		return errors.Wrap(err, "Error in creating config file")
	}
}

// InitLogger initializes stderr and logger to file
func InitLogger() {
	logLevel := logrus.InfoLevel
	if viperConfig.GetBool("debug") {
		logLevel = logrus.DebugLevel
		logrus.SetReportCaller(true)
	}
	logrus.SetLevel(logLevel)

	logrus.SetOutput(ioutil.Discard)

	fileHook, _ := NewLogFileHook(
		LogFileConfig{
			Filename: logFile,
			MaxSize:  5, // MiB
			Level:    logLevel,
			Formatter: &logrus.JSONFormatter{
				PrettyPrint: true,
			},
		},
	)
	logrus.AddHook(fileHook)

	consoleFormatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
		ForceColors:     true,
	}

	if viperConfig.GetBool("verbose") {
		stdoutHook := &ConsoleWriterHook{
			Writer: os.Stdout,
			LogLevels: []logrus.Level{
				logrus.InfoLevel,
				logrus.DebugLevel,
				logrus.WarnLevel,
			},
			Formatter: consoleFormatter,
		}

		logrus.AddHook(stdoutHook)
	}

	stderrHook := &ConsoleWriterHook{
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		},
		Formatter: consoleFormatter,
	}

	logrus.AddHook(stderrHook)

	logrus.Debugf("%s is running in debug mode", AppName)
}

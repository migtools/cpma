package env

import (
	"io/ioutil"
	"os"
	"path"
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

	// OnlyReportMode holds name of cpma mode
	OnlyReportMode = "report"
	// OnlyManifestsMode holds name of cpma mode
	OnlyManifestsMode = "manifests"
)

var (
	// ConfigFile - keeps full path to the configuration file
	ConfigFile string

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
func InitConfig() (err error) {
	// Fill in environment variables that match
	viperConfig.SetEnvPrefix("CPMA")
	viperConfig.AutomaticEnv()

	if err = setConfigLocation(); err != nil {
		return err
	}

	// If a config file is found, read it in.
	readConfigErr := viperConfig.ReadInConfig()

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

// setConfigLocation sets location for CPMA configuration
func setConfigLocation() (err error) {
	var home string
	// Find home directory.
	home, err = homedir.Dir()
	if err != nil {
		return errors.Wrap(err, "Can't detect home user directory")
	}
	viperConfig.Set("home", home)

	// Try to find config file if it wasn't provided as a flag
	if ConfigFile == "" {
		ConfigFile = path.Join(home, "cpma.yaml")
	}
	viperConfig.SetConfigFile(ConfigFile)
	return
}

func surveyMissingValues() error {
	err := surveyCPMAMode()
	if err != nil {
		return err
	}

	err = surveyConfigSource()
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

	if viperConfig.GetString("WorkDir") == "" {
		workDir := "."
		prompt := &survey.Input{
			Message: "Path to application data, skip to use current directory",
			Default: ".",
		}
		err := survey.AskOne(prompt, &workDir, nil)
		if err != nil {
			return err
		}

		viperConfig.Set("WorkDir", workDir)
	}

	return nil
}

func surveyCPMAMode() error {
	mode := viperConfig.GetString("Mode")
	if !viperConfig.InConfig("mode") && mode == "" {
		prompt := &survey.Select{
			Message: "Should CPMA generate only report, only manifests or both?",
			Options: []string{"Both", "Reports only", "Manifests only"},
		}
		err := survey.AskOne(prompt, &mode, nil)
		if err != nil {
			return err
		}

		switch mode {
		case "Reports only":
			viperConfig.Set("Mode", OnlyReportMode)
		case "Manifests only":
			viperConfig.Set("Mode", OnlyManifestsMode)
		}
	}

	switch viperConfig.GetString("Mode") {
	case "report":
	case "manifests":
	case "":
	default:
		return errors.New("Accepted values for mode are: manifests or report")
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
	hostname := viperConfig.GetString("Hostname")
	if !viperConfig.InConfig("hostname") && hostname == "" {
		discoverCluster := ""
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

	clusterName := viperConfig.GetString("ClusterName")
	if !viperConfig.InConfig("clustername") && clusterName == "" {
		prompt := &survey.Input{
			Message: "Cluster name",
		}
		err := survey.AskOne(prompt, &clusterName, nil)
		if err != nil {
			return err
		}

		viperConfig.Set("ClusterName", clusterName)
	}

	login := viperConfig.GetString("SSHLogin")
	if !viperConfig.InConfig("sshlogin") && login == "" {
		prompt := &survey.Input{
			Message: "SSH login",
			Default: "root",
		}
		err := survey.AskOne(prompt, &login, nil)
		if err != nil {
			return err
		}

		viperConfig.Set("SSHLogin", login)
	}

	port := viperConfig.GetString("SSHPort")
	if !viperConfig.InConfig("sshport") && port == "" {
		prompt := &survey.Input{
			Message: "SSH Port",
			Default: "22",
		}
		err := survey.AskOne(prompt, &port, nil)
		if err != nil {
			return err
		}

		viperConfig.Set("SSHPort", port)
	}

	privatekey := viperConfig.GetString("SSHPrivateKey")
	if !viperConfig.InConfig("sshprivatekey") && privatekey == "" {
		prompt := &survey.Input{
			Message: "Path to private SSH key",
		}
		err := survey.AskOne(prompt, &privatekey, survey.ComposeValidators(survey.Required))
		if err != nil {
			return err
		}

		viperConfig.Set("SSHPrivateKey", privatekey)
	}

	// set defaults for remote config files paths
	viperConfig.SetDefault("CrioConfigFile", "/etc/crio/crio.conf")
	viperConfig.SetDefault("ETCDConfigFile", "/etc/etcd/etcd.conf")
	viperConfig.SetDefault("MasterConfigFile", "/etc/origin/master/master-config.yaml")
	viperConfig.SetDefault("NodeConfigFile", "/etc/origin/node/node-config.yaml")
	viperConfig.SetDefault("RegistriesConfigFile", "/etc/containers/registries.conf")
	return nil
}

func surveyConfigPaths() error {
	var config string
	config = viperConfig.GetString("CrioConfigFile")
	if !viperConfig.InConfig("crioconfigfile") && config == "" {
		prompt := &survey.Input{
			Message: "Path to crio config file, example: /path/crio/crio.conf",
		}
		err := survey.AskOne(prompt, &config, nil)
		if err != nil {
			return err
		}
		viperConfig.Set("CrioConfigFile", config)
	}

	config = viperConfig.GetString("ETCDConfigFile")
	if !viperConfig.InConfig("etcdconfigfile") && config == "" {
		prompt := &survey.Input{
			Message: "Path to etcd config file, example: /path/etcd/etcd.conf",
		}
		err := survey.AskOne(prompt, &config, nil)
		if err != nil {
			return err
		}
		viperConfig.Set("ETCDConfigFile", config)
	}

	config = viperConfig.GetString("MasterConfigFile")
	if !viperConfig.InConfig("masterconfigfile") && config == "" {
		prompt := &survey.Input{
			Message: "Path to master config file, example: /path/etcd/master-config.yaml",
		}
		err := survey.AskOne(prompt, &config, nil)
		if err != nil {
			return err
		}
		viperConfig.Set("MasterConfigFile", config)
	}

	config = viperConfig.GetString("NodeConfigFile")
	if !viperConfig.InConfig("nodeconfigfile") && config == "" {
		prompt := &survey.Input{
			Message: "Path to node config file, example: /path/node/node-config.yaml",
		}
		err := survey.AskOne(prompt, &config, nil)
		if err != nil {
			return err
		}
		viperConfig.Set("NodeConfigFile", config)
	}

	config = viperConfig.GetString("RegistriesConfigFile")
	if !viperConfig.InConfig("registriesconfigfile") && config == "" {
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

func surveyCreateConfigFile() (err error) {
	createConfig := viperConfig.GetString("CreateConfig")
	if createConfig == "" {
		prompt := &survey.Select{
			Message: "No config file found, do you wish to create one for future use?",
			Options: []string{"yes", "no"},
		}
		err = survey.AskOne(prompt, &createConfig, nil)
		if err != nil {
			return err
		}
		viperConfig.Set("CreateConfig", createConfig)

	}

	if createConfig == "yes" {
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
			logrus.WarnLevel,
		},
		Formatter: consoleFormatter,
	}

	logrus.AddHook(stderrHook)

	logrus.Debugf("%s is running in debug mode", AppName)
}

package env

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
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

	if err := setConfigLocation(); err != nil {
		return err
	}

	// If a config file is found, read it in.
	readConfigErr := viperConfig.ReadInConfig()
	// If no config file and save config file is undetermined, ask to create or save it for future use
	if readConfigErr != nil && viperConfig.GetString("SaveConfig") != "false" {
		if err := surveySaveConfig(); err != nil {
			return handleInterrupt(err)
		}
		logrus.Debug("Can't read config file, all values were prompted and new config was asked to be created, err: ", readConfigErr)
	}

	// Parse kubeconfig for creating api client later
	if err := api.ParseKubeConfig(); err != nil {
		return errors.Wrap(err, "kubeconfig parsing failed")
	}

	// Ask for all values that are missing in ENV, flags or config yaml
	if err := surveyMissingValues(); err != nil {
		return handleInterrupt(err)
	}

	if viperConfig.GetString("SaveConfig") == "true" {
		viperConfig.WriteConfig()
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
	if err := surveySaveConfig(); err != nil {
		return err
	}

	if err := surveyManifests(); err != nil {
		return err
	}

	if err := surveyReporting(); err != nil {
		return err
	}

	if err := surveyConfigSource(); err != nil {
		return err
	}

	if err := surveyConfigPaths(); err != nil {
		return err
	}

	switch viperConfig.GetString("ConfigSource") {
	case "remote":
		if err := surveySSHConfigValues(); err != nil {
			return err
		}
		viperConfig.Set("FetchFromRemote", true)
	case "local":
		if err := surveyHostname(); err != nil {
			return err
		}
		viperConfig.Set("FetchFromRemote", false)
	default:
		return errors.New("Accepted values for config-source are: remote or local")
	}

	if err := createAPIClients(); err != nil {
		return err
	}

	if viperConfig.GetString("WorkDir") == "" {
		workDir := "."
		prompt := &survey.Input{
			Message: "Path to application data, skip to use current directory",
			Default: ".",
		}
		if err := survey.AskOne(prompt, &workDir); err != nil {
			return err
		}

		viperConfig.Set("WorkDir", workDir)
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
		if err := survey.AskOne(prompt, &configSource); err != nil {
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

func surveyManifests() error {
	manifests := viperConfig.GetString("Manifests")
	if !viperConfig.InConfig("manifests") && manifests == "" {
		prompt := &survey.Select{
			Message: "Would you like to generate manifests?",
			Options: []string{"true", "false"},
		}
		if err := survey.AskOne(prompt, &manifests); err != nil {
			return err
		}
		if manifests == "false" {
			viperConfig.Set("Manifests", false)
		} else {
			viperConfig.Set("Manifests", true)
		}
	}
	return nil
}

func surveyReporting() error {
	reporting := viperConfig.GetString("Reporting")
	if !viperConfig.InConfig("reporting") && reporting == "" {
		prompt := &survey.Select{
			Message: "Would you like reporting?",
			Options: []string{"true", "false"},
		}
		if err := survey.AskOne(prompt, &reporting); err != nil {
			return err
		}
		if reporting == "false" {
			viperConfig.Set("Reporting", false)
		} else {
			viperConfig.Set("Reporting", true)
		}

	}
	return nil
}

func surveyHostname() error {
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
		if err := survey.AskOne(prompt, &discoverCluster); err != nil {
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
			if err := survey.AskOne(prompt, &hostname, survey.WithValidator(survey.Required)); err != nil {
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
		if err := survey.AskOne(prompt, &clusterName); err != nil {
			return err
		}

		viperConfig.Set("ClusterName", clusterName)
	}

	return nil
}

func surveySSHConfigValues() error {
	if err := surveyHostname(); err != nil {
		return err
	}

	login := viperConfig.GetString("SSHLogin")
	if !viperConfig.InConfig("sshlogin") && login == "" {
		prompt := &survey.Input{
			Message: "SSH login",
			Default: "root",
		}
		if err := survey.AskOne(prompt, &login); err != nil {
			return err
		}

		viperConfig.Set("SSHLogin", login)
	}

	port := ""
	if !viperConfig.InConfig("sshport") && viperConfig.GetInt("SSHPort") == 0 {
		prompt := &survey.Input{
			Message: "SSH Port",
			Default: "22",
		}
		if err := survey.AskOne(prompt, &port); err != nil {
			return err
		}

		if p, err := strconv.ParseInt(port, 10, 16); err == nil {
			viperConfig.Set("SSHPort", p)
		}
	}

	privatekey := viperConfig.GetString("SSHPrivateKey")
	if !viperConfig.InConfig("sshprivatekey") && privatekey == "" {
		prompt := &survey.Input{
			Message: "Path to private SSH key",
		}
		if err := survey.AskOne(prompt, &privatekey, survey.WithValidator(survey.Required)); err != nil {
			return err
		}

		viperConfig.Set("SSHPrivateKey", privatekey)
	}

	return nil
}

func surveyConfigPaths() error {
	config := viperConfig.GetString("CrioConfigFile")
	if !viperConfig.InConfig("crioconfigfile") && config == "" {
		prompt := &survey.Input{
			Message: "Path to crio config file",
			Default: "/etc/crio/crio.conf",
		}
		if err := survey.AskOne(prompt, &config); err != nil {
			return err
		}
		viperConfig.Set("CrioConfigFile", config)
	}

	config = viperConfig.GetString("ETCDConfigFile")
	if !viperConfig.InConfig("etcdconfigfile") && config == "" {
		prompt := &survey.Input{
			Message: "Path to etcd config file",
			Default: "/etc/etcd/etcd.conf",
		}
		if err := survey.AskOne(prompt, &config); err != nil {
			return err
		}
		viperConfig.Set("ETCDConfigFile", config)
	}

	config = viperConfig.GetString("MasterConfigFile")
	if !viperConfig.InConfig("masterconfigfile") && config == "" {
		prompt := &survey.Input{
			Message: "Path to master config file",
			Default: "/etc/origin/master/master-config.yaml",
		}
		if err := survey.AskOne(prompt, &config); err != nil {
			return err
		}
		viperConfig.Set("MasterConfigFile", config)
	}

	config = viperConfig.GetString("NodeConfigFile")
	if !viperConfig.InConfig("nodeconfigfile") && config == "" {
		prompt := &survey.Input{
			Message: "Path to node config file",
			Default: "/etc/origin/node/node-config.yaml",
		}
		if err := survey.AskOne(prompt, &config); err != nil {
			return err
		}
		viperConfig.Set("NodeConfigFile", config)
	}

	config = viperConfig.GetString("RegistriesConfigFile")
	if !viperConfig.InConfig("registriesconfigfile") && config == "" {
		prompt := &survey.Input{
			Message: "Path to registries config file",
			Default: "/etc/containers/registries.conf",
		}
		if err := survey.AskOne(prompt, &config); err != nil {
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
			Options: []string{"Current kubeconfig context", "Select kubeconfig context", "Prompt"},
		}
		if err := survey.AskOne(prompt, &contextSource); err != nil {
			return err
		}

		clusterName := ""
		if contextSource == "Prompt" {
			prompt := &survey.Input{
				Message: "Cluster name",
			}
			if err := survey.AskOne(prompt, &clusterName); err != nil {
				return err
			}
			// set current context to cluster name for connecting to cluster using client-go
			api.KubeConfig.CurrentContext = api.ClusterNames[clusterName]
		} else if contextSource == "Current kubeconfig context" {
			// get cluster name from current context for future use
			for key, value := range api.ClusterNames {
				if value == api.KubeConfig.CurrentContext {
					clusterName = key
				}
			}
		} else {
			clusterName = clusterdiscovery.SurveyClusters()
			api.KubeConfig.CurrentContext = api.ClusterNames[clusterName]
		}

		viperConfig.Set("ClusterName", clusterName)
	}

	if err := api.CreateK8sClient(viperConfig.GetString("ClusterName")); err != nil {
		return errors.Wrap(err, "k8s api client failed to create")
	}

	if err := api.CreateO7tClient(viperConfig.GetString("ClusterName")); err != nil {
		return errors.Wrap(err, "OpenShift api client failed to create")
	}

	return nil
}

func surveySaveConfig() (err error) {
	saveConfig := viperConfig.GetString("SaveConfig")
	if saveConfig == "" {
		prompt := &survey.Select{
			Message: "Do you wish to save configuration for future use?",
			Options: []string{"true", "false"},
		}
		if err := survey.AskOne(prompt, &saveConfig); err != nil {
			return err
		}
	}
	if saveConfig == "true" {
		viperConfig.Set("SaveConfig", true)
	} else {
		viperConfig.Set("SaveConfig", false)
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

	if !viperConfig.GetBool("silent") {
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

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

	viperConfig.SetDefault("CrioConfigFile", "/etc/crio/crio.conf")
	viperConfig.SetDefault("ETCDConfigFile", "/etc/etcd/etcd.conf")
	viperConfig.SetDefault("MasterConfigFile", "/etc/origin/master/master-config.yaml")
	viperConfig.SetDefault("NodeConfigFile", "/etc/origin/node/node-config.yaml")
	viperConfig.SetDefault("RegistriesConfigFile", "/etc/containers/registries.conf")

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

	err = api.ParseKubeConfig()
	if err != nil {
		return errors.Wrap(err, "kubeconfig parsing failed")
	}

	getNestedArgValues()

	err = surveyMissingValues()
	if err != nil {
		switch {
		case err.Error() == "interrupt":
			return errors.Wrap(err, "ctrl-C pressed. Exiting.")
		default:
			return errors.Wrap(err, "Error in reading missing values")
		}
	}

	if api.Client == nil {
		err = api.CreateAPIClient(viperConfig.GetString("ClusterName"))
		if err != nil {
			return errors.Wrap(err, "k8s api client failed to create")
		}
	}

	if readConfigErr != nil {
		err = surveyCreateConfigFile()
		if err != nil {
			switch {
			case err.Error() == "interrupt":
				return errors.Wrap(err, "ctrl-C pressed. Exiting.")
			default:
				return errors.Wrap(err, "Error in creating config file")
			}
		}

		logrus.Debug("Can't read config file, all values were prompted and new config was asked to be created, err: ", readConfigErr)
	}

	return nil
}

func surveyMissingValues() error {
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

	if viperConfig.GetString("ClusterName") == "" {
		clusterName := ""
		prompt := &survey.Input{
			Message: "Cluster name",
		}
		err := survey.AskOne(prompt, &clusterName, nil)
		if err != nil {
			return err
		}

		viperConfig.Set("ClusterName", clusterName)
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

	if viperConfig.GetString("OutputDir") == "" {
		outPutDir := "."
		prompt := &survey.Input{
			Message: "Path to output, skip to use current directory",
			Default: ".",
		}
		err := survey.AskOne(prompt, &outPutDir, nil)
		if err != nil {
			return err
		}

		viperConfig.Set("OutputDir", outPutDir)
	}

	viperConfig.Set("SSHCreds", sshCreds)

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

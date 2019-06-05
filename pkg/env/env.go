package env

import (
	"errors"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	// AppName holds the name of this application
	AppName = "CPMA"

	// logFile - keeps full path to the logging file
	logFile = "cpma.log.json" // TODO: we may want this configurable
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
		return errors.New("Can't detect home user directory")
	}
	viperConfig.Set("home", home)

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
	err = viperConfig.ReadInConfig()

	getNestedArgValues()

	cliPromptMissingValues()

	if err != nil {
		cliCreateConfigFile()
		logrus.Debug("Can't read config file, all values were prompted and new config was asked to be created, err: ", err)
	}

	return nil
}

func cliPromptMissingValues() {
	if viperConfig.GetString("Source") == "" {
		hostname := ""
		prompt := &survey.Input{
			Message: "OCP3 Cluster hostname",
		}
		survey.AskOne(prompt, &hostname, nil)
		viperConfig.Set("Source", hostname)
	}

	sshCreds := viperConfig.GetStringMapString("SSHCreds")
	if sshCreds["login"] == "" {
		login := ""
		prompt := &survey.Input{
			Message: "SSH login",
			Default: "root",
		}
		survey.AskOne(prompt, &login, nil)
		sshCreds["login"] = login
	}

	if sshCreds["privatekey"] == "" {
		privatekey := ""
		prompt := &survey.Input{
			Message: "Path to private SSH key",
		}
		survey.AskOne(prompt, &privatekey, nil)
		sshCreds["privatekey"] = privatekey
	}

	if sshCreds["port"] == "" {
		port := ""
		prompt := &survey.Input{
			Message: "SSH Port",
			Default: "22",
		}
		survey.AskOne(prompt, &port, nil)
		sshCreds["port"] = port
	}

	if viperConfig.GetString("OutputDir") == "." {
		outPutDir := ""
		prompt := &survey.Input{
			Message: "Path to output, skip to use current directory",
			Default: ".",
		}
		survey.AskOne(prompt, &outPutDir, nil)
		viperConfig.Set("OutputDir", outPutDir)
	}

	viperConfig.Set("SSHCreds", sshCreds)
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

func cliCreateConfigFile() {
	createConfig := ""
	prompt := &survey.Select{
		Message: "No config file found, do you wish to create one for future use?",
		Options: []string{"yes", "no"},
	}
	survey.AskOne(prompt, &createConfig, nil)

	if createConfig == "yes" {
		viperConfig.SetConfigFile("cpma.yaml")
		viperConfig.WriteConfig()
	}
}

// InitLogger initializes stderr and logger to file
func InitLogger() {
	logLevel := logrus.InfoLevel
	if viperConfig.GetBool("debug") {
		logLevel = logrus.DebugLevel
		logrus.SetReportCaller(true)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	logrus.SetLevel(logLevel)
	logrus.Debugf("%s is running in debug mode", AppName)

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
}

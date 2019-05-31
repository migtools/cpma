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
	if err := viperConfig.ReadInConfig(); err != nil {
		logrus.Debug("Can't read config file, all values will be prompted, err: ", err)
	}

	promptMissingValues()

	return nil
}

func promptMissingValues() {
	if Config().GetString("Source") == "" {
		hostname := ""
		prompt := &survey.Input{
			Message: "OCP3 Cluster hostname",
		}
		survey.AskOne(prompt, &hostname, nil)
		viperConfig.Set("Source", hostname)
	}

	sshCreds := Config().GetStringMapString("SSHCreds")
	if sshCreds["login"] == "" {
		login := ""
		prompt := &survey.Input{
			Message: "SSH login",
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
		}
		survey.AskOne(prompt, &port, nil)
		sshCreds["port"] = port
	}

	if Config().GetString("OutputDir") == "." {
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

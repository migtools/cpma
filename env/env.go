package env

import (
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	viperConfig *viper.Viper

	// ConfigFile - keeps full path to the configuration file
	ConfigFile string
)

func init() {
	viperConfig = viper.New()
}

// Config returns pointer to the viper config
func Config() *viper.Viper {
	return viperConfig
}

// InitConfig initializes config
func InitConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		logrus.Fatal(err)
	}
	viperConfig.Set("home", home)

	if ConfigFile != "" {
		// Use config file from the flag.
		viperConfig.SetConfigFile(ConfigFile)
	} else {
		// Search config in home directory with name ".cpma" (without extension).
		viperConfig.AddConfigPath(home)
		viperConfig.SetConfigName(".cpma")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viperConfig.ReadInConfig(); err != nil {
		logrus.Fatal(err)
	}
}

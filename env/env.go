package env

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var viperConfig *viper.Viper

// ConfigFile Path to config file
var ConfigFile string

// Config returns viper config
func Config() *viper.Viper {
	return viperConfig
}

// InitConfig initialize config
func InitConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
}

func init() {
	viperConfig = viper.New()
}

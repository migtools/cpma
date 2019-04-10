package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var defaultConfig *viper.Viper

// ConfigFile Path to config file
var ConfigFile string

// Config returns config
func Config() *viper.Viper {
	return defaultConfig
}

// InitConfig initialize config
func InitConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defaultConfig.Set("home", home)

	if ConfigFile != "" {
		// Use config file from the flag.
		defaultConfig.SetConfigFile(ConfigFile)
	} else {
		// Search config in home directory with name ".cpma" (without extension).
		defaultConfig.AddConfigPath(home)
		defaultConfig.SetConfigName(".cpma")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := defaultConfig.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	defaultConfig = viper.New()
}

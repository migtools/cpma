package config

import "io/ioutil"

type Parser interface {
	ParseYAML([]byte) error
}

// Load reads configuration file from disk
func Load(configFile string, p Parser) error {
	// Read the config file
	yamlBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	// Parse the config
	if err := p.ParseYAML(yamlBytes); err != nil {
		return err
	}

	return nil
}

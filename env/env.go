package env

import (
	"github.com/fusor/cpma/config"
	yaml "gopkg.in/yaml.v2"
)

// Info contains configuration data
type Info struct {
	Source struct {
		UserName string `yaml:"UserName"`
		HostName string `yaml:"HostName"`
		SSHKey   string `yaml:"SSHKey"`
	} `yaml:"Source"`
}

// ParseYAML unmarshals yaml into configuration structure
func (c *Info) ParseYAML(b []byte) error {
	return yaml.Unmarshal([]byte(b), &c)
}

// LoadConfig reads the yaml configuration file.
func LoadConfig(filename string) (*Info, error) {
	cfg := &Info{}

	err := config.Load(filename, cfg)

	return cfg, err
}

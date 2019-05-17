package config

import (
	"path/filepath"

	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/sirupsen/logrus"
)

// reference:
// https://docs.openshift.com/container-platform/3.11/install_config/master_node_configuration.html

// Config contains CPMA configuration information
type Config struct {
	OutputDir string
	Hostname  string
}

// Fetch files from the OCP3 cluster
func (c *Config) Fetch(path string) ([]byte, error) {
	dst := filepath.Join(c.OutputDir, c.Hostname, path)
	logrus.Infof("Fetching file: %s", dst)
	f, err := io.GetFile(c.Hostname, path, dst)
	if err != nil {
		return nil, err
	}
	logrus.Infof("File:loaded: %v", dst)

	return f, nil
}

// LoadConfig collects and stores configuration for CPMA
func LoadConfig() Config {
	logrus.Info("Loaded config")

	return Config{
		OutputDir: env.Config().GetString("OutputDir"),
		Hostname:  env.Config().GetString("Source"),
	}
}

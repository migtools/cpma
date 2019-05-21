package configmaps

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ConfigMap represent configmap definition
type ConfigMap struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   MetaData `yaml:"metadata"`
	Data       Data     `yaml:"data"`
}

// MetaData configmap's metadata
type MetaData struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

// Data contains CA
type Data struct {
	CAData string `yaml:"ca"`
}

const (
	// APIVersion is the apiVersion string
	APIVersion = "v1"
	// Kind is config map resource type
	Kind = "ConfigMap"
)

// GenConfigMap generates a secret
func GenConfigMap(name string, namespace string, CAData []byte) *ConfigMap {
	return &ConfigMap{
		APIVersion: APIVersion,
		Data: Data{
			CAData: string(CAData),
		},
		Kind: Kind,
		Metadata: MetaData{
			Name:      name,
			Namespace: namespace,
		},
	}
}

// GenYAML returns a YAML of the configMap
func (configMap *ConfigMap) GenYAML() ([]byte, error) {
	yamlBytes, err := yaml.Marshal(&configMap)
	if err != nil {
		logrus.Debugf("Error in config map, config map - %+v", configMap)
		return nil, err
	}
	return yamlBytes, nil
}

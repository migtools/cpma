package secrets

import (
	"log"

	"gopkg.in/yaml.v2"
)

type FileSecret struct {
	HTPasswd string `yaml:"htpasswd"`
}

type LiteralSecret struct {
	ClientSecret string `yaml:"clientSecret`
}

type Secret struct {
	ApiVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Type       string      `yam:"type"`
	MetaData   Metadata    `yaml:"metaData"`
	Data       interface{} `yaml:"data"`
}

type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

var APIVersion = "v1"

func GenSecretFile(name string, encodedSecret string, namespace string) *Secret {
	var secret = Secret{
		ApiVersion: APIVersion,
		Data:       FileSecret{HTPasswd: encodedSecret},
		Kind:       "Secret",
		Type:       "Opaque",
		MetaData: Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
	return &secret
}

func GenSecretLiteral(name string, clientSecret string, namespace string) *Secret {
	var secret = Secret{
		ApiVersion: APIVersion,
		Data:       LiteralSecret{ClientSecret: clientSecret},
		Kind:       "Secret",
		Type:       "Opaque",
		MetaData: Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
	return &secret
}

func (secret *Secret) PrintCRD() string {
	yamlBytes, err := yaml.Marshal(&secret)
	if err != nil {
		log.Fatal(err)
	}
	return string(yamlBytes)
}

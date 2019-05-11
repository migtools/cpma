package secrets

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// TODO: Comment exported functions and structures.
// We may want to unexport some...

type HTPasswdFileSecret struct {
	HTPasswd string `yaml:"htpasswd"`
}

type KeystoneFileSecret struct {
	Keystone string `yaml:"keystone"`
}

type LiteralSecret struct {
	ClientSecret string `yaml:"clientSecret"`
}
type BasicAuthFileSecret struct {
	BasicAuth string `yaml:"basicAuth"`
}

type Secret struct {
	APIVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Type       string      `yaml:"type"`
	Metadata   MetaData    `yaml:"metadata"`
	Data       interface{} `yaml:"data"`
}

type MetaData struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

var APIVersion = "v1"

func GenSecret(name string, secretContent string, namespace string, secretType string) *Secret {
	data := buildData(secretType, secretContent)

	var secret = Secret{
		APIVersion: APIVersion,
		Data:       data,
		Kind:       "Secret",
		Type:       "Opaque",
		Metadata: MetaData{
			Name:      name,
			Namespace: namespace,
		},
	}
	return &secret
}

func buildData(secretType, secretContent string) interface{} {
	var data interface{}

	switch secretType {
	case "keystone":
		data = KeystoneFileSecret{Keystone: secretContent}
	case "htpasswd":
		data = HTPasswdFileSecret{HTPasswd: secretContent}
	case "literal":
		data = LiteralSecret{ClientSecret: secretContent}
	case "basicauth":
		data = BasicAuthFileSecret{BasicAuth: secretContent}
	default:
		logrus.Fatal("Not valid secret type ", secretType)
	}

	return data
}

// GenYAML returns a YAML of the OAuthCRD
func (secret *Secret) GenYAML() []byte {
	yamlBytes, err := yaml.Marshal(&secret)
	if err != nil {
		logrus.WithError(err).Fatal("Cannot generate CRD")
		logrus.Debugf("%+v", secret)
	}
	return yamlBytes
}

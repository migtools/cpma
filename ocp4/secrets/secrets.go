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

type Secret struct {
	APIVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Type       string      `yaml:"type"`
	MetaData   MetaData    `yaml:"metaData"`
	Data       interface{} `yaml:"data"`
}

type MetaData struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

var APIVersion = "v1"

func GenSecretFile(name string, encodedSecret string, namespace string, idenityProviderType string) *Secret {
	data := buildData(idenityProviderType, encodedSecret)

	var secret = Secret{
		APIVersion: APIVersion,
		Data:       data,
		Kind:       "Secret",
		Type:       "Opaque",
		MetaData: MetaData{
			Name:      name,
			Namespace: namespace,
		},
	}
	return &secret
}

func buildData(idenityProviderType, encodedSecret string) interface{} {
	var data interface{}

	switch idenityProviderType {
	case "keystone":
		data = KeystoneFileSecret{Keystone: encodedSecret}
	case "htpasswd":
		data = HTPasswdFileSecret{HTPasswd: encodedSecret}
	default:
		logrus.Fatal("Not valid idenity provider type ", idenityProviderType)
	}

	return data
}

func GenSecretLiteral(name string, clientSecret string, namespace string) *Secret {
	var secret = Secret{
		APIVersion: APIVersion,
		Data:       LiteralSecret{ClientSecret: clientSecret},
		Kind:       "Secret",
		Type:       "Opaque",
		MetaData: MetaData{
			Name:      name,
			Namespace: namespace,
		},
	}
	return &secret
}

// GenYAML returns a YAML of the OAuthCRD
func (secret *Secret) GenYAML() string {
	yamlBytes, err := yaml.Marshal(&secret)
	if err != nil {
		logrus.WithError(err).Fatal("Cannot generate CRD")
		logrus.Debugf("%+v", secret)
	}
	return string(yamlBytes)
}

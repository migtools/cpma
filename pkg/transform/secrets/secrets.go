package secrets

import (
	"errors"
)

// HTPasswdFileSecret is an htpasswd secret
type HTPasswdFileSecret struct {
	HTPasswd string `yaml:"htpasswd"`
}

// KeystoneFileSecret is a keystone secret
type KeystoneFileSecret struct {
	Keystone string `yaml:"keystone"`
}

// LiteralSecret is a literal secret
type LiteralSecret struct {
	ClientSecret string `yaml:"clientSecret"`
}

// BasicAuthFileSecret is a basic auth secret
type BasicAuthFileSecret struct {
	BasicAuth string `yaml:"basicAuth"`
}

// Secret contains a secret
type Secret struct {
	APIVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Type       string      `yaml:"type"`
	Metadata   MetaData    `yaml:"metadata"`
	Data       interface{} `yaml:"data"`
}

// MetaData is the Metadata for a secret
type MetaData struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

// SecretType is an enumerator for secret types
type SecretType int

const (
	// KeystoneSecretType - keystone type for Secret
	KeystoneSecretType = iota
	// HtpasswdSecretType - htpasswd type for Secret
	HtpasswdSecretType
	// LiteralSecretType - literal type for Secret
	LiteralSecretType
	// BasicAuthSecretType - basicauth type for Secret
	BasicAuthSecretType
)

var typeArray = []string{
	"KeystoneSecretType",
	"HtpasswdSecretType",
	"LiteralSecretType",
	"BasicAuthSecretType",
}

// APIVersion is the apiVersion string
var APIVersion = "v1"

// GenSecret generates a secret
func GenSecret(name string, secretContent string, namespace string, secretType SecretType) (*Secret, error) {
	data, err := buildData(secretType, secretContent)
	if err != nil {
		return nil, err
	}

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
	return &secret, nil
}

func buildData(secretType SecretType, secretContent string) (interface{}, error) {
	var data interface{}

	switch secretType {
	case KeystoneSecretType:
		data = KeystoneFileSecret{Keystone: secretContent}
	case HtpasswdSecretType:
		data = HTPasswdFileSecret{HTPasswd: secretContent}
	case LiteralSecretType:
		data = LiteralSecret{ClientSecret: secretContent}
	case BasicAuthSecretType:
		data = BasicAuthFileSecret{BasicAuth: secretContent}
	default:
		return nil, errors.New("Not a valid secret type " + secretType.String())
	}

	return data, nil
}

// SecretType.String returns a string representation for SecretType enum
func (secType SecretType) String() string {
	if secType >= KeystoneSecretType && int(secType) < len(typeArray) {
		return typeArray[secType]
	}
	return "unknown"
}

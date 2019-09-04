package secrets

import (
	"errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation"
)

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

const (
	apiVersion      = "v1"
	secretNameError = `Secret name is no valid, make sure it consists of lower case alphanumeric characters, ‘-’ or ‘.’,` +
		`and must start and end with an alphanumeric character (e.g. ‘example.com’, regex used for validation is ‘[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*’)`
)

// GenTLSSecret generates a TLS secret
func GenTLSSecret(name string, namespace string, cert []byte, key []byte) (*corev1.Secret, error) {
	nameErrors := validation.IsDNS1123Label(name)
	if nameErrors != nil {
		return nil, errors.New(secretNameError)
	}

	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiVersion,
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: "kubernetes.io/tls",
		Data: map[string][]byte{
			"tls.cert": cert,
			"tls.key":  key,
		},
	}

	return secret, nil
}

// GenSecret generates a secret
func GenSecret(name string, secretContent string, namespace string, secretType SecretType) (*corev1.Secret, error) {
	nameErrors := validation.IsDNS1123Label(name)
	if nameErrors != nil {
		return nil, errors.New(secretNameError)
	}

	data, err := buildData(secretType, secretContent)
	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiVersion,
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: "Opaque",
		Data: data,
	}

	return secret, nil
}

func buildData(secretType SecretType, secretContent string) (map[string][]byte, error) {
	var data map[string][]byte

	switch secretType {
	case KeystoneSecretType:
		data = map[string][]byte{
			"keystone": []byte(secretContent),
		}
	case HtpasswdSecretType:
		data = map[string][]byte{
			"htpasswd": []byte(secretContent),
		}
	case LiteralSecretType:
		data = map[string][]byte{
			"clientSecret": []byte(secretContent),
		}
	case BasicAuthSecretType:
		data = map[string][]byte{
			"basicAuth": []byte(secretContent),
		}
	default:
		return nil, errors.New("Unknown secret type")
	}

	return data, nil
}

package secrets

import (
	"errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation"
)

const (
	apiVersion      = "v1"
	secretNameError = `Secret name is no valid, make sure it consists of lower case alphanumeric characters, ‘-’ or ‘.’,` +
		`and must start and end with an alphanumeric character (e.g. ‘example.com’, regex used for validation is ‘[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*’)`
)

// new creates a secret core without Type and Data
func new(name string, namespace string, secretType corev1.SecretType, data map[string][]byte) (*corev1.Secret, error) {
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
		Type: secretType,
		Data: data,
	}

	return secret, nil
}

// TLS generates a TLS secret
func TLS(name string, namespace string, cert []byte, key []byte) (*corev1.Secret, error) {
	data := map[string][]byte{
		"tls.cert": cert,
		"tls.key":  key,
	}

	return new(name, namespace, corev1.SecretTypeTLS, data)
}

// Opaque generates an Opaque secret
func Opaque(name string, secretContent []byte, namespace string, dataName string) (*corev1.Secret, error) {
	if secretContent == nil {
		secretContent = []byte("")
	}

	data := map[string][]byte{
		dataName: secretContent,
	}

	return new(name, namespace, corev1.SecretTypeOpaque, data)
}

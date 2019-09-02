package apicert

import (
	"path/filepath"

	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/secrets"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
)

// Translate ImagePolicyConfig definitions
func Translate(servingInfo legacyconfigv1.ServingInfo) (*corev1.Secret, error) {
	const (
		secretName      = "api-server-cert-secret"
		namespace       = "openshift-config"
		defaultCertPath = "/etc/origin/master"
	)

	if servingInfo.CertFile == "" || servingInfo.KeyFile == "" {
		return nil, errors.New("No Secret available")
	}

	path := defaultCertPath
	dir, file := filepath.Split(servingInfo.CertFile)
	if dir != "" {
		path = dir
	}
	certFile := filepath.Join(path, file)
	certContent, err := io.FetchFile(certFile)
	if err != nil {
		return nil, err
	}

	path = defaultCertPath
	dir, file = filepath.Split(servingInfo.KeyFile)
	if dir != "" {
		path = dir
	}
	keyFile := filepath.Join(path, file)
	keyContent, err := io.FetchFile(keyFile)
	if err != nil {
		return nil, err
	}

	tlsSecret, err := secrets.GenTLSSecret(secretName, namespace, certContent, keyContent)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate TLS secret, see error")
	}

	return tlsSecret, nil
}

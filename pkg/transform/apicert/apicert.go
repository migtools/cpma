package apicert

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"path/filepath"
	"strings"

	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/secrets"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
)

// certSigner gets certificate CN
func certSigner(certContent []byte) string {
	block, _ := pem.Decode(certContent)
	if block == nil || block.Type != "CERTIFICATE" {
		return ""
	}

	certif, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("Can't read certif: %s", err)
	}
	return certif.Issuer.CommonName
}

func getCert(certName, keyName string) ([]byte, []byte) {
	const defaultCertPath = "/etc/origin/master"

	crtFile := setDefaultPath(certName, defaultCertPath)
	keyFile := setDefaultPath(keyName, defaultCertPath)

	crtContent, err := io.FetchFile(crtFile)
	if err != nil {
		return nil, nil
	}

	keyContent, err := io.FetchFile(keyFile)
	if err != nil {
		return nil, nil
	}

	return crtContent, keyContent
}

func getLocalCert(certName string) []byte {
	const defaultCertPath = "/etc/origin/master"

	crtFile := setDefaultPath(certName, defaultCertPath)
	crtContent, err := io.FetchFromLocal(crtFile)
	if err != nil {
		return nil
	}

	return crtContent
}

// setDefaultPath determines a default path for given filename
func setDefaultPath(name string, defaultPath string) string {
	path := defaultPath
	dir, file := filepath.Split(name)
	if dir != "" {
		path = dir
	}
	return filepath.Join(path, file)
}

// OCPSigned tells if locally stored certificate is OCP signed or not
func OCPSigned(crtFile string) bool {
	crtContent := getLocalCert(crtFile)
	if strings.Contains(certSigner(crtContent), "openshift-signer@") {
		return true
	}
	return false
}

// Translate ImagePolicyConfig definitions
func Translate(servingInfo legacyconfigv1.ServingInfo) (*corev1.Secret, error) {
	const (
		namespace  = "openshift-config"
		secretName = "api-server-cert-secret"
	)

	if servingInfo.CertFile == "" || servingInfo.KeyFile == "" {
		return nil, errors.New("No Secret available")
	}

	crtContent, keyContent := getCert(servingInfo.CertFile, servingInfo.KeyFile)
	if strings.Contains(certSigner(crtContent), "openshift-signer@") {
		logrus.Info("APITransform: API certificate is OpenShift signed, not ported")
		return nil, nil
	}

	tlsSecret, err := secrets.GenTLSSecret(secretName, namespace, crtContent, keyContent)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate TLS secret, see error")
	}

	return tlsSecret, nil
}

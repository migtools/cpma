package oauth

import (
	"encoding/base64"
	"path/filepath"

	"github.com/fusor/cpma/env"
	"github.com/fusor/cpma/pkg/transform/secrets"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

// IdentityProviderKeystone is a Keystone specific identity provider
type IdentityProviderKeystone struct {
	identityProviderCommon `yaml:",inline"`
	Keystone               struct {
		DomainName string `yaml:"domainName"`
		URL        string `yaml:"url"`
		CA         struct {
			Name string `yaml:"name"`
		} `yaml:"ca"`
		TLSClientCert struct {
			Name string `yaml:"name"`
		} `yaml:"tlsClientCert"`
		TLSClientKey struct {
			Name string `yaml:"name"`
		} `yaml:"tlsClientKey"`
	} `yaml:"keystone"`
}

func buildKeystoneIP(serializer *json.Serializer, p IdentityProvider) (IdentityProviderKeystone, secrets.Secret, secrets.Secret) {
	var (
		idP        IdentityProviderKeystone
		keystone   configv1.KeystonePasswordIdentityProvider
		certSecret = new(secrets.Secret)
		keySecret  = new(secrets.Secret)
	)

	_, _, _ = serializer.Decode(p.Provider.Raw, nil, &keystone)

	idP.Type = "Keystone"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.Keystone.DomainName = keystone.DomainName
	idP.Keystone.URL = keystone.URL
	idP.Keystone.CA.Name = keystone.CA

	if keystone.UseKeystoneIdentity {
		logrus.Warn("Keystone useKeystoneIdentity value is not supported in OCP4")
	}

	if keystone.CertFile != "" {
		outputDir := env.Config().GetString("OutputDir")
		host := env.Config().GetString("Source")

		certSecretName := p.Name + "-client-cert-secret"
		idP.Keystone.TLSClientCert.Name = certSecretName
		src := filepath.Join(keystone.KeyFile)
		dst := filepath.Join(outputDir, host, keystone.CertFile)
		certFile := GetFile(host, src, dst)
		encoded := base64.StdEncoding.EncodeToString([]byte(certFile))
		certSecret = secrets.GenSecret(certSecretName, encoded, "openshift-config", "keystone")

		keySecretName := p.Name + "-client-key-secret"
		idP.Keystone.TLSClientKey.Name = keySecretName
		src = filepath.Join(keystone.KeyFile)
		dst = filepath.Join(outputDir, host, keystone.KeyFile)
		keyFile := GetFile(host, src, dst)
		encoded = base64.StdEncoding.EncodeToString([]byte(keyFile))
		keySecret = secrets.GenSecret(keySecretName, encoded, "openshift-config", "keystone")
	}

	return idP, *certSecret, *keySecret
}

package oauth

import (
	"encoding/base64"

	"github.com/fusor/cpma/pkg/transform/secrets"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

// IdentityProviderKeystone is a Keystone specific identity provider
type IdentityProviderKeystone struct {
	identityProviderCommon `yaml:",inline"`
	Keystone               Keystone `yaml:"keystone"`
}

// Keystone specific Provider data
type Keystone struct {
	DomainName    string        `yaml:"domainName"`
	URL           string        `yaml:"url"`
	CA            CA            `yaml:"ca"`
	TLSClientCert TLSClientCert `yaml:"tlsClientCert"`
	TLSClientKey  TLSClientKey  `yaml:"tlsClientKey"`
}

// CA secret name
type CA struct {
	Name string `yaml:"name"`
}

// TLSClientCert secret name
type TLSClientCert struct {
	Name string `yaml:"name"`
}

// TLSClientKey secret name
type TLSClientKey struct {
	Name string `yaml:"name"`
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
		certSecretName := p.Name + "-client-cert-secret"
		idP.Keystone.TLSClientCert.Name = certSecretName
		encoded := base64.StdEncoding.EncodeToString(p.CrtData)
		certSecret = secrets.GenSecret(certSecretName, encoded, "openshift-config", "keystone")

		keySecretName := p.Name + "-client-key-secret"
		idP.Keystone.TLSClientKey.Name = keySecretName
		encoded = base64.StdEncoding.EncodeToString(p.KeyData)
		keySecret = secrets.GenSecret(keySecretName, encoded, "openshift-config", "keystone")
	}

	return idP, *certSecret, *keySecret
}

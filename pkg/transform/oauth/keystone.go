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

func buildKeystoneIP(serializer *json.Serializer, p IdentityProvider) (IdentityProviderKeystone, secrets.Secret, secrets.Secret, error) {
	var (
		idP        IdentityProviderKeystone
		keystone   configv1.KeystonePasswordIdentityProvider
		certSecret = new(secrets.Secret)
		keySecret  = new(secrets.Secret)
		err        error
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
		certSecret, err = secrets.GenSecret(certSecretName, encoded, "openshift-config", "keystone")
		if err != nil {
			return idP, *certSecret, *keySecret, nil
		}

		keySecretName := p.Name + "-client-key-secret"
		idP.Keystone.TLSClientKey.Name = keySecretName
		encoded = base64.StdEncoding.EncodeToString(p.KeyData)
		keySecret, err = secrets.GenSecret(keySecretName, encoded, "openshift-config", "keystone")
		if err != nil {
			return idP, *certSecret, *keySecret, nil
		}
	}

	return idP, *certSecret, *keySecret, nil
}

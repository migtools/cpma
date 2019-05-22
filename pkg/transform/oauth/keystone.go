package oauth

import (
	"encoding/base64"

	"github.com/fusor/cpma/pkg/transform/configmaps"
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
	DomainName    string         `yaml:"domainName"`
	URL           string         `yaml:"url"`
	CA            *CA            `yaml:"ca,omitempty"`
	TLSClientCert *TLSClientCert `yaml:"tlsClientCert,omitempty"`
	TLSClientKey  *TLSClientKey  `yaml:"tlsClientKey,omitempty"`
}

func buildKeystoneIP(serializer *json.Serializer, p IdentityProvider) (IdentityProviderKeystone, secrets.Secret, secrets.Secret, *configmaps.ConfigMap, error) {
	var (
		idP         IdentityProviderKeystone
		certSecret  = new(secrets.Secret)
		keySecret   = new(secrets.Secret)
		caConfigmap *configmaps.ConfigMap
		err         error
		keystone    configv1.KeystonePasswordIdentityProvider
	)
	_, _, err = serializer.Decode(p.Provider.Raw, nil, &keystone)
	if err != nil {
		return idP, *certSecret, *keySecret, nil, err
	}

	idP.Type = "Keystone"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.Keystone.DomainName = keystone.DomainName
	idP.Keystone.URL = keystone.URL

	if keystone.CA != "" {
		caConfigmap = configmaps.GenConfigMap("keystone-configmap", OAuthNamespace, p.CAData)
		idP.Keystone.CA = &CA{Name: caConfigmap.Metadata.Name}
	}

	if keystone.UseKeystoneIdentity {
		logrus.Warn("Keystone useKeystoneIdentity value is not supported in OCP4")
	}

	if keystone.CertFile != "" {
		certSecretName := p.Name + "-client-cert-secret"
		idP.Keystone.TLSClientCert = &TLSClientCert{Name: certSecretName}
		encoded := base64.StdEncoding.EncodeToString(p.CrtData)
		certSecret, err = secrets.GenSecret(certSecretName, encoded, OAuthNamespace, secrets.KeystoneSecretType)
		if err != nil {
			return idP, *certSecret, *keySecret, nil, nil
		}

		keySecretName := p.Name + "-client-key-secret"
		idP.Keystone.TLSClientKey = &TLSClientKey{Name: keySecretName}
		encoded = base64.StdEncoding.EncodeToString(p.KeyData)
		keySecret, err = secrets.GenSecret(keySecretName, encoded, OAuthNamespace, secrets.KeystoneSecretType)
		if err != nil {
			return idP, *certSecret, *keySecret, nil, nil
		}
	}

	return idP, *certSecret, *keySecret, caConfigmap, nil
}

package oauth

import (
	"encoding/base64"

	"github.com/fusor/cpma/pkg/transform/configmaps"
	"github.com/fusor/cpma/pkg/transform/secrets"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

// IdentityProviderBasicAuth is a basic auth specific identity provider
type IdentityProviderBasicAuth struct {
	identityProviderCommon `yaml:",inline"`
	BasicAuth              BasicAuth `yaml:"basicAuth"`
}

// BasicAuth provider specific data
// BasicAuth provider specific data
type BasicAuth struct {
	URL           string         `yaml:"url"`
	CA            *CA            `yaml:"ca,omitempty"`
	TLSClientCert *TLSClientCert `yaml:"tlsClientCert,omitempty"`
	TLSClientKey  *TLSClientKey  `yaml:"tlsClientKey,omitempty"`
}

func buildBasicAuthIP(serializer *json.Serializer, p IdentityProvider) (IdentityProviderBasicAuth, secrets.Secret, secrets.Secret, *configmaps.ConfigMap, error) {
	var (
		err         error
		idP         IdentityProviderBasicAuth
		certSecret  = &secrets.Secret{}
		keySecret   = &secrets.Secret{}
		caConfigmap *configmaps.ConfigMap
		basicAuth   configv1.BasicAuthPasswordIdentityProvider
	)

	_, _, err = serializer.Decode(p.Provider.Raw, nil, &basicAuth)
	if err != nil {
		return idP, *certSecret, *keySecret, nil, err
	}

	idP.Type = "BasicAuth"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.BasicAuth.URL = basicAuth.URL

	if basicAuth.CA != "" {
		caConfigmap = configmaps.GenConfigMap("basicauth-configmap", OAuthNamespace, p.CAData)
		idP.BasicAuth.CA = &CA{Name: caConfigmap.Metadata.Name}
	}

	if basicAuth.CertFile != "" {
		certSecretName := p.Name + "-client-cert-secret"
		idP.BasicAuth.TLSClientCert = &TLSClientCert{Name: certSecretName}

		encoded := base64.StdEncoding.EncodeToString(p.CrtData)
		certSecret, err = secrets.GenSecret(certSecretName, encoded, OAuthNamespace, secrets.BasicAuthSecretType)
		if err != nil {
			return idP, *certSecret, *keySecret, nil, err
		}

		keySecretName := p.Name + "-client-key-secret"
		idP.BasicAuth.TLSClientKey = &TLSClientKey{Name: keySecretName}

		encoded = base64.StdEncoding.EncodeToString(p.KeyData)
		keySecret, err = secrets.GenSecret(keySecretName, encoded, OAuthNamespace, secrets.BasicAuthSecretType)
		if err != nil {
			return idP, *certSecret, *keySecret, nil, err
		}
	}

	return idP, *certSecret, *keySecret, caConfigmap, nil
}

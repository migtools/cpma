package oauth

import (
	"encoding/base64"

	"github.com/fusor/cpma/pkg/transform/secrets"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

// IdentityProviderBasicAuth is a basic auth specific identity provider
type IdentityProviderBasicAuth struct {
	identityProviderCommon `yaml:",inline"`
	BasicAuth              struct {
		URL string `yaml:"url"`
		CA  struct {
			Name string `yaml:"name"`
		} `yaml:"ca"`
		TLSClientCert struct {
			Name string `yaml:"name"`
		} `yaml:"tlsClientCert"`
		TLSClientKey struct {
			Name string `yaml:"name"`
		} `yaml:"tlsClientKey"`
	} `yaml:"basicAuth"`
}

func buildBasicAuthIP(serializer *json.Serializer, p IdentityProvider) (IdentityProviderBasicAuth, secrets.Secret, secrets.Secret, error) {
	var idP IdentityProviderBasicAuth
	var certSecret *secrets.Secret
	var keySecret *secrets.Secret

	var basicAuth configv1.BasicAuthPasswordIdentityProvider
	_, _, _ = serializer.Decode(p.Provider.Raw, nil, &basicAuth)

	idP.Type = "BasicAuth"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.BasicAuth.URL = basicAuth.URL
	idP.BasicAuth.CA.Name = basicAuth.CA

	certSecretName := p.Name + "-client-cert-secret"
	idP.BasicAuth.TLSClientCert.Name = certSecretName

	// TODO: Fetch cert and key
	certFile := "456This is pretend content"
	encoded := base64.StdEncoding.EncodeToString([]byte(certFile))
	certSecret, err := secrets.GenSecret(certSecretName, encoded, "openshift-config", "basicauth")
	if err != nil {
		return idP, *certSecret, *keySecret, err
	}

	keySecretName := p.Name + "-client-key-secret"
	idP.BasicAuth.TLSClientKey.Name = keySecretName

	keyFile := "123This is pretend content"
	encoded = base64.StdEncoding.EncodeToString([]byte(keyFile))
	keySecret, err = secrets.GenSecret(keySecretName, encoded, "openshift-config", "basicauth")
	if err != nil {
		return idP, *certSecret, *keySecret, err
	}

	return idP, *certSecret, *keySecret, nil
}

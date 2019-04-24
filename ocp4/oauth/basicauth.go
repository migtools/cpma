package oauth

import (
	"encoding/base64"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	"github.com/fusor/cpma/ocp4/secrets"
	configv1 "github.com/openshift/api/legacyconfig/v1"
)

type identityProviderBasicAuth struct {
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

func buildBasicAuthIP(serializer *json.Serializer, p configv1.IdentityProvider) (identityProviderBasicAuth, secrets.Secret, secrets.Secret) {
	var idP identityProviderBasicAuth
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
	certSecret := secrets.GenSecret(certSecretName, encoded, "openshift-config", "basicauth")

	keySecretName := p.Name + "-client-key-secret"
	idP.BasicAuth.TLSClientKey.Name = keySecretName

	keyFile := "123This is pretend content"
	encoded = base64.StdEncoding.EncodeToString([]byte(keyFile))
	keySecret := secrets.GenSecret(keySecretName, encoded, "openshift-config", "basicauth")

	return idP, *certSecret, *keySecret
}

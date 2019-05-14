package oauth

import (
	"github.com/fusor/cpma/pkg/transform/secrets"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

//IdentityProviderGoogle is a Google specific identity provider
type IdentityProviderGoogle struct {
	identityProviderCommon `yaml:",inline"`
	Google                 struct {
		ClientID     string `yaml:"clientID"`
		ClientSecret struct {
			Name string `yaml:"name"`
		} `yaml:"clientSecret"`
		HostedDomain string `yaml:"hostedDomain"`
	} `yaml:"google"`
}

func buildGoogleIP(serializer *json.Serializer, p IdentityProvider) (IdentityProviderGoogle, secrets.Secret) {
	var idP IdentityProviderGoogle
	var google configv1.GoogleIdentityProvider
	_, _, _ = serializer.Decode(p.Provider.Raw, nil, &google)

	idP.Type = "Google"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.Google.ClientID = google.ClientID
	idP.Google.HostedDomain = google.HostedDomain

	secretName := p.Name + "-secret"
	idP.Google.ClientSecret.Name = secretName
	secret := secrets.GenSecret(secretName, google.ClientSecret.Value, "openshift-config", "literal")

	return idP, *secret
}

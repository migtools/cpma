package oauth

import (
	"github.com/fusor/cpma/pkg/transform/secrets"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

// IdentityProviderOpenID is an Open ID specific identity provider
type IdentityProviderOpenID struct {
	identityProviderCommon `yaml:",inline"`
	OpenID                 struct {
		ClientID     string `yaml:"clientID"`
		ClientSecret struct {
			Name string `yaml:"name"`
		} `yaml:"clientSecret"`
		Claims struct {
			PreferredUsername []string `yaml:"preferredUsername"`
			Name              []string `yaml:"name"`
			Email             []string `yaml:"email"`
		} `yaml:"claims"`
		URLs struct {
			Authorize string `yaml:"authorize"`
			Token     string `yaml:"token"`
		} `yaml:"urls"`
	} `yaml:"openID"`
}

func buildOpenIDIP(serializer *json.Serializer, p IdentityProvider) (IdentityProviderOpenID, secrets.Secret, error) {
	var idP IdentityProviderOpenID
	var openID configv1.OpenIDIdentityProvider
	_, _, _ = serializer.Decode(p.Provider.Raw, nil, &openID)

	idP.Type = "OpenID"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.OpenID.ClientID = openID.ClientID
	idP.OpenID.Claims.PreferredUsername = openID.Claims.PreferredUsername
	idP.OpenID.Claims.Name = openID.Claims.Name
	idP.OpenID.Claims.Email = openID.Claims.Email
	idP.OpenID.URLs.Authorize = openID.URLs.Authorize
	idP.OpenID.URLs.Token = openID.URLs.Token

	secretName := p.Name + "-secret"
	idP.OpenID.ClientSecret.Name = secretName
	secret, err := secrets.GenSecret(secretName, openID.ClientSecret.Value, "openshift-config", "literal")
	if err != nil {
		return idP, *secret, err
	}

	return idP, *secret, nil
}

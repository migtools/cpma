package oauth

import (
	"encoding/base64"
	"errors"

	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/secrets"
	configv1 "github.com/openshift/api/legacyconfig/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

// IdentityProviderOpenID is an Open ID specific identity provider
type IdentityProviderOpenID struct {
	identityProviderCommon `yaml:",inline"`
	OpenID                 OpenID `yaml:"openID"`
}

// OpenID provider specific data
type OpenID struct {
	ClientID     string       `yaml:"clientID"`
	ClientSecret ClientSecret `yaml:"clientSecret"`
	Claims       OpenIDClaims `yaml:"claims"`
	URLs         OpenIDURLs   `yaml:"urls"`
}

// OpenIDClaims are the claims for an OpenID provider
type OpenIDClaims struct {
	PreferredUsername []string `yaml:"preferredUsername"`
	Name              []string `yaml:"name"`
	Email             []string `yaml:"email"`
}

// OpenIDURLs are the URLs for an OpenID provider
type OpenIDURLs struct {
	Authorize string `yaml:"authorize"`
	Token     string `yaml:"token"`
}

func buildOpenIDIP(serializer *json.Serializer, p IdentityProvider) (*IdentityProviderOpenID, *secrets.Secret, error) {
	var (
		err    error
		secret *secrets.Secret
		idP    = &IdentityProviderOpenID{}
		openID configv1.OpenIDIdentityProvider
	)
	_, _, err = serializer.Decode(p.Provider.Raw, nil, &openID)
	if err != nil {
		return nil, nil, err
	}

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
	secretContent, err := io.FetchStringSource(openID.ClientSecret)
	if err != nil {
		return nil, nil, err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(secretContent))
	secret, err = secrets.GenSecret(secretName, encoded, OAuthNamespace, secrets.LiteralSecretType)
	if err != nil {
		return nil, nil, err
	}

	return idP, secret, nil
}

func validateOpenIDProvider(serializer *json.Serializer, p IdentityProvider) error {
	var openID configv1.OpenIDIdentityProvider

	_, _, err := serializer.Decode(p.Provider.Raw, nil, &openID)
	if err != nil {
		return err
	}

	if p.Name == "" {
		return errors.New("Name can't be empty")
	}

	if err := validateMappingMethod(p.MappingMethod); err != nil {
		return err
	}

	if openID.ClientSecret.KeyFile != "" {
		return errors.New("Usage of encrypted files as secret value is not supported")
	}

	if err := validateClientData(openID.ClientID, openID.ClientSecret); err != nil {
		return err
	}

	if len(openID.Claims.ID) == 0 && len(openID.Claims.PreferredUsername) == 0 && len(openID.Claims.Name) == 0 && len(openID.Claims.Email) == 0 {
		return errors.New("All claims are empty. At least one is required")
	}

	if openID.URLs.Authorize == "" {
		return errors.New("Authorization endpoint can't be empty")
	}

	if openID.URLs.Token == "" {
		return errors.New("Token endpoint can't be empty")
	}

	return nil
}

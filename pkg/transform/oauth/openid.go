package oauth

import (
	"encoding/base64"

	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/secrets"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

// IdentityProviderOpenID is an Open ID specific identity provider
type IdentityProviderOpenID struct {
	identityProviderCommon `json:",inline"`
	OpenID                 OpenID `json:"openID"`
}

// OpenID provider specific data
type OpenID struct {
	ClientID     string       `json:"clientID"`
	ClientSecret ClientSecret `json:"clientSecret"`
	Claims       OpenIDClaims `json:"claims"`
	URLs         OpenIDURLs   `json:"urls"`
}

// OpenIDClaims are the claims for an OpenID provider
type OpenIDClaims struct {
	PreferredUsername []string `json:"preferredUsername"`
	Name              []string `json:"name"`
	Email             []string `json:"email"`
}

// OpenIDURLs are the URLs for an OpenID provider
type OpenIDURLs struct {
	Authorize string `json:"authorize"`
	Token     string `json:"token"`
}

func buildOpenIDIP(serializer *json.Serializer, p IdentityProvider) (*IdentityProviderOpenID, *secrets.Secret, error) {
	var (
		err    error
		secret *secrets.Secret
		idP    = &IdentityProviderOpenID{}
		openID legacyconfigv1.OpenIDIdentityProvider
	)

	if _, _, err = serializer.Decode(p.Provider.Raw, nil, &openID); err != nil {
		return nil, nil, errors.Wrap(err, "Something is wrong in decoding openID")
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
		return nil, nil, errors.Wrap(err, "Something is wrong in fetching client secret for openID")
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(secretContent))

	if secret, err = secrets.GenSecret(secretName, encoded, OAuthNamespace, secrets.LiteralSecretType); err != nil {
		return nil, nil, errors.Wrap(err, "Something is wrong in generating secret for openID")
	}

	return idP, secret, nil
}

func validateOpenIDProvider(serializer *json.Serializer, p IdentityProvider) error {
	var openID legacyconfigv1.OpenIDIdentityProvider

	if _, _, err := serializer.Decode(p.Provider.Raw, nil, &openID); err != nil {
		return errors.Wrap(err, "Something is wrong in decoding openID")
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

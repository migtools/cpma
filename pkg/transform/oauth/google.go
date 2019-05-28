package oauth

import (
	"encoding/base64"
	"errors"

	"github.com/fusor/cpma/pkg/config"
	"github.com/fusor/cpma/pkg/transform/secrets"
	configv1 "github.com/openshift/api/legacyconfig/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

//IdentityProviderGoogle is a Google specific identity provider
type IdentityProviderGoogle struct {
	identityProviderCommon `yaml:",inline"`
	Google                 Google `yaml:"google"`
}

// Google provider specific data
type Google struct {
	ClientID     string       `yaml:"clientID"`
	ClientSecret ClientSecret `yaml:"clientSecret"`
	HostedDomain string       `yaml:"hostedDomain,omitempty"`
}

func buildGoogleIP(serializer *json.Serializer, p IdentityProvider, config *config.Config) (*IdentityProviderGoogle, *secrets.Secret, error) {
	var (
		err    error
		idP    = &IdentityProviderGoogle{}
		secret *secrets.Secret
		google configv1.GoogleIdentityProvider
	)
	_, _, err = serializer.Decode(p.Provider.Raw, nil, &google)
	if err != nil {
		return nil, nil, err
	}

	idP.Type = "Google"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.Google.ClientID = google.ClientID
	idP.Google.HostedDomain = google.HostedDomain

	secretName := p.Name + "-secret"
	idP.Google.ClientSecret.Name = secretName
	secretContent, err := fetchStringSource(google.ClientSecret, config)

	encoded := base64.StdEncoding.EncodeToString([]byte(secretContent))
	secret, err = secrets.GenSecret(secretName, encoded, OAuthNamespace, secrets.LiteralSecretType)
	if err != nil {
		return nil, nil, err
	}

	return idP, secret, nil
}

func validateGoogleProvider(serializer *json.Serializer, p IdentityProvider) error {
	var google configv1.GoogleIdentityProvider

	_, _, err := serializer.Decode(p.Provider.Raw, nil, &google)
	if err != nil {
		return err
	}

	if p.Name == "" {
		return errors.New("Name can't be empty")
	}

	if err := validateMappingMethod(p.MappingMethod); err != nil {
		return err
	}

	if err := validateClientData(google.ClientID, google.ClientSecret); err != nil {
		return err
	}

	return nil
}

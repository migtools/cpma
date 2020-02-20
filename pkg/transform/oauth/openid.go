package oauth

import (
	"github.com/konveyor/cpma/pkg/io"
	"github.com/konveyor/cpma/pkg/transform/secrets"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func buildOpenIDIP(serializer *json.Serializer, p IdentityProvider) (*ProviderResources, error) {
	var (
		err             error
		providerSecrets []*corev1.Secret
		idP             = &configv1.IdentityProvider{}
		openID          legacyconfigv1.OpenIDIdentityProvider
	)

	if _, _, err = serializer.Decode(p.Provider.Raw, nil, &openID); err != nil {
		return nil, errors.Wrap(err, "Failed to decode openID, see error")
	}

	idP.Type = "OpenID"
	idP.Name = p.Name
	idP.MappingMethod = configv1.MappingMethodType(p.MappingMethod)
	idP.OpenID = &configv1.OpenIDIdentityProvider{}
	idP.OpenID.ClientID = openID.ClientID
	idP.OpenID.Claims.PreferredUsername = openID.Claims.PreferredUsername
	idP.OpenID.Claims.Name = openID.Claims.Name
	idP.OpenID.Claims.Email = openID.Claims.Email

	secretName := "openid-secret"
	idP.OpenID.ClientSecret.Name = secretName
	secretContent, err := io.FetchStringSource(openID.ClientSecret)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch client secret for openID, see error")
	}

	secret, err := secrets.Opaque(secretName, []byte(secretContent), OAuthNamespace, "clientSecret")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate secret for openID, see error")
	}
	providerSecrets = append(providerSecrets, secret)

	return &ProviderResources{
		IDP:     idP,
		Secrets: providerSecrets,
	}, nil
}

func validateOpenIDProvider(serializer *json.Serializer, p IdentityProvider) error {
	var openID legacyconfigv1.OpenIDIdentityProvider

	if _, _, err := serializer.Decode(p.Provider.Raw, nil, &openID); err != nil {
		return errors.Wrap(err, "Failed to decode openID, see error")
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

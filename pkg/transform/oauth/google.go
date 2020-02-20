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

func buildGoogleIP(serializer *json.Serializer, p IdentityProvider) (*ProviderResources, error) {
	var (
		err             error
		idP             = &configv1.IdentityProvider{}
		providerSecrets []*corev1.Secret
		google          legacyconfigv1.GoogleIdentityProvider
	)

	if _, _, err = serializer.Decode(p.Provider.Raw, nil, &google); err != nil {
		return nil, errors.Wrap(err, "Failed to decode google, see error")
	}

	idP.Type = "Google"
	idP.Name = p.Name
	idP.MappingMethod = configv1.MappingMethodType(p.MappingMethod)
	idP.Google = &configv1.GoogleIdentityProvider{}
	idP.Google.ClientID = google.ClientID
	idP.Google.HostedDomain = google.HostedDomain

	secretName := "google-secret"
	idP.Google.ClientSecret.Name = secretName
	secretContent, err := io.FetchStringSource(google.ClientSecret)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch client secret for google, see error")
	}

	secret, err := secrets.Opaque(secretName, []byte(secretContent), OAuthNamespace, "clientSecret")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate client secret for google, see error")
	}
	providerSecrets = append(providerSecrets, secret)

	return &ProviderResources{
		IDP:     idP,
		Secrets: providerSecrets,
	}, nil
}

func validateGoogleProvider(serializer *json.Serializer, p IdentityProvider) error {
	var google legacyconfigv1.GoogleIdentityProvider

	if _, _, err := serializer.Decode(p.Provider.Raw, nil, &google); err != nil {
		return errors.Wrap(err, "Failed to decode google, see error")
	}

	if p.Name == "" {
		return errors.New("Name can't be empty")
	}

	if err := validateMappingMethod(p.MappingMethod); err != nil {
		return err
	}

	if google.ClientSecret.KeyFile != "" {
		return errors.New("Usage of encrypted files as secret value is not supported")
	}

	if err := validateClientData(google.ClientID, google.ClientSecret); err != nil {
		return err
	}

	return nil
}

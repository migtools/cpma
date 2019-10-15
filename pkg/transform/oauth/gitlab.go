package oauth

import (
	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/configmaps"
	"github.com/fusor/cpma/pkg/transform/secrets"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func buildGitLabIP(serializer *json.Serializer, p IdentityProvider) (*ProviderResources, error) {
	var (
		err                error
		idP                = &configv1.IdentityProvider{}
		providerSecrets    []*corev1.Secret
		providerConfigMaps []*corev1.ConfigMap
		gitlab             legacyconfigv1.GitLabIdentityProvider
	)

	if _, _, err = serializer.Decode(p.Provider.Raw, nil, &gitlab); err != nil {
		return nil, errors.Wrap(err, "Failed to decode gitlab, see error")
	}

	idP.Type = "GitLab"
	idP.Name = p.Name
	idP.MappingMethod = configv1.MappingMethodType(p.MappingMethod)
	idP.GitLab = &configv1.GitLabIdentityProvider{}
	idP.GitLab.URL = gitlab.URL
	idP.GitLab.ClientID = gitlab.ClientID

	if gitlab.CA != "" {
		caConfigmap := configmaps.GenConfigMap("gitlab-configmap", OAuthNamespace, p.CAData)
		idP.GitLab.CA = configv1.ConfigMapNameReference{Name: caConfigmap.ObjectMeta.Name}
		providerConfigMaps = append(providerConfigMaps, caConfigmap)
	}

	secretName := "gitlab-secret"
	idP.GitLab.ClientSecret.Name = secretName
	secretContent, err := io.FetchStringSource(gitlab.ClientSecret)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch client secret for gitlab, see error")
	}

	secret, err := secrets.Opaque(secretName, []byte(secretContent), OAuthNamespace, "clientSecret")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate secret for gitlab, see error")
	}
	providerSecrets = append(providerSecrets, secret)

	return &ProviderResources{
		IDP:        idP,
		Secrets:    providerSecrets,
		ConfigMaps: providerConfigMaps,
	}, nil
}

func validateGitLabProvider(serializer *json.Serializer, p IdentityProvider) error {
	var gitlab legacyconfigv1.GitLabIdentityProvider

	if _, _, err := serializer.Decode(p.Provider.Raw, nil, &gitlab); err != nil {
		return errors.Wrap(err, "Failed to decode gitlab, see error")
	}

	if p.Name == "" {
		return errors.New("Name can't be empty")
	}

	if err := validateMappingMethod(p.MappingMethod); err != nil {
		return err
	}

	if gitlab.URL == "" {
		return errors.New("URL can't be empty")
	}

	if gitlab.ClientSecret.KeyFile != "" {
		return errors.New("Usage of encrypted files as secret value is not supported")
	}

	if err := validateClientData(gitlab.ClientID, gitlab.ClientSecret); err != nil {
		return err
	}

	return nil
}

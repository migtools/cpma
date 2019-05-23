package oauth

import (
	"errors"

	"github.com/fusor/cpma/pkg/transform/configmaps"
	"github.com/fusor/cpma/pkg/transform/secrets"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

// IdentityProviderGitLab is a Gitlab specific identity provider
type IdentityProviderGitLab struct {
	identityProviderCommon `yaml:",inline"`
	GitLab                 GitLab `yaml:"gitlab"`
}

// GitLab provider specific data
type GitLab struct {
	URL          string       `yaml:"url"`
	CA           *CA          `yaml:"ca,omitempty"`
	ClientID     string       `yaml:"clientID"`
	ClientSecret ClientSecret `yaml:"clientSecret"`
}

func buildGitLabIP(serializer *json.Serializer, p IdentityProvider) (*IdentityProviderGitLab, *secrets.Secret, *configmaps.ConfigMap, error) {
	var (
		err         error
		idP         = &IdentityProviderGitLab{}
		secret      *secrets.Secret
		caConfigmap *configmaps.ConfigMap
		gitlab      configv1.GitLabIdentityProvider
	)
	_, _, err = serializer.Decode(p.Provider.Raw, nil, &gitlab)
	if err != nil {
		return nil, nil, nil, err
	}

	idP.Type = "GitLab"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.GitLab.URL = gitlab.URL
	idP.GitLab.ClientID = gitlab.ClientID

	if gitlab.CA != "" {
		caConfigmap = configmaps.GenConfigMap("gitlab-configmap", OAuthNamespace, p.CAData)
		idP.GitLab.CA = &CA{Name: caConfigmap.Metadata.Name}
	}

	secretName := p.Name + "-secret"
	idP.GitLab.ClientSecret.Name = secretName
	secret, err = secrets.GenSecret(secretName, gitlab.ClientSecret.Value, OAuthNamespace, secrets.LiteralSecretType)
	if err != nil {
		return nil, nil, nil, err
	}

	return idP, secret, caConfigmap, nil
}

func validateGitLabProvider(serializer *json.Serializer, p IdentityProvider) error {
	var gitlab configv1.GitLabIdentityProvider

	_, _, err := serializer.Decode(p.Provider.Raw, nil, &gitlab)
	if err != nil {
		return err
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

	if err := validateClientData(gitlab.ClientID, gitlab.ClientSecret); err != nil {
		return err
	}

	return nil
}

package oauth

import (
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
	CA           CA           `yaml:"ca"`
	ClientID     string       `yaml:"clientID"`
	ClientSecret ClientSecret `yaml:"clientSecret"`
}

func buildGitLabIP(serializer *json.Serializer, p IdentityProvider) (IdentityProviderGitLab, secrets.Secret, error) {
	var idP IdentityProviderGitLab
	var gitlab configv1.GitLabIdentityProvider
	_, _, _ = serializer.Decode(p.Provider.Raw, nil, &gitlab)

	idP.Type = "GitLab"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.GitLab.URL = gitlab.URL
	idP.GitLab.CA.Name = gitlab.CA
	idP.GitLab.ClientID = gitlab.ClientID

	secretName := p.Name + "-secret"
	idP.GitLab.ClientSecret.Name = secretName
	secret, err := secrets.GenSecret(secretName, gitlab.ClientSecret.Value, "openshift-config", "literal")
	if err != nil {
		return idP, *secret, err
	}

	return idP, *secret, nil
}

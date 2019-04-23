package oauth

import (
	"github.com/fusor/cpma/ocp4/secrets"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

type identityProviderGitLab struct {
	identityProviderCommon `yaml:",inline"`
	GitLab                 struct {
		URL string `yaml:"url"`
		CA  struct {
			Name string `yaml:"name"`
		} `yaml:"ca"`
		ClientID     string `yaml:"clientID"`
		ClientSecret struct {
			Name string `yaml:"name"`
		} `yaml:"clientSecret"`
	} `yaml:"gitlab"`
}

func buildGitLabIP(serializer *json.Serializer, p configv1.IdentityProvider) (identityProviderGitLab, secrets.Secret) {
	var idP identityProviderGitLab
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
	secret := secrets.GenSecret(secretName, gitlab.ClientSecret.Value, "openshift-config", "literal")

	return idP, *secret
}

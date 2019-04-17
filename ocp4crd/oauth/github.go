package oauth

import (
	configv1 "github.com/openshift/api/legacyconfig/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

type identityProviderGitHub struct {
	Name          string `yaml:"name"`
	Challenge     bool   `yaml:"challenge"`
	Login         bool   `yaml:"login"`
	MappingMethod string `yaml:"mappingMethod"`
	Type          string `yaml:"type"`
	GitHub        struct {
		HostName string `yaml:"hostname"`
		CA       struct {
			Name string `yaml:"name"`
		} `yaml:"ca"`
		ClientID     string `yaml:"clientID"`
		ClientSecret struct {
			Name string `yaml:"name"`
		} `yaml:"clientSecret"`
		Organizations []string `yaml:"organizations"`
		Teams         []string `yaml:"teams"`
	} `yaml:"github"`
}

func buildGitHubIP(serializer *json.Serializer, p configv1.IdentityProvider) identityProviderGitHub {
	var idP identityProviderGitHub
	var github configv1.GitHubIdentityProvider
	_, _, _ = serializer.Decode(p.Provider.Raw, nil, &github)

	idP.Type = "GitHub"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.GitHub.HostName = github.Hostname
	idP.GitHub.CA.Name = github.CA
	idP.GitHub.ClientID = github.ClientID
	idP.GitHub.Organizations = github.Organizations
	idP.GitHub.Teams = github.Teams
	// TODO: Learn how to handle secrets
	idP.GitHub.ClientSecret.Name = github.ClientSecret.Value
	return idP
}

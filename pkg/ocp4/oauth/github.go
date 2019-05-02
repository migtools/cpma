package oauth

import (
	"encoding/base64"

	"github.com/fusor/cpma/pkg/ocp4/secrets"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

type identityProviderGitHub struct {
	identityProviderCommon `yaml:",inline"`
	GitHub                 struct {
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

func buildGitHubIP(serializer *json.Serializer, p configv1.IdentityProvider) (identityProviderGitHub, secrets.Secret) {
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

	secretName := p.Name + "-secret"
	idP.GitHub.ClientSecret.Name = secretName

	encoded := base64.StdEncoding.EncodeToString([]byte(github.ClientSecret.Value))
	secret := secrets.GenSecret(secretName, encoded, "openshift-config", "literal")

	return idP, *secret
}

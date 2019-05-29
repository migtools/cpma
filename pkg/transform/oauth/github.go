package oauth

import (
	"encoding/base64"
	"errors"

	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/configmaps"
	"github.com/fusor/cpma/pkg/transform/secrets"
	configv1 "github.com/openshift/api/legacyconfig/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

//IdentityProviderGitHub is a Github specific identity provider
type IdentityProviderGitHub struct {
	identityProviderCommon `yaml:",inline"`
	GitHub                 GitHub `yaml:"github"`
}

// GitHub provider specific data
type GitHub struct {
	HostName      string       `yaml:"hostname,omitempty"`
	CA            *CA          `yaml:"ca,omitempty"`
	ClientID      string       `yaml:"clientID"`
	ClientSecret  ClientSecret `yaml:"clientSecret"`
	Organizations []string     `yaml:"organizations,omitempty"`
	Teams         []string     `yaml:"teams,omitempty"`
}

func buildGitHubIP(serializer *json.Serializer, p IdentityProvider) (*IdentityProviderGitHub, *secrets.Secret, *configmaps.ConfigMap, error) {
	var (
		err         error
		idP         = &IdentityProviderGitHub{}
		secret      *secrets.Secret
		caConfigmap *configmaps.ConfigMap
		github      configv1.GitHubIdentityProvider
	)

	_, _, err = serializer.Decode(p.Provider.Raw, nil, &github)
	if err != nil {
		return nil, nil, nil, err
	}

	idP.Type = "GitHub"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.GitHub.HostName = github.Hostname
	idP.GitHub.ClientID = github.ClientID
	idP.GitHub.Organizations = github.Organizations
	idP.GitHub.Teams = github.Teams

	if github.CA != "" {
		caConfigmap = configmaps.GenConfigMap("github-configmap", OAuthNamespace, p.CAData)
		idP.GitHub.CA = &CA{Name: caConfigmap.Metadata.Name}
	}

	secretName := p.Name + "-secret"
	idP.GitHub.ClientSecret.Name = secretName
	secretContent, err := io.FetchStringSource(github.ClientSecret)
	if err != nil {
		return nil, nil, nil, err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(secretContent))
	secret, err = secrets.GenSecret(secretName, encoded, OAuthNamespace, secrets.LiteralSecretType)
	if err != nil {
		return nil, nil, nil, err
	}

	return idP, secret, caConfigmap, nil
}

func validateGithubProvider(serializer *json.Serializer, p IdentityProvider) error {
	var github configv1.GitHubIdentityProvider

	_, _, err := serializer.Decode(p.Provider.Raw, nil, &github)
	if err != nil {
		return err
	}

	if p.Name == "" {
		return errors.New("Name can't be empty")
	}

	if err := validateMappingMethod(p.MappingMethod); err != nil {
		return err
	}

	if github.ClientSecret.KeyFile != "" {
		return errors.New("Usage of encrypted files as secret value is not supported")
	}

	if err := validateClientData(github.ClientID, github.ClientSecret); err != nil {
		return err
	}

	return nil
}

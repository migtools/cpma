package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fusor/cpma/ocp3"
)

func TestTranslateMasterConfigGithub(t *testing.T) {
	masterConfig := ocp3.ParseMaster("../../test/oauth/github-test-master-config.yaml")

	var expectedCrd OAuthCRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = "openshift-config"

	var githubIDP identityProviderGitHub
	githubIDP.Type = "GitHub"
	githubIDP.Challenge = false
	githubIDP.Login = true
	githubIDP.MappingMethod = "claim"
	githubIDP.Name = "github123456789"
	githubIDP.GitHub.HostName = "test.example.com"
	githubIDP.GitHub.CA.Name = "github.crt"
	githubIDP.GitHub.ClientID = "2d85ea3f45d6777bffd7"
	githubIDP.GitHub.Organizations = []string{"myorganization1", "myorganization2"}
	githubIDP.GitHub.Teams = []string{"myorganization1/team-a", "myorganization2/team-b"}
	githubIDP.GitHub.ClientSecret.Name = "github123456789-secret"
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, githubIDP)

	resCrd, _, err := Translate(masterConfig.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, &expectedCrd, resCrd)
}

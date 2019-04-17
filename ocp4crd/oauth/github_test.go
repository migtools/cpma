package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocp3 "github.com/fusor/cpma/ocp3config"
)

func TestGenerateMasterConfigGithub(t *testing.T) {
	testConfig := ocp3.Config{
		Masterf: "../../test/oauth/github-test-master-config.yaml",
	}
	masterConfig := testConfig.ParseMaster()

	var expectedCrd v4OAuthCRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.MetaData.Name = "cluster"

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
	githubIDP.GitHub.ClientSecret.Name = "e16a59ad33d7c29fd4354f46059f0950c609a7ea"
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, githubIDP)

	resCrd, err := Generate(masterConfig)
	require.NoError(t, err)
	assert.Equal(t, &expectedCrd, resCrd)
}

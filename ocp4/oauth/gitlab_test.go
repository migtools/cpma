package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fusor/cpma/ocp3"
)

func TestTranslateMasterConfigGitlab(t *testing.T) {
	masterConfig := ocp3.ParseMaster("../../test/oauth/gitlab-test-master-config.yaml")

	var expectedCrd OAuthCRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = "openshift-config"

	var gitlabIDP identityProviderGitLab
	gitlabIDP.Type = "GitLab"
	gitlabIDP.Challenge = true
	gitlabIDP.Login = true
	gitlabIDP.MappingMethod = "claim"
	gitlabIDP.Name = "gitlab123456789"
	gitlabIDP.GitLab.URL = "https://gitlab.com/"
	gitlabIDP.GitLab.CA.Name = "gitlab.crt"
	gitlabIDP.GitLab.ClientID = "fake-id"
	gitlabIDP.GitLab.ClientSecret.Name = "gitlab123456789-secret"
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, gitlabIDP)

	resCrd, _, err := Translate(masterConfig.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, &expectedCrd, resCrd)
}

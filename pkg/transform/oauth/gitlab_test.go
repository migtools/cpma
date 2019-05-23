package oauth_test

import (
	"testing"

	"github.com/fusor/cpma/pkg/transform/oauth"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigGitlab(t *testing.T) {
	identityProviders, err := cpmatest.LoadIPTestData("testdata/gitlab-test-master-config.yaml")
	require.NoError(t, err)

	var expectedCrd oauth.CRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = oauth.OAuthNamespace

	var gitlabIDP = &oauth.IdentityProviderGitLab{}
	gitlabIDP.Type = "GitLab"
	gitlabIDP.Challenge = true
	gitlabIDP.Login = true
	gitlabIDP.MappingMethod = "claim"
	gitlabIDP.Name = "gitlab123456789"
	gitlabIDP.GitLab.URL = "https://gitlab.com/"
	gitlabIDP.GitLab.CA = &oauth.CA{Name: "gitlab-configmap"}
	gitlabIDP.GitLab.ClientID = "fake-id"
	gitlabIDP.GitLab.ClientSecret.Name = "gitlab123456789-secret"
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, gitlabIDP)

	testCases := []struct {
		name        string
		expectedCrd *oauth.CRD
	}{
		{
			name:        "build gitlab provider",
			expectedCrd: &expectedCrd,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resCrd, _, _, err := oauth.Translate(identityProviders)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedCrd, resCrd)
		})
	}
}

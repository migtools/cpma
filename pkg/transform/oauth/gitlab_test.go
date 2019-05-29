package oauth_test

import (
	"errors"
	"testing"

	"github.com/fusor/cpma/pkg/transform/oauth"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigGitlab(t *testing.T) {
	identityProviders, err := cpmatest.LoadIPTestData("testdata/gitlab/test-master-config.yaml")
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

func TestGitlabValidation(t *testing.T) {
	testCases := []struct {
		name         string
		requireError bool
		inputFile    string
		expectedErr  error
	}{
		{
			name:         "validate gitlab provider",
			requireError: false,
			inputFile:    "testdata/gitlab/test-master-config.yaml",
		},
		{
			name:         "fail on invalid name in gitlab provider",
			requireError: true,
			inputFile:    "testdata/gitlab/invalid-name-master-config.yaml",
			expectedErr:  errors.New("Name can't be empty"),
		},
		{
			name:         "fail on invalid mapping method in gitlab provider",
			requireError: true,
			inputFile:    "testdata/gitlab/invalid-mapping-master-config.yaml",
			expectedErr:  errors.New("Not valid mapping method"),
		},
		{
			name:         "fail on invalid url in gitlab provider",
			requireError: true,
			inputFile:    "testdata/gitlab/invalid-url-master-config.yaml",
			expectedErr:  errors.New("URL can't be empty"),
		},
		{
			name:         "fail on invalid clientid in gitlab provider",
			requireError: true,
			inputFile:    "testdata/gitlab/invalid-clientid-master-config.yaml",
			expectedErr:  errors.New("Client ID can't be empty"),
		},
		{
			name:         "fail on invalid client secret in gitlab provider",
			requireError: true,
			inputFile:    "testdata/gitlab/invalid-clientsecret-master-config.yaml",
			expectedErr:  errors.New("Client Secret can't be empty"),
		},
		{
			name:         "fail if key file is present for client secret in gitlab provider",
			requireError: true,
			inputFile:    "testdata/gitlab/invalid-keyfile-master-config.yaml",
			expectedErr:  errors.New("Usage of encrypted files as secret value is not supported"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			identityProvider, err := cpmatest.LoadIPTestData(tc.inputFile)
			require.NoError(t, err)

			err = oauth.Validate(identityProvider)

			if tc.requireError {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

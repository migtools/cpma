package oauth_test

import (
	"errors"
	"testing"

	"github.com/fusor/cpma/pkg/transform/oauth"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigGithub(t *testing.T) {
	identityProviders, err := cpmatest.LoadIdentityProvidersTestData("testdata/github/test-master-config.yaml")
	require.NoError(t, err)

	var expectedCrd oauth.CRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = oauth.OAuthNamespace

	var githubIDP = &oauth.IdentityProviderGitHub{}
	githubIDP.Type = "GitHub"
	githubIDP.Challenge = false
	githubIDP.Login = true
	githubIDP.MappingMethod = "claim"
	githubIDP.Name = "github123456789"
	githubIDP.GitHub.HostName = "test.example.com"
	githubIDP.GitHub.CA = &oauth.CA{Name: "github-configmap"}
	githubIDP.GitHub.ClientID = "2d85ea3f45d6777bffd7"
	githubIDP.GitHub.Organizations = []string{"myorganization1", "myorganization2"}
	githubIDP.GitHub.Teams = []string{"myorganization1/team-a", "myorganization2/team-b"}
	githubIDP.GitHub.ClientSecret.Name = "github123456789-secret"
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, githubIDP)

	testCases := []struct {
		name        string
		expectedCrd *oauth.CRD
	}{
		{
			name:        "build github provider",
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

func TestGithubValidation(t *testing.T) {
	validIdentityProviders, err := cpmatest.LoadIdentityProvidersTestData("testdata/github/test-master-config.yaml")
	require.NoError(t, err)

	invalidNameIdentityProviders, err := cpmatest.LoadIdentityProvidersTestData("testdata/github/invalid-name-master-config.yaml")
	require.NoError(t, err)

	invalidMappingMethodIdentityProviders, err := cpmatest.LoadIdentityProvidersTestData("testdata/github/invalid-mapping-master-config.yaml")
	require.NoError(t, err)

	invalidClientIDProviders, err := cpmatest.LoadIdentityProvidersTestData("testdata/github/invalid-clientid-master-config.yaml")
	require.NoError(t, err)

	invalidClientSecretProviders, err := cpmatest.LoadIdentityProvidersTestData("testdata/github/invalid-clientsecret-master-config.yaml")
	require.NoError(t, err)

	testCases := []struct {
		name         string
		requireError bool
		inputData    []oauth.IdentityProvider
		expectedErr  error
	}{
		{
			name:         "validate github provider",
			requireError: false,
			inputData:    validIdentityProviders,
		},
		{
			name:         "fail on invalid name in github provider",
			requireError: true,
			inputData:    invalidNameIdentityProviders,
			expectedErr:  errors.New("Name can't be empty"),
		},
		{
			name:         "fail on invalid mapping method in github provider",
			requireError: true,
			inputData:    invalidMappingMethodIdentityProviders,
			expectedErr:  errors.New("Not valid mapping method"),
		},
		{
			name:         "fail on invalid clientid in github provider",
			requireError: true,
			inputData:    invalidClientIDProviders,
			expectedErr:  errors.New("Client ID can't be empty"),
		},
		{
			name:         "fail on invalid client secret in github provider",
			requireError: true,
			inputData:    invalidClientSecretProviders,
			expectedErr:  errors.New("Client Secret can't be empty"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = oauth.Validate(tc.inputData)

			if tc.requireError {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

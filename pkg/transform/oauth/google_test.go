package oauth_test

import (
	"errors"
	"testing"

	"github.com/fusor/cpma/pkg/transform/oauth"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigGoogle(t *testing.T) {
	identityProviders, err := cpmatest.LoadIPTestData("testdata/google/master_config.yaml")
	require.NoError(t, err)

	var expectedCrd configv1.OAuth
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Name = "cluster"
	expectedCrd.Namespace = oauth.OAuthNamespace

	var googleIDP = &configv1.IdentityProvider{}
	googleIDP.Type = "Google"
	googleIDP.MappingMethod = "claim"
	googleIDP.Name = "google123456789123456789"
	googleIDP.Google = &configv1.GoogleIdentityProvider{}
	googleIDP.Google.ClientID = "82342890327-tf5lqn4eikdf4cb4edfm85jiqotvurpq.apps.googleusercontent.com"
	googleIDP.Google.ClientSecret.Name = "google123456789123456789-secret"
	googleIDP.Google.HostedDomain = "test.example.com"
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, *googleIDP)

	testCases := []struct {
		name        string
		expectedCrd *configv1.OAuth
	}{
		{
			name:        "build google provider",
			expectedCrd: &expectedCrd,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oauthResources, err := oauth.Translate(identityProviders, oauth.TokenConfig{})
			require.NoError(t, err)
			assert.Equal(t, tc.expectedCrd, oauthResources.OAuthCRD)
		})
	}
}

func TestGoogleValidation(t *testing.T) {
	testCases := []struct {
		name         string
		requireError bool
		inputFile    string
		expectedErr  error
	}{
		{
			name:         "validate google provider",
			requireError: false,
			inputFile:    "testdata/google/master_config.yaml",
		},
		{
			name:         "fail on invalid name in google provider",
			requireError: true,
			inputFile:    "testdata/google/invalid-name-master-config.yaml",
			expectedErr:  errors.New("Name can't be empty"),
		},
		{
			name:         "fail on invalid mapping method in google provider",
			requireError: true,
			inputFile:    "testdata/google/invalid-mapping-master-config.yaml",
			expectedErr:  errors.New("Not valid mapping method"),
		},
		{
			name:         "fail on invalid clientid in google provider",
			requireError: true,
			inputFile:    "testdata/google/invalid-clientid-master-config.yaml",
			expectedErr:  errors.New("Client ID can't be empty"),
		},
		{
			name:         "fail on invalid client secret in google provider",
			requireError: true,
			inputFile:    "testdata/google/invalid-clientsecret-master-config.yaml",
			expectedErr:  errors.New("Client Secret can't be empty"),
		},
		{
			name:         "fail if key file is present for client secret in google provider",
			requireError: true,
			inputFile:    "testdata/google/invalid-keyfile-master-config.yaml",
			expectedErr:  errors.New("Usage of encrypted files as secret value is not supported"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			identityProvider, err := cpmatest.LoadIPTestData(tc.inputFile)
			require.NoError(t, err)

			err = oauth.Validate(identityProvider)

			if tc.requireError {
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

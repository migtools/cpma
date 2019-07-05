package oauth_test

import (
	"errors"
	"testing"

	"github.com/fusor/cpma/pkg/transform/oauth"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigGoogle(t *testing.T) {
	t.Parallel()
	identityProviders, _, err := cpmatest.LoadIPTestData("testdata/google/master_config.yaml")
	require.NoError(t, err)

	expectedCrd, err := loadExpectedOAuth("testdata/google/expected-CR-oauth.yaml")
	require.NoError(t, err)

	testCases := []struct {
		name        string
		expectedCrd *configv1.OAuth
	}{
		{
			name:        "build google provider",
			expectedCrd: expectedCrd,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oauthResources, err := oauth.Translate(identityProviders, oauth.TokenConfig{}, legacyconfigv1.OAuthTemplates{})
			require.NoError(t, err)
			assert.Equal(t, tc.expectedCrd, oauthResources.OAuthCRD)
		})
	}
}

func TestGoogleValidation(t *testing.T) {
	t.Parallel()
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
			identityProvider, _, err := cpmatest.LoadIPTestData(tc.inputFile)
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

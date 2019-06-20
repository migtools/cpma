package oauth_test

import (
	"testing"

	"github.com/fusor/cpma/pkg/transform/oauth"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigBasicAuth(t *testing.T) {
	identityProviders, _, err := cpmatest.LoadIPTestData("testdata/basicauth/master_config.yaml")
	require.NoError(t, err)

	expectedCrd, err := loadExpectedOAuth("testdata/basicauth/expected-CR-oauth.yaml")
	require.NoError(t, err)

	testCases := []struct {
		name        string
		expectedCrd *configv1.OAuth
	}{
		{
			name:        "build basic auth provider",
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

func TestBasicAuthValidation(t *testing.T) {
	testCases := []struct {
		name         string
		requireError bool
		inputFile    string
		expectedErr  error
	}{
		{
			name:         "validate basic auth provider",
			requireError: false,
			inputFile:    "testdata/basicauth/master_config.yaml",
		},
		{
			name:         "fail on invalid name in basic auth provider",
			requireError: true,
			inputFile:    "testdata/basicauth/invalid-name-master-config.yaml",
			expectedErr:  errors.New("Name can't be empty"),
		},
		{
			name:         "fail on invalid mapping method in basic auth provider",
			requireError: true,
			inputFile:    "testdata/basicauth/invalid-mapping-master-config.yaml",
			expectedErr:  errors.New("Not valid mapping method"),
		},
		{
			name:         "fail on invalid url in basic auth provider",
			requireError: true,
			inputFile:    "testdata/basicauth/invalid-url-master-config.yaml",
			expectedErr:  errors.New("URL can't be empty"),
		},
		{
			name:         "fail on invalid key file in basic auth provider",
			requireError: true,
			inputFile:    "testdata/basicauth/invalid-keyfile-master-config.yaml",
			expectedErr:  errors.New("Key file can't be empty if cert file is specified"),
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

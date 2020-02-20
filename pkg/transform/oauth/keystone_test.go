package oauth_test

import (
	"errors"
	"testing"

	cpmatest "github.com/fusor/cpma/pkg/transform/internal/test"
	"github.com/fusor/cpma/pkg/transform/oauth"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigKeystone(t *testing.T) {
	t.Parallel()
	identityProviders, _, err := cpmatest.LoadIPTestData("testdata/keystone/master_config.yaml")
	require.NoError(t, err)

	expectedCrd, err := loadExpectedOAuth("testdata/keystone/expected-CR-oauth.yaml")
	require.NoError(t, err)

	testCases := []struct {
		name        string
		expectedCrd *configv1.OAuth
	}{
		{
			name:        "build keystone provider",
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

func TestKeystoneValidation(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		requireError bool
		inputFile    string
		expectedErr  error
	}{
		{
			name:         "validate keystone provider",
			requireError: false,
			inputFile:    "testdata/keystone/master_config.yaml",
		},
		{
			name:         "fail on invalid name in keystone provider",
			requireError: true,
			inputFile:    "testdata/keystone/invalid-name-master-config.yaml",
			expectedErr:  errors.New("Name can't be empty"),
		},
		{
			name:         "fail on invalid mapping method in keystone provider",
			requireError: true,
			inputFile:    "testdata/keystone/invalid-mapping-master-config.yaml",
			expectedErr:  errors.New("Not valid mapping method"),
		},
		{
			name:         "fail on invalid url in keystone provider",
			requireError: true,
			inputFile:    "testdata/keystone/invalid-url-master-config.yaml",
			expectedErr:  errors.New("URL can't be empty"),
		},
		{
			name:         "fail on invalid domain name in keystone provider",
			requireError: true,
			inputFile:    "testdata/keystone/invalid-domainname-master-config.yaml",
			expectedErr:  errors.New("Domain name can't be empty"),
		},
		{
			name:         "fail on invalid key file in keystone provider",
			requireError: true,
			inputFile:    "testdata/keystone/invalid-keyfile-master-config.yaml",
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

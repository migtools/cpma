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

func TestTransformMasterConfigLDAP(t *testing.T) {
	t.Parallel()
	identityProviders, _, err := cpmatest.LoadIPTestData("testdata/ldap/master_config.yaml")
	require.NoError(t, err)

	expectedCrd, err := loadExpectedOAuth("testdata/ldap/expected-CR-oauth.yaml")
	require.NoError(t, err)

	testCases := []struct {
		name        string
		expectedCrd *configv1.OAuth
	}{
		{
			name:        "build ldap provider",
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

func TestLDAPValidation(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		requireError bool
		inputFile    string
		expectedErr  error
	}{
		{
			name:         "validate ldap provider",
			requireError: false,
			inputFile:    "testdata/ldap/master_config.yaml",
		},
		{
			name:         "fail on invalid name in ldap provider",
			requireError: true,
			inputFile:    "testdata/ldap/invalid-name-master-config.yaml",
			expectedErr:  errors.New("Name can't be empty"),
		},
		{
			name:         "fail on invalid mapping method in ldap provider",
			requireError: true,
			inputFile:    "testdata/ldap/invalid-mapping-master-config.yaml",
			expectedErr:  errors.New("Not valid mapping method"),
		},
		{
			name:         "fail on invalid ids in ldap provider",
			requireError: true,
			inputFile:    "testdata/ldap/invalid-ids-master-config.yaml",
			expectedErr:  errors.New("ID can't be empty"),
		},
		{
			name:         "fail on invalid emails in ldap provider",
			requireError: true,
			inputFile:    "testdata/ldap/invalid-emails-master-config.yaml",
			expectedErr:  errors.New("Email can't be empty"),
		},
		{
			name:         "fail on invalid names in ldap provider",
			requireError: true,
			inputFile:    "testdata/ldap/invalid-names-master-config.yaml",
			expectedErr:  errors.New("Name can't be empty"),
		},
		{
			name:         "fail on invalid preferred usernames in ldap provider",
			requireError: true,
			inputFile:    "testdata/ldap/invalid-usernames-master-config.yaml",
			expectedErr:  errors.New("Preferred username can't be empty"),
		},
		{
			name:         "fail on invalid url in ldap provider",
			requireError: true,
			inputFile:    "testdata/ldap/invalid-url-master-config.yaml",
			expectedErr:  errors.New("URL can't be empty"),
		},
		{
			name:         "fail if key file is present for bind password in ldap provider",
			requireError: true,
			inputFile:    "testdata/ldap/invalid-bpass-master-config.yaml",
			expectedErr:  errors.New("Usage of encrypted files as bind password value is not supported"),
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

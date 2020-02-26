package oauth_test

import (
	"errors"
	"testing"

	cpmatest "github.com/konveyor/cpma/pkg/transform/internal/test"
	"github.com/konveyor/cpma/pkg/transform/oauth"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigRequestHeader(t *testing.T) {
	t.Parallel()
	identityProviders, _, err := cpmatest.LoadIPTestData("testdata/requestheader/master_config.yaml")
	require.NoError(t, err)

	expectedCrd, err := loadExpectedOAuth("testdata/requestheader/expected-CR-oauth.yaml")
	require.NoError(t, err)

	testCases := []struct {
		name        string
		expectedCrd *configv1.OAuth
	}{
		{
			name:        "build request header provider",
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

func TestRequestHeaderValidation(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		requireError bool
		inputFile    string
		expectedErr  error
	}{
		{
			name:         "validate requestheader provider",
			requireError: false,
			inputFile:    "testdata/requestheader/master_config.yaml",
		},
		{
			name:         "fail on invalid name in requestheader provider",
			requireError: true,
			inputFile:    "testdata/requestheader/invalid-name-master-config.yaml",
			expectedErr:  errors.New("Name can't be empty"),
		},
		{
			name:         "fail on invalid mapping method in requestheader provider",
			requireError: true,
			inputFile:    "testdata/requestheader/invalid-mapping-master-config.yaml",
			expectedErr:  errors.New("Not valid mapping method"),
		},
		{
			name:         "fail on invalid headers in requestheader provider",
			requireError: true,
			inputFile:    "testdata/requestheader/invalid-headers-master-config.yaml",
			expectedErr:  errors.New("Headers can't be empty"),
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

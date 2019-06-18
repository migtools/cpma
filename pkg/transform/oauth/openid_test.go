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

func TestTransformMasterConfigOpenID(t *testing.T) {
	identityProviders, err := cpmatest.LoadIPTestData("testdata/openid/master_config.yaml")
	require.NoError(t, err)

	var expectedCrd configv1.OAuth
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Name = "cluster"
	expectedCrd.Namespace = oauth.OAuthNamespace

	var openidIDP = &configv1.IdentityProvider{}
	openidIDP.Type = "OpenID"
	openidIDP.MappingMethod = "claim"
	openidIDP.Name = "my_openid_connect"
	openidIDP.OpenID = &configv1.OpenIDIdentityProvider{}
	openidIDP.OpenID.ClientID = "testid"
	openidIDP.OpenID.Claims.PreferredUsername = []string{"preferred_username", "email"}
	openidIDP.OpenID.Claims.Name = []string{"nickname", "given_name", "name"}
	openidIDP.OpenID.Claims.Email = []string{"custom_email_claim", "email"}
	openidIDP.OpenID.ClientSecret.Name = "my_openid_connect-secret"

	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, *openidIDP)

	testCases := []struct {
		name        string
		expectedCrd *configv1.OAuth
	}{
		{
			name:        "build openid provider",
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

func TestOpenIDValidation(t *testing.T) {
	testCases := []struct {
		name         string
		requireError bool
		inputFile    string
		expectedErr  error
	}{
		{
			name:         "validate openid provider",
			requireError: false,
			inputFile:    "testdata/openid/master_config.yaml",
		},
		{
			name:         "fail on invalid name in openid provider",
			requireError: true,
			inputFile:    "testdata/openid/invalid-name-master-config.yaml",
			expectedErr:  errors.New("Name can't be empty"),
		},
		{
			name:         "fail on invalid mapping method in openid provider",
			requireError: true,
			inputFile:    "testdata/openid/invalid-mapping-master-config.yaml",
			expectedErr:  errors.New("Not valid mapping method"),
		},
		{
			name:         "fail on invalid clientid in openid provider",
			requireError: true,
			inputFile:    "testdata/openid/invalid-clientid-master-config.yaml",
			expectedErr:  errors.New("Client ID can't be empty"),
		},
		{
			name:         "fail on invalid client secret in openid provider",
			requireError: true,
			inputFile:    "testdata/openid/invalid-clientsecret-master-config.yaml",
			expectedErr:  errors.New("Client Secret can't be empty"),
		},
		{
			name:         "fail on invalid claims in openid provider",
			requireError: true,
			inputFile:    "testdata/openid/invalid-claims-master-config.yaml",
			expectedErr:  errors.New("All claims are empty. At least one is required"),
		},
		{
			name:         "fail on invalid auth endpoint in openid provider",
			requireError: true,
			inputFile:    "testdata/openid/invalid-auth-master-config.yaml",
			expectedErr:  errors.New("Authorization endpoint can't be empty"),
		},
		{
			name:         "fail on invalid token endpoint in openid provider",
			requireError: true,
			inputFile:    "testdata/openid/invalid-token-master-config.yaml",
			expectedErr:  errors.New("Token endpoint can't be empty"),
		},
		{
			name:         "fail if key file is present for client secret in openid provider",
			requireError: true,
			inputFile:    "testdata/openid/invalid-keyfile-master-config.yaml",
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

package oauth_test

import (
	"errors"
	"testing"

	"github.com/fusor/cpma/pkg/config"
	"github.com/fusor/cpma/pkg/transform/oauth"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigOpenID(t *testing.T) {
	config := config.LoadConfig()
	identityProviders, err := cpmatest.LoadIPTestData("testdata/openid/test-master-config.yaml")
	require.NoError(t, err)

	var expectedCrd oauth.CRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = oauth.OAuthNamespace

	var openidIDP = &oauth.IdentityProviderOpenID{}
	openidIDP.Type = "OpenID"
	openidIDP.Challenge = false
	openidIDP.Login = true
	openidIDP.MappingMethod = "claim"
	openidIDP.Name = "my_openid_connect"
	openidIDP.OpenID.ClientID = "testid"
	openidIDP.OpenID.Claims.PreferredUsername = []string{"preferred_username", "email"}
	openidIDP.OpenID.Claims.Name = []string{"nickname", "given_name", "name"}
	openidIDP.OpenID.Claims.Email = []string{"custom_email_claim", "email"}
	openidIDP.OpenID.URLs.Authorize = "https://myidp.example.com/oauth2/authorize"
	openidIDP.OpenID.URLs.Token = "https://myidp.example.com/oauth2/token"
	openidIDP.OpenID.ClientSecret.Name = "my_openid_connect-secret"

	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, openidIDP)

	testCases := []struct {
		name        string
		expectedCrd *oauth.CRD
	}{
		{
			name:        "build openid provider",
			expectedCrd: &expectedCrd,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resCrd, _, _, err := oauth.Translate(identityProviders, &config)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedCrd, resCrd)
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
			inputFile:    "testdata/openid/test-master-config.yaml",
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
				assert.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

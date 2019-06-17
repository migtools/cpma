package oauth_test

import (
	"testing"

	"github.com/fusor/cpma/pkg/transform/oauth"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigBasicAuth(t *testing.T) {
	identityProviders, err := cpmatest.LoadIPTestData("testdata/basicauth/test-master-config.yaml")
	require.NoError(t, err)

	var expectedCrd oauth.CRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = oauth.OAuthNamespace

	var basicAuthIDP = &oauth.IdentityProviderBasicAuth{}
	basicAuthIDP.Type = "BasicAuth"
	basicAuthIDP.Challenge = true
	basicAuthIDP.Login = true
	basicAuthIDP.Name = "my_remote_basic_auth_provider"
	basicAuthIDP.MappingMethod = "claim"
	basicAuthIDP.BasicAuth.URL = "https://www.example.com/"
	basicAuthIDP.BasicAuth.TLSClientCert = &oauth.TLSClientCert{Name: "my_remote_basic_auth_provider-client-cert-secret"}
	basicAuthIDP.BasicAuth.TLSClientKey = &oauth.TLSClientKey{Name: "my_remote_basic_auth_provider-client-key-secret"}
	basicAuthIDP.BasicAuth.CA = &oauth.CA{Name: "basicauth-configmap"}

	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, basicAuthIDP)

	testCases := []struct {
		name        string
		expectedCrd *oauth.CRD
	}{
		{
			name:        "build basic auth provider",
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
			inputFile:    "testdata/basicauth/test-master-config.yaml",
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

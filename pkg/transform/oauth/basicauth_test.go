package oauth_test

import (
	"errors"
	"testing"

	"github.com/fusor/cpma/pkg/transform/oauth"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigBasicAuth(t *testing.T) {
	identityProviders, err := cpmatest.LoadIdentityProvidersTestData("testdata/basicauth/test-master-config.yaml")
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
			resCrd, _, _, err := oauth.Translate(identityProviders)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedCrd, resCrd)
		})
	}
}

func TestBasicAuthValidation(t *testing.T) {
	validIdentityProviders, err := cpmatest.LoadIdentityProvidersTestData("testdata/basicauth/test-master-config.yaml")
	require.NoError(t, err)

	invalidNameIdentityProviders, err := cpmatest.LoadIdentityProvidersTestData("testdata/basicauth/invalid-name-master-config.yaml")
	require.NoError(t, err)

	invalidMappingMethodIdentityProviders, err := cpmatest.LoadIdentityProvidersTestData("testdata/basicauth/invalid-mapping-master-config.yaml")
	require.NoError(t, err)

	invalidURLIdentityProviders, err := cpmatest.LoadIdentityProvidersTestData("testdata/basicauth/invalid-url-master-config.yaml")
	require.NoError(t, err)

	invalidKeyFileIdentityProviders, err := cpmatest.LoadIdentityProvidersTestData("testdata/basicauth/invalid-keyfile-master-config.yaml")
	require.NoError(t, err)

	testCases := []struct {
		name         string
		requireError bool
		inputData    []oauth.IdentityProvider
		expectedErr  error
	}{
		{
			name:         "validate basic auth provider",
			requireError: false,
			inputData:    validIdentityProviders,
		},
		{
			name:         "fail on invalid name in basic auth provider",
			requireError: true,
			inputData:    invalidNameIdentityProviders,
			expectedErr:  errors.New("Name can't be empty"),
		},
		{
			name:         "fail on invalid mapping method in basic auth provider",
			requireError: true,
			inputData:    invalidMappingMethodIdentityProviders,
			expectedErr:  errors.New("Not valid mapping method"),
		},
		{
			name:         "fail on invalid url in basic auth provider",
			requireError: true,
			inputData:    invalidURLIdentityProviders,
			expectedErr:  errors.New("URL can't be empty"),
		},
		{
			name:         "fail on invalid key file in basic auth provider",
			requireError: true,
			inputData:    invalidKeyFileIdentityProviders,
			expectedErr:  errors.New("Key file can't be empty if cert file is specified"),
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

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

func TestTransformMasterConfigRequestHeader(t *testing.T) {
	config := config.LoadConfig()
	identityProviders, err := cpmatest.LoadIPTestData("testdata/requestheader/test-master-config.yaml")
	require.NoError(t, err)

	var expectedCrd oauth.CRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = oauth.OAuthNamespace

	var requestHeaderIDP = &oauth.IdentityProviderRequestHeader{}

	requestHeaderIDP.Type = "RequestHeader"
	requestHeaderIDP.Name = "my_request_header_provider"
	requestHeaderIDP.Challenge = true
	requestHeaderIDP.Login = true
	requestHeaderIDP.MappingMethod = "claim"
	requestHeaderIDP.RequestHeader.ChallengeURL = "https://example.com"
	requestHeaderIDP.RequestHeader.LoginURL = "https://example.com"
	requestHeaderIDP.RequestHeader.CA = &oauth.CA{Name: "requestheader-configmap"}
	requestHeaderIDP.RequestHeader.ClientCommonNames = []string{"my-auth-proxy"}
	requestHeaderIDP.RequestHeader.Headers = []string{"X-Remote-User", "SSO-User"}
	requestHeaderIDP.RequestHeader.EmailHeaders = []string{"X-Remote-User-Email"}
	requestHeaderIDP.RequestHeader.NameHeaders = []string{"X-Remote-User-Display-Name"}
	requestHeaderIDP.RequestHeader.PreferredUsernameHeaders = []string{"X-Remote-User-Login"}
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, requestHeaderIDP)

	testCases := []struct {
		name        string
		expectedCrd *oauth.CRD
	}{
		{
			name:        "build request header provider",
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

func TestRequestHeaderValidation(t *testing.T) {
	testCases := []struct {
		name         string
		requireError bool
		inputFile    string
		expectedErr  error
	}{
		{
			name:         "validate requestheader provider",
			requireError: false,
			inputFile:    "testdata/requestheader/test-master-config.yaml",
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

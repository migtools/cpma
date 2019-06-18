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

func TestTransformMasterConfigRequestHeader(t *testing.T) {
	identityProviders, err := cpmatest.LoadIPTestData("testdata/requestheader/master_config.yaml")
	require.NoError(t, err)

	var expectedCrd configv1.OAuth
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Name = "cluster"
	expectedCrd.Namespace = oauth.OAuthNamespace

	var requestHeaderIDP = &configv1.IdentityProvider{}

	requestHeaderIDP.Type = "RequestHeader"
	requestHeaderIDP.Name = "my_request_header_provider"
	requestHeaderIDP.MappingMethod = "claim"
	requestHeaderIDP.RequestHeader = &configv1.RequestHeaderIdentityProvider{}
	requestHeaderIDP.RequestHeader.ChallengeURL = "https://example.com"
	requestHeaderIDP.RequestHeader.LoginURL = "https://example.com"
	requestHeaderIDP.RequestHeader.ClientCA = configv1.ConfigMapNameReference{Name: "requestheader-configmap"}
	requestHeaderIDP.RequestHeader.ClientCommonNames = []string{"my-auth-proxy"}
	requestHeaderIDP.RequestHeader.Headers = []string{"X-Remote-User", "SSO-User"}
	requestHeaderIDP.RequestHeader.EmailHeaders = []string{"X-Remote-User-Email"}
	requestHeaderIDP.RequestHeader.NameHeaders = []string{"X-Remote-User-Display-Name"}
	requestHeaderIDP.RequestHeader.PreferredUsernameHeaders = []string{"X-Remote-User-Login"}
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, *requestHeaderIDP)

	testCases := []struct {
		name        string
		expectedCrd *configv1.OAuth
	}{
		{
			name:        "build request header provider",
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

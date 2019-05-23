package oauth_test

import (
	"testing"

	"github.com/fusor/cpma/pkg/transform/oauth"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigRequestHeader(t *testing.T) {
	identityProviders, err := cpmatest.LoadIdentityProvidersTestData("testdata/requestheader-test-master-config.yaml")
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
			resCrd, _, _, err := oauth.Translate(identityProviders)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedCrd, resCrd)
		})
	}
}

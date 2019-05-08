package oauth

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fusor/cpma/pkg/ocp3"
)

func TestTransformMasterConfigRequestHeader(t *testing.T) {
	file := "testdata/requestheader-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)

	masterV3 := ocp3.Master{}
	masterV3.Decode(content)

	var expectedCrd OAuthCRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = "openshift-config"

	var requestHeaderIDP identityProviderRequestHeader

	requestHeaderIDP.Type = "RequestHeader"
	requestHeaderIDP.Name = "my_request_header_provider"
	requestHeaderIDP.Challenge = true
	requestHeaderIDP.Login = true
	requestHeaderIDP.MappingMethod = "claim"
	requestHeaderIDP.RequestHeader.ChallengeURL = "https://example.com"
	requestHeaderIDP.RequestHeader.LoginURL = "https://example.com"
	requestHeaderIDP.RequestHeader.CA.Name = "cert.crt"
	requestHeaderIDP.RequestHeader.ClientCommonNames = []string{"my-auth-proxy"}
	requestHeaderIDP.RequestHeader.Headers = []string{"X-Remote-User", "SSO-User"}
	requestHeaderIDP.RequestHeader.EmailHeaders = []string{"X-Remote-User-Email"}
	requestHeaderIDP.RequestHeader.NameHeaders = []string{"X-Remote-User-Display-Name"}
	requestHeaderIDP.RequestHeader.PreferredUsernameHeaders = []string{"X-Remote-User-Login"}
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, requestHeaderIDP)

	resCrd, _, err := Transform(masterV3.Config.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, &expectedCrd, resCrd)
}

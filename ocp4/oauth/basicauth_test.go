package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fusor/cpma/ocp3"
)

func TestTranslateMasterConfigBasicAuth(t *testing.T) {
	testConfig := ocp3.Config{
		Masterf: "../../test/oauth/basicauth-test-master-config.yaml",
	}
	masterConfig := testConfig.ParseMaster()

	var expectedCrd OAuthCRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.MetaData.Name = "cluster"
	expectedCrd.MetaData.NameSpace = "openshift-config"

	var basicAuthIDP identityProviderBasicAuth
	basicAuthIDP.Type = "BasicAuth"
	basicAuthIDP.Challenge = true
	basicAuthIDP.Login = true
	basicAuthIDP.Name = "my_remote_basic_auth_provider"
	basicAuthIDP.MappingMethod = "claim"
	basicAuthIDP.BasicAuth.URL = "https://www.example.com/"
	basicAuthIDP.BasicAuth.CA.Name = "ca.file"
	basicAuthIDP.BasicAuth.TLSClientCert.Name = "my_remote_basic_auth_provider-client-cert-secret"
	basicAuthIDP.BasicAuth.TLSClientKey.Name = "my_remote_basic_auth_provider-client-key-secret"

	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, basicAuthIDP)

	resCrd, _, err := Translate(masterConfig.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, &expectedCrd, resCrd)
}

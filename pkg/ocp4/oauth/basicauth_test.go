package oauth

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fusor/cpma/pkg/ocp3"
)

func TestTranslateMasterConfigBasicAuth(t *testing.T) {
	file := "testdata/basicauth-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)

	masterV3 := ocp3.Master{}
	masterV3.Decode(content)

	var expectedCrd OAuthCRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = "openshift-config"

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

	resCrd, _, err := Translate(masterV3.Config.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, &expectedCrd, resCrd)
}

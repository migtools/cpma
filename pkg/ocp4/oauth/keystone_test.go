package oauth_test

import (
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/fusor/cpma/pkg/ocp4/oauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigKeystone(t *testing.T) {
	defer func() { io.GetFile = _GetFile }()
	oauth.GetFile = mockGetFile

	file := "testdata/keystone-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)
	masterV3 := ocp3.MasterDecode(content)

	var expectedCrd oauth.OAuthCRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = "openshift-config"

	var keystoneIDP oauth.IdentityProviderKeystone
	keystoneIDP.Type = "Keystone"
	keystoneIDP.Challenge = true
	keystoneIDP.Login = true
	keystoneIDP.Name = "my_keystone_provider"
	keystoneIDP.MappingMethod = "claim"
	keystoneIDP.Keystone.DomainName = "default"
	keystoneIDP.Keystone.URL = "http://fake.url:5000"
	keystoneIDP.Keystone.CA.Name = "keystone.pem"
	keystoneIDP.Keystone.TLSClientCert.Name = "my_keystone_provider-client-cert-secret"
	keystoneIDP.Keystone.TLSClientKey.Name = "my_keystone_provider-client-key-secret"

	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, keystoneIDP)

	resCrd, _, err := oauth.Transform(masterV3.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, &expectedCrd, resCrd)
}

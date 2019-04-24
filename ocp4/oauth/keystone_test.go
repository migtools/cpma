package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fusor/cpma/ocp3"
)

func TestTranslateMasterConfigKeystone(t *testing.T) {
	masterConfig := ocp3.ParseMaster("../../test/oauth/keystone-test-master-config.yaml")

	var expectedCrd OAuthCRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = "openshift-config"

	var keystoneIDP identityProviderKeystone
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

	resCrd, _, err := Translate(masterConfig.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, &expectedCrd, resCrd)
}

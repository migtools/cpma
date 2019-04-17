package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocp3 "github.com/fusor/cpma/ocp3config"
)

func TestGenerateMasterConfigHtpasswd(t *testing.T) {
	testConfig := ocp3.Config{
		Masterf: "../../test/oauth/htpasswd-test-master-config.yaml",
	}
	masterConfig := testConfig.ParseMaster()

	var expectedCrd v4OAuthCRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.MetaData.Name = "cluster"

	var htpasswdIDP identityProviderHTPasswd
	htpasswdIDP.Type = "HTPasswd"
	htpasswdIDP.Challenge = true
	htpasswdIDP.Login = true
	htpasswdIDP.MappingMethod = "claim"
	htpasswdIDP.HTPasswd.FileData.Name = "/etc/origin/master/htpasswd"
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, htpasswdIDP)

	resCrd, err := Generate(masterConfig)
	require.NoError(t, err)
	assert.Equal(t, &expectedCrd, resCrd)
}

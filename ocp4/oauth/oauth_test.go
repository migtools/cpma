package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocp3 "github.com/fusor/cpma/ocp3config"
)

func TestTranslateMasterConfig(t *testing.T) {
	testConfig := ocp3.Config{
		Masterf: "../../test/oauth/bulk-test-master-config.yaml",
	}
	masterConfig := testConfig.ParseMaster()

	resCrd, _, err := Translate(masterConfig.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, resCrd.Spec.IdentityProviders[0].(identityProviderHTPasswd).Type, "HTPasswd")
	assert.Equal(t, resCrd.Spec.IdentityProviders[1].(identityProviderGitHub).Type, "GitHub")
}

package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fusor/cpma/ocp3"
)

func TestTranslateMasterConfig(t *testing.T) {
	testConfig := ocp3.Config{
		Masterf: "../../test/oauth/bulk-test-master-config.yaml",
	}
	masterConfig := testConfig.ParseMaster()

	resCrd, _, err := Translate(masterConfig.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, resCrd.Spec.IdentityProviders[0].(identityProviderGitHub).Type, "GitHub")
	assert.Equal(t, resCrd.Spec.IdentityProviders[1].(identityProviderGitLab).Type, "GitLab")
	assert.Equal(t, resCrd.Spec.IdentityProviders[2].(identityProviderHTPasswd).Type, "HTPasswd")
}

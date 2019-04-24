package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fusor/cpma/ocp3"
)

func TestTranslateMasterConfigHtpasswd(t *testing.T) {
	masterConfig := ocp3.ParseMaster("../../test/oauth/htpasswd-test-master-config.yaml")

	var expectedCrd OAuthCRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = "openshift-config"

	var htpasswdIDP identityProviderHTPasswd
	htpasswdIDP.Name = "htpasswd_auth"
	htpasswdIDP.Type = "HTPasswd"
	htpasswdIDP.Challenge = true
	htpasswdIDP.Login = true
	htpasswdIDP.MappingMethod = "claim"
	htpasswdIDP.HTPasswd.FileData.Name = "htpasswd_auth-secret"
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, htpasswdIDP)

	resCrd, _, err := Translate(masterConfig.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, &expectedCrd, resCrd)
}

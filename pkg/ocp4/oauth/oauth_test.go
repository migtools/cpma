package oauth

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fusor/cpma/internal/io"
	"github.com/fusor/cpma/pkg/ocp3"
)

var _GetFile = io.GetFile

func mockGetFile(host, src, dst string) []byte {
	return []byte("This is test file content")
}

func TestTranslateMasterConfig(t *testing.T) {
	defer func() { GetFile = _GetFile }()
	GetFile = mockGetFile

	file := "testdata/bulk-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)

	masterV3 := ocp3.Master{}
	masterV3.Decode(content)

	resCrd, _, err := Translate(masterV3.Config.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, len(resCrd.Spec.IdentityProviders), 9)
	assert.Equal(t, resCrd.Spec.IdentityProviders[0].(identityProviderBasicAuth).Type, "BasicAuth")
	assert.Equal(t, resCrd.Spec.IdentityProviders[1].(identityProviderGitHub).Type, "GitHub")
	assert.Equal(t, resCrd.Spec.IdentityProviders[2].(identityProviderGitLab).Type, "GitLab")
	assert.Equal(t, resCrd.Spec.IdentityProviders[3].(identityProviderGoogle).Type, "Google")
	assert.Equal(t, resCrd.Spec.IdentityProviders[4].(identityProviderHTPasswd).Type, "HTPasswd")
	assert.Equal(t, resCrd.Spec.IdentityProviders[5].(identityProviderKeystone).Type, "Keystone")
	assert.Equal(t, resCrd.Spec.IdentityProviders[6].(identityProviderLDAP).Type, "LDAP")
	assert.Equal(t, resCrd.Spec.IdentityProviders[7].(identityProviderRequestHeader).Type, "RequestHeader")
	assert.Equal(t, resCrd.Spec.IdentityProviders[8].(identityProviderOpenID).Type, "OpenID")
}

func TestGenYAML(t *testing.T) {
	defer func() { GetFile = _GetFile }()
	GetFile = mockGetFile

	file := "testdata/bulk-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)

	masterV3 := ocp3.Master{}
	masterV3.Decode(content)

	crd, manifests, err := Translate(masterV3.Config.OAuthConfig)
	require.NoError(t, err)

	CRD := crd.GenYAML()
	expectedYaml, _ := ioutil.ReadFile("testdata/expected-bulk-test-masterconfig-oauth.yaml")

	assert.Equal(t, len(manifests), 9)
	assert.Equal(t, expectedYaml, CRD)
}

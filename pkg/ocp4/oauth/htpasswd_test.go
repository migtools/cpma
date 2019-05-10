package oauth_test

import (
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/internal/io"
	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/fusor/cpma/pkg/ocp4/oauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigHtpasswd(t *testing.T) {
	defer func() { io.GetFile = _GetFile }()
	oauth.GetFile = mockGetFile

	file := "testdata/htpasswd-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)
	masterV3 := ocp3.MasterDecode(content)

	var expectedCrd oauth.OAuthCRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = "openshift-config"

	var htpasswdIDP oauth.IdentityProviderHTPasswd
	htpasswdIDP.Name = "htpasswd_auth"
	htpasswdIDP.Type = "HTPasswd"
	htpasswdIDP.Challenge = true
	htpasswdIDP.Login = true
	htpasswdIDP.MappingMethod = "claim"
	htpasswdIDP.HTPasswd.FileData.Name = "htpasswd_auth-secret"
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, htpasswdIDP)

	resCrd, _, err := oauth.Transform(masterV3.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, &expectedCrd, resCrd)
}

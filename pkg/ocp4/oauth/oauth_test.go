package oauth_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/fusor/cpma/pkg/ocp4/oauth"
)

var _GetFile = io.GetFile

func mockGetFile(a, b, c string) []byte {
	return []byte("This is test file content")
}

func TestTransformMasterConfig(t *testing.T) {
	defer func() { io.GetFile = _GetFile }()
	oauth.GetFile = mockGetFile

	file := "testdata/bulk-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)
	masterV3 := ocp3.MasterDecode(content)

	resCrd, _, err := oauth.Transform(masterV3.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, len(resCrd.Spec.IdentityProviders), 9)
	assert.Equal(t, resCrd.Spec.IdentityProviders[0].(oauth.IdentityProviderBasicAuth).Type, "BasicAuth")
	assert.Equal(t, resCrd.Spec.IdentityProviders[1].(oauth.IdentityProviderGitHub).Type, "GitHub")
	assert.Equal(t, resCrd.Spec.IdentityProviders[2].(oauth.IdentityProviderGitLab).Type, "GitLab")
	assert.Equal(t, resCrd.Spec.IdentityProviders[3].(oauth.IdentityProviderGoogle).Type, "Google")
	assert.Equal(t, resCrd.Spec.IdentityProviders[4].(oauth.IdentityProviderHTPasswd).Type, "HTPasswd")
	assert.Equal(t, resCrd.Spec.IdentityProviders[5].(oauth.IdentityProviderKeystone).Type, "Keystone")
	assert.Equal(t, resCrd.Spec.IdentityProviders[6].(oauth.IdentityProviderLDAP).Type, "LDAP")
	assert.Equal(t, resCrd.Spec.IdentityProviders[7].(oauth.IdentityProviderRequestHeader).Type, "RequestHeader")
	assert.Equal(t, resCrd.Spec.IdentityProviders[8].(oauth.IdentityProviderOpenID).Type, "OpenID")
}

func TestGenYAML(t *testing.T) {
	defer func() { oauth.GetFile = _GetFile }()
	oauth.GetFile = mockGetFile

	file := "testdata/bulk-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)
	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	var masterV3 configv1.MasterConfig
	_, _, _ = serializer.Decode(content, nil, &masterV3)

	var identityProviders []ocp3.IdentityProvider
	for _, identityProvider := range masterV3.OAuthConfig.IdentityProviders {
		providerJSON, _ := identityProvider.Provider.MarshalJSON()
		provider := ocp.Provider{}
		json.Unmarshal(providerJSON, &provider)

		identityProviders = append(identityProviders,
			ocp3.IdentityProvider{
				provider.Kind,
				provider.APIVersion,
				identityProvider.MappingMethod,
				identityProvider.Name,
				identityProvider.Provider,
				provider.File,
				nil,
				identityProvider.UseAsChallenger,
				identityProvider.UseAsLogin,
			})
	}

	crd, manifests, err := oauth.Translate(identityProviders)

	require.NoError(t, err)

	CRD := crd.GenYAML()
	expectedYaml, _ := ioutil.ReadFile("testdata/expected-bulk-test-masterconfig-oauth.yaml")

	assert.Equal(t, len(manifests), 9)
	assert.Equal(t, expectedYaml, CRD)
}

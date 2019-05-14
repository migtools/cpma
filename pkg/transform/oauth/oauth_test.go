package oauth_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/oauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
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
	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	var masterV3 configv1.MasterConfig
	_, _, _ = serializer.Decode(content, nil, &masterV3)

	var identityProviders []oauth.IdentityProvider
	for _, identityProvider := range masterV3.OAuthConfig.IdentityProviders {
		providerJSON, _ := identityProvider.Provider.MarshalJSON()
		provider := oauth.Provider{}
		json.Unmarshal(providerJSON, &provider)

		identityProviders = append(identityProviders,
			oauth.IdentityProvider{
				provider.Kind,
				provider.APIVersion,
				identityProvider.MappingMethod,
				identityProvider.Name,
				identityProvider.Provider,
				provider.File,
				nil,
				nil,
				nil,
				identityProvider.UseAsChallenger,
				identityProvider.UseAsLogin,
			})
	}

	resCrd, _, err := oauth.Translate(identityProviders)
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

	var identityProviders []oauth.IdentityProvider
	for _, identityProvider := range masterV3.OAuthConfig.IdentityProviders {
		providerJSON, _ := identityProvider.Provider.MarshalJSON()
		provider := oauth.Provider{}
		json.Unmarshal(providerJSON, &provider)

		identityProviders = append(identityProviders,
			oauth.IdentityProvider{
				provider.Kind,
				provider.APIVersion,
				identityProvider.MappingMethod,
				identityProvider.Name,
				identityProvider.Provider,
				provider.File,
				nil,
				nil,
				nil,
				identityProvider.UseAsChallenger,
				identityProvider.UseAsLogin,
			})
	}

	crd, manifests, err := oauth.Translate(identityProviders)

	require.NoError(t, err)

	CRD := crd.GenYAML()
	expectedYaml, _ := ioutil.ReadFile("testdata/expected-bulk-test-masterconfig-oauth.yaml")

	assert.Equal(t, 10, len(manifests))
	assert.Equal(t, expectedYaml, CRD)
}

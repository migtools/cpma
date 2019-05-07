package ocp4

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/fusor/cpma/pkg/ocp4/oauth"
)

var _GetFile = oauth.GetFile

func mockGetFile(a, b, c string) []byte {
	return []byte("This is test file content")
}

func TestClusterTranslate(t *testing.T) {
	defer func() { oauth.GetFile = _GetFile }()
	oauth.GetFile = mockGetFile

	masterV4 := Master{}
	file := "../testdata/common-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)

	masterV3 := ocp3.Master{}
	masterV3.Decode(content)
	masterV4.Translate(masterV3.Config)

	assert.Equal(t, "cluster", masterV4.OAuth.Metadata.Name)
	assert.Equal(t, 2, len(masterV4.OAuth.Spec.IdentityProviders))

	assert.Equal(t, 2, len(masterV4.Secrets))
	assert.Equal(t, "htpasswd_auth-secret", masterV4.Secrets[0].Metadata.Name)
	assert.Equal(t, "github123456789-secret", masterV4.Secrets[1].Metadata.Name)
}

func TestClusterGenYaml(t *testing.T) {
	defer func() { oauth.GetFile = _GetFile }()
	oauth.GetFile = mockGetFile

	masterV4 := Master{}
	file := "../testdata/common-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)

	masterV3 := ocp3.Master{}
	masterV3.Decode(content)
	masterV4.Translate(masterV3.Config)
	manifests := masterV4.GenYAML()

	// Test manifest names
	assert.Equal(t, "100_CPMA-cluster-config-oauth.yaml", manifests[0].Name)
	assert.Equal(t, "100_CPMA-cluster-config-secret-htpasswd_auth-secret.yaml", manifests[1].Name)
	assert.Equal(t, "100_CPMA-cluster-config-secret-github123456789-secret.yaml", manifests[2].Name)
	assert.Equal(t, "100_CPMA-cluster-config-sdn.yaml", manifests[3].Name)

	// Test Oauth CR contents
	expectedOauthCR, _ := ioutil.ReadFile("testdata/expected-test-oauth-master.yaml")
	assert.Equal(t, expectedOauthCR, manifests[0].CRD)

	// Test secrets contents
	expectedSecretHtpasswd, _ := ioutil.ReadFile("testdata/expected-test-secret-httpasswd.yaml")
	expectedSecretGitHub, _ := ioutil.ReadFile("testdata/expected-test-secret-github.yaml")
	assert.Equal(t, expectedSecretHtpasswd, manifests[1].CRD)
	assert.Equal(t, expectedSecretGitHub, manifests[2].CRD)

	// Test network CR contents
	expectedNetworkCR, _ := ioutil.ReadFile("testdata/expected-test-network-cr-master.yaml")
	assert.Equal(t, expectedNetworkCR, manifests[3].CRD)
}

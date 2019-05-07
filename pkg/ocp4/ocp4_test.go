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

	// Test Oauth CR contents
	expectedOauthCR := `apiVersion: config.openshift.io/v1
kind: OAuth
metadata:
  name: cluster
  namespace: openshift-config
spec:
  identityProviders:
  - name: htpasswd_auth
    challenge: true
    login: true
    mappingMethod: claim
    type: HTPasswd
    htpasswd:
      fileData:
        name: htpasswd_auth-secret
  - name: github123456789
    challenge: false
    login: true
    mappingMethod: claim
    type: GitHub
    github:
      hostname: test.example.com
      ca:
        name: github.crt
      clientID: 2d85ea3f45d6777bffd7
      clientSecret:
        name: github123456789-secret
      organizations:
      - myorganization1
      - myorganization2
      teams:
      - myorganization1/team-a
      - myorganization2/team-b
`
	assert.Equal(t, expectedOauthCR, string(manifests[0].CRD))

	// Test secrets contents
	expectedSecretHtpasswd := `apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: htpasswd_auth-secret
  namespace: openshift-config
data:
  htpasswd: VGhpcyBpcyB0ZXN0IGZpbGUgY29udGVudA==
`

	expectedSecretGitHub := `apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: github123456789-secret
  namespace: openshift-config
data:
  clientSecret: ZTE2YTU5YWQzM2Q3YzI5ZmQ0MzU0ZjQ2MDU5ZjA5NTBjNjA5YTdlYQ==
`

	assert.Equal(t, expectedSecretHtpasswd, string(manifests[1].CRD))
	assert.Equal(t, expectedSecretGitHub, string(manifests[2].CRD))
}

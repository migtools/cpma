package ocp4

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fusor/cpma/ocp3"
)

func TestClusterTranslateMaster(t *testing.T) {
	parsedTestConfig := ocp3.ParseMaster("../test/common-test-master-config.yaml")

	clusterOCP4 := Cluster{}
	clusterOCP4.Master.TranslateMaster(parsedTestConfig)

	assert.Equal(t, "cluster", clusterOCP4.Master.OAuth.Metadata.Name)
	assert.Equal(t, 2, len(clusterOCP4.Master.OAuth.Spec.IdentityProviders))

	assert.Equal(t, 2, len(clusterOCP4.Master.Secrets))
	assert.Equal(t, "htpasswd_auth-secret", clusterOCP4.Master.Secrets[0].Metadata.Name)
	assert.Equal(t, "github123456789-secret", clusterOCP4.Master.Secrets[1].Metadata.Name)
}

func TestClusterMasterGenYaml(t *testing.T) {
	parsedTestConfig := ocp3.ParseMaster("../test/common-test-master-config.yaml")
	clusterV4 := Cluster{}
	clusterV4.Master.TranslateMaster(parsedTestConfig)
	manifests := clusterV4.GenYAML()

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
  htpasswd: VGhpcyBpcyBwcmV0ZW5kIGNvbnRlbnQ=
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

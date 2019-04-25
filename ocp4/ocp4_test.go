package ocp4

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fusor/cpma/ocp3"
)

func TestClusterTranslate(t *testing.T) {
	testConfig := ocp3.Config{
		Masterf: "../test/common-test-master-config.yaml",
	}
	parsedTestConfig := testConfig.ParseMaster()
	clusterV4 := Cluster{}
	clusterV4.Translate(parsedTestConfig)

	assert.Equal(t, "cluster", clusterV4.Master.OAuth.MetaData.Name)
	assert.Equal(t, 2, len(clusterV4.Master.OAuth.Spec.IdentityProviders))

	assert.Equal(t, 2, len(clusterV4.Master.Secrets))
	assert.Equal(t, "htpasswd_auth-secret", clusterV4.Master.Secrets[0].MetaData.Name)
	assert.Equal(t, "github123456789-secret", clusterV4.Master.Secrets[1].MetaData.Name)
}

func TestClusterGenYaml(t *testing.T) {
	testConfig := ocp3.Config{
		Masterf: "../test/common-test-master-config.yaml",
	}
	parsedTestConfig := testConfig.ParseMaster()
	clusterV4 := Cluster{}
	clusterV4.Translate(parsedTestConfig)
	manifests := clusterV4.GenYAML()

	// Test manifest names
	assert.Equal(t, "CPMA-cluster-config-oauth.yaml", manifests[0].Name)
	assert.Equal(t, "CPMA-cluster-config-secret-htpasswd_auth-secret.yaml", manifests[1].Name)
	assert.Equal(t, "CPMA-cluster-config-secret-github123456789-secret.yaml", manifests[2].Name)

	// Test Oauth CR contents
	expectedOauthCR := `apiVersion: config.openshift.io/v1
kind: OAuth
metaData:
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
metaData:
  name: htpasswd_auth-secret
  namespace: openshift-config
data:
  htpasswd: VGhpcyBpcyBwcmV0ZW5kIGNvbnRlbnQ=
`

	expectedSecretGitHub := `apiVersion: v1
kind: Secret
type: Opaque
metaData:
  name: github123456789-secret
  namespace: openshift-config
data:
  clientSecret: e16a59ad33d7c29fd4354f46059f0950c609a7ea
`

	assert.Equal(t, expectedSecretHtpasswd, string(manifests[1].CRD))
	assert.Equal(t, expectedSecretGitHub, string(manifests[2].CRD))
}

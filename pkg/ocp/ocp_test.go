package ocp

import (
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/env"
	"github.com/fusor/cpma/internal/io"
	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/fusor/cpma/pkg/ocp4/oauth"
	"github.com/stretchr/testify/assert"
)

var _GetFile = io.GetFile

func mockGetFile(a, b, c string) []byte {
	return []byte("This is test file content")
}

func TestAddOAuthConfig(t *testing.T) {
	// Init config with default master config paths
	ocpMaster := OAuthConfig{}
	ocpMaster.Add("example.com")

	assert.Equal(t, OAuthConfig{
		ConfigFile: ConfigFile{
			Hostname: "example.com",
			Path:     "/etc/origin/master/master-config.yaml",
		},
	}, ocpMaster)

	// Init config with different master config path
	env.Config().Set("MasterConfigFile", "/test/path/master.yml")
	ocpMaster = OAuthConfig{}
	ocpMaster.Add("example.com")

	assert.Equal(t, OAuthConfig{
		ConfigFile: ConfigFile{
			Hostname: "example.com",
			Path:     "/test/path/master.yml",
		},
	}, ocpMaster)
	env.Config().Set("MasterConfigFile", "/etc/origin/master/master-config.yaml")
}

func TestAddSDNConfig(t *testing.T) {
	// Init config with default node config paths
	ocpMaster := SDNConfig{}
	ocpMaster.Add("example.com")

	assert.Equal(t, SDNConfig{
		ConfigFile: ConfigFile{
			Hostname: "example.com",
			Path:     "/etc/origin/master/master-config.yaml",
		},
	}, ocpMaster)

	// Init config with different node config paths
	env.Config().Set("MasterConfigFile", "/test/path/another.yml")
	ocpMaster = SDNConfig{}
	ocpMaster.Add("example.com")

	assert.Equal(t, SDNConfig{
		ConfigFile: ConfigFile{
			Hostname: "example.com",
			Path:     "/test/path/another.yml",
		},
	}, ocpMaster)
}

func TestTransformOAuth(t *testing.T) {
	defer func() { io.GetFile = _GetFile }()
	oauth.GetFile = mockGetFile

	file := "../testdata/common-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)
	masterV3 := ocp3.MasterDecode(content)
	oauthConfig, secrets, _ := oauth.Transform(masterV3.OAuthConfig)

	assert.Equal(t, "cluster", oauthConfig.Metadata.Name)
	assert.Equal(t, 2, len(oauthConfig.Spec.IdentityProviders))

	assert.Equal(t, 2, len(secrets))
	assert.Equal(t, "htpasswd_auth-secret", secrets[0].Metadata.Name)
	assert.Equal(t, "github123456789-secret", secrets[1].Metadata.Name)
}

func TestGenYamlOAuth(t *testing.T) {
	defer func() { io.GetFile = _GetFile }()
	oauth.GetFile = mockGetFile

	file := "../testdata/common-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)
	masterV3 := ocp3.MasterDecode(content)

	oauthConfig := OAuthConfig{}
	oauthConfig.OCP3.OAuthConfig = masterV3.OAuthConfig
	oauthConfig.Transform()

	manifests := oauthConfig.GenYAML()

	// Test number of manifests
	assert.Equal(t, len(manifests), 3)

	// Test manifest names
	assert.Equal(t, "100_CPMA-cluster-config-oauth.yaml", manifests[0].Name)
	assert.Equal(t, "100_CPMA-cluster-config-secret-htpasswd_auth-secret.yaml", manifests[1].Name)
	assert.Equal(t, "100_CPMA-cluster-config-secret-github123456789-secret.yaml", manifests[2].Name)

	// Test Oauth CR contents
	expectedOauthCR, _ := ioutil.ReadFile("testdata/expected-test-oauth-master.yaml")
	assert.Equal(t, expectedOauthCR, manifests[0].CRD)

	// Test secrets contents
	expectedSecretHtpasswd, _ := ioutil.ReadFile("testdata/expected-test-secret-httpasswd.yaml")
	expectedSecretGitHub, _ := ioutil.ReadFile("testdata/expected-test-secret-github.yaml")
	assert.Equal(t, expectedSecretHtpasswd, manifests[1].CRD)
	assert.Equal(t, expectedSecretGitHub, manifests[2].CRD)
}

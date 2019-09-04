package cmd_test

import (
	"os"
	"testing"

	_ "github.com/fusor/cpma/cmd"
	"github.com/fusor/cpma/pkg/env"
	"github.com/stretchr/testify/assert"
)

func TestInitDefaults(t *testing.T) {
	assert.Equal(t, "", env.Config().GetString("ConfigSource"))
	assert.Equal(t, "", env.Config().GetString("ClusterName"))
	assert.Equal(t, "", env.Config().GetString("CRIOConfigFile"))
	assert.Equal(t, false, env.Config().Get("Debug"))
	assert.Equal(t, "", env.Config().GetString("ETCDConfigfile"))
	assert.Equal(t, "", env.Config().GetString("Hostname"))
	assert.Equal(t, false, env.Config().Get("InsecureHostKey"))
	assert.Equal(t, "", env.Config().GetString("NodeConfigFile"))
	assert.Equal(t, "", env.Config().GetString("MasterConfigFile"))
	assert.Equal(t, true, env.Config().Get("Manifests"))
	assert.Equal(t, "", env.Config().GetString("RegistriesConfigFile"))
	assert.Equal(t, true, env.Config().Get("Reporting"))
	assert.Equal(t, "", env.Config().GetString("SSHPrivateKey"))
	assert.Equal(t, "", env.Config().GetString("SSHLogin"))
	assert.Equal(t, "0", env.Config().GetString("SSHPort"))
	assert.Equal(t, false, env.Config().Get("Silent"))
	assert.Equal(t, "", env.Config().GetString("WorkDIr"))
}

func TestInitSetValues(t *testing.T) {
	defer func() {
		os.Unsetenv("CPMA_CONFIGSOURCE")
		os.Unsetenv("CPMA_CLUSTERNAME")
		os.Unsetenv("CPMA_CRIOCONFIGFILE")
		os.Unsetenv("CPMA_DEBUG")
		os.Unsetenv("CPMA_ETCDCONFIGFILE")
		os.Unsetenv("CPMA_HOSTNAME")
		os.Unsetenv("CPMA_INSECUREHOSTKEY")
		os.Unsetenv("CPMA_NODECONFIGFILE")
		os.Unsetenv("CPMA_MANIFESTS")
		os.Unsetenv("CPMA_MASTERCONFIGFILE")
		os.Unsetenv("CPMA_REGISTRIESCONFIGFILE")
		os.Unsetenv("CPMA_REPORTING")
		os.Unsetenv("CPMA_SSHPRIVATEKEY")
		os.Unsetenv("CPMA_SSHLOGIN")
		os.Unsetenv("CPMA_SSHPORT")
		os.Unsetenv("CPMA_SILENT")
		os.Unsetenv("CPMA_WORKDIR")
	}()

	os.Setenv("CPMA_CONFIGSOURCE", "remote")
	os.Setenv("CPMA_CLUSTERNAME", "cluster1-example-com")
	os.Setenv("CPMA_CRIOCONFIGFILE", "/tmp/crio.conf")
	os.Setenv("CPMA_DEBUG", "true")
	os.Setenv("CPMA_ETCDCONFIGFILE", "/tmp/etcd.conf")
	os.Setenv("CPMA_HOSTNAME", "www.example.com")
	os.Setenv("CPMA_INSECUREHOSTKEY", "true")
	os.Setenv("CPMA_NODECONFIGFILE", "/tmp/node-config.yaml")
	os.Setenv("CPMA_MANIFESTS", "false")
	os.Setenv("CPMA_MASTERCONFIGFILE", "/tmp/master-config.yaml")
	os.Setenv("CPMA_REGISTRIESCONFIGFILE", "/tmp/registries.conf")
	os.Setenv("CPMA_REPORTING", "false")
	os.Setenv("CPMA_SSHPRIVATEKEY", "/test/.ssh/test")
	os.Setenv("CPMA_SSHLOGIN", "testuser")
	os.Setenv("CPMA_SSHPORT", "8080")
	os.Setenv("CPMA_SILENT", "true")
	os.Setenv("CPMA_WORKDIR", "./testdir")
	env.InitConfig()

	assert.Equal(t, "remote", env.Config().GetString("ConfigSource"))
	assert.Equal(t, "cluster1-example-com", env.Config().GetString("ClusterName"))
	assert.Equal(t, "/tmp/crio.conf", env.Config().GetString("CRIOConfigFile"))
	assert.Equal(t, "/tmp/etcd.conf", env.Config().GetString("ETCDConfigfile"))
	assert.Equal(t, true, env.Config().GetBool("Debug"))
	assert.Equal(t, "www.example.com", env.Config().GetString("Hostname"))
	assert.Equal(t, true, env.Config().GetBool("InsecureHostKey"))
	assert.Equal(t, "/tmp/node-config.yaml", env.Config().GetString("NodeConfigFile"))
	assert.Equal(t, "/tmp/master-config.yaml", env.Config().GetString("MasterConfigFile"))
	assert.Equal(t, false, env.Config().GetBool("Manifests"))
	assert.Equal(t, "/tmp/registries.conf", env.Config().GetString("RegistriesConfigFile"))
	assert.Equal(t, false, env.Config().GetBool("Reporting"))
	assert.Equal(t, "/test/.ssh/test", env.Config().GetString("SSHPrivateKey"))
	assert.Equal(t, "testuser", env.Config().GetString("SSHLogin"))
	assert.Equal(t, "8080", env.Config().GetString("SSHPort"))
	assert.Equal(t, true, env.Config().GetBool("Silent"))
	assert.Equal(t, "./testdir", env.Config().GetString("WorkDIr"))
}

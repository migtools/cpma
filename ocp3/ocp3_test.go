package ocp3

import (
	"testing"

	configv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/stretchr/testify/assert"
)

func TestConfigParseMaster(t *testing.T) {
	testConfig := Config{
		Masterf: "../test/common-test-master-config.yaml",
	}

	expectedMasterConfig := configv1.MasterConfig{
		AuthConfig: configv1.MasterAuthConfig{
			RequestHeader: &configv1.RequestHeaderAuthenticationOptions{
				ClientCA: "front-proxy-ca.crt",
			},
		},
		EtcdClientInfo: configv1.EtcdConnectionInfo{
			CA:   "master.etcd-ca.crt",
			URLs: []string{"https://master-0.gildub2.lab.pnq2.cee.redhat.com:2379"},
		},
		OAuthConfig: &configv1.OAuthConfig{
			MasterURL: "https://openshift.internal.gildub2.lab.pnq2.cee.redhat.com:443",
			IdentityProviders: []configv1.IdentityProvider{
				configv1.IdentityProvider{
					Name: "htpasswd_auth",
				},
				configv1.IdentityProvider{
					Name: "github123456789",
				},
			},
		},
	}

	resMasterConfig := testConfig.ParseMaster()

	assert.Equal(t, expectedMasterConfig.AuthConfig.RequestHeader.ClientCA, resMasterConfig.AuthConfig.RequestHeader.ClientCA)
	assert.Equal(t, expectedMasterConfig.EtcdClientInfo.CA, resMasterConfig.EtcdClientInfo.CA)
	assert.Equal(t, expectedMasterConfig.EtcdClientInfo.URLs, resMasterConfig.EtcdClientInfo.URLs)
	assert.Equal(t, expectedMasterConfig.OAuthConfig.MasterURL, resMasterConfig.OAuthConfig.MasterURL)
	assert.Equal(t, expectedMasterConfig.OAuthConfig.IdentityProviders[0].Name, resMasterConfig.OAuthConfig.IdentityProviders[0].Name)
	assert.Equal(t, expectedMasterConfig.OAuthConfig.IdentityProviders[1].Name, resMasterConfig.OAuthConfig.IdentityProviders[1].Name)
}

func TestNewConfig(t *testing.T) {
	config := New()

	assert.Equal(t, &Config{
		Masterf: "/etc/origin/master/master-config.yaml",
		Nodef:   "/etc/origin/node/node-config.yaml",
	}, config)
}

package transform

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

func TestTransformMasterConfig(t *testing.T) {
	file := "testdata/network-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)

	networkCR := SDNTranslate(content)

	// Check if network CR was translated correctly
	assert.Equal(t, networkCR.APIVersion, "operator.openshift.io/v1")
	assert.Equal(t, networkCR.Kind, "Network")
	assert.Equal(t, networkCR.Spec.ClusterNetworks[0].CIDR, "10.128.0.0/14")
	assert.Equal(t, networkCR.Spec.ClusterNetworks[0].HostPrefix, uint32(9))
	assert.Equal(t, networkCR.Spec.ServiceNetwork, "172.30.0.0/16")
	assert.Equal(t, networkCR.Spec.DefaultNetwork.Type, "OpenShiftSDN")
	assert.Equal(t, networkCR.Spec.DefaultNetwork.OpenshiftSDNConfig.Mode, "Subnet")
}

func TestSelectNetworkPlugin(t *testing.T) {
	resPluginName, err := SelectNetworkPlugin("redhat/openshift-ovs-multitenant")
	require.NoError(t, err)
	assert.Equal(t, "Multitenant", resPluginName)

	resPluginName, err = SelectNetworkPlugin("redhat/openshift-ovs-networkpolicy")
	require.NoError(t, err)
	assert.Equal(t, "NetworkPolicy", resPluginName)

	resPluginName, err = SelectNetworkPlugin("redhat/openshift-ovs-subnet")
	require.NoError(t, err)
	assert.Equal(t, "Subnet", resPluginName)

	_, err = SelectNetworkPlugin("123")
	expectedErr := errors.New("Network plugin not supported")
	assert.Error(t, expectedErr, err)
}

func TestTransformClusterNetworks(t *testing.T) {
	var clusterNeworkEntries []configv1.ClusterNetworkEntry
	clusterNetwork1 := configv1.ClusterNetworkEntry{CIDR: "10.128.0.0/14", HostSubnetLength: uint32(9)}
	clusterNetwork2 := configv1.ClusterNetworkEntry{CIDR: "10.127.0.0/14", HostSubnetLength: uint32(10)}
	clusterNeworkEntries = append(clusterNeworkEntries, clusterNetwork1, clusterNetwork2)

	translatedClusterNetworks := TranslateClusterNetworks(clusterNeworkEntries)
	assert.Equal(t, translatedClusterNetworks[0].CIDR, "10.128.0.0/14")
	assert.Equal(t, translatedClusterNetworks[0].HostPrefix, uint32(9))
	assert.Equal(t, translatedClusterNetworks[1].CIDR, "10.127.0.0/14")
	assert.Equal(t, translatedClusterNetworks[1].HostPrefix, uint32(10))
}

func TestGenYAML(t *testing.T) {
	file := "testdata/network-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)

	networkCR := SDNTranslate(content)
	networkCRYAML := GenYAML(networkCR)

	expectedYaml, _ := ioutil.ReadFile("testdata/expected-network-cr-master.yaml")
	assert.Equal(t, expectedYaml, networkCRYAML)
}

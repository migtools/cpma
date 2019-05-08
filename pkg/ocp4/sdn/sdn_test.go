package sdn

import (
	"errors"
	"io/ioutil"
	"testing"

	configv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fusor/cpma/pkg/ocp3"
)

func TestTransformMasterConfig(t *testing.T) {
	file := "testdata/network-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)

	masterV3 := ocp3.Master{}
	masterV3.Decode(content)

	resNetworkCR := Transform(masterV3.Config.NetworkConfig)

	// Check if network CR was translated correctly
	assert.Equal(t, resNetworkCR.APIVersion, "operator.openshift.io/v1")
	assert.Equal(t, resNetworkCR.Kind, "Network")
	assert.Equal(t, resNetworkCR.Spec.ClusterNetworks[0].CIDR, "10.128.0.0/14")
	assert.Equal(t, resNetworkCR.Spec.ClusterNetworks[0].HostPrefix, uint32(9))
	assert.Equal(t, resNetworkCR.Spec.ServiceNetwork, "172.30.0.0/16")
	assert.Equal(t, resNetworkCR.Spec.DefaultNetwork.Type, "OpenShiftSDN")
	assert.Equal(t, resNetworkCR.Spec.DefaultNetwork.OpenshiftSDNConfig.Mode, "Subnet")
}

func TestSelectNetworkPlugin(t *testing.T) {
	resPluginName, err := selectNetworkPlugin("redhat/openshift-ovs-multitenant")
	require.NoError(t, err)
	assert.Equal(t, "Multitenant", resPluginName)

	resPluginName, err = selectNetworkPlugin("redhat/openshift-ovs-networkpolicy")
	require.NoError(t, err)
	assert.Equal(t, "NetworkPolicy", resPluginName)

	resPluginName, err = selectNetworkPlugin("redhat/openshift-ovs-subnet")
	require.NoError(t, err)
	assert.Equal(t, "Subnet", resPluginName)

	_, err = selectNetworkPlugin("123")
	expectedErr := errors.New("Network plugin not supported")
	assert.Error(t, expectedErr, err)
}

func TestTransformClusterNetworks(t *testing.T) {
	var clusterNeworkEntries []configv1.ClusterNetworkEntry
	clusterNetwork1 := configv1.ClusterNetworkEntry{CIDR: "10.128.0.0/14", HostSubnetLength: uint32(9)}
	clusterNetwork2 := configv1.ClusterNetworkEntry{CIDR: "10.127.0.0/14", HostSubnetLength: uint32(10)}
	clusterNeworkEntries = append(clusterNeworkEntries, clusterNetwork1, clusterNetwork2)

	translatedClusterNetworks := translateClusterNetworks(clusterNeworkEntries)
	assert.Equal(t, translatedClusterNetworks[0].CIDR, "10.128.0.0/14")
	assert.Equal(t, translatedClusterNetworks[0].HostPrefix, uint32(9))
	assert.Equal(t, translatedClusterNetworks[1].CIDR, "10.127.0.0/14")
	assert.Equal(t, translatedClusterNetworks[1].HostPrefix, uint32(10))
}

func TestGenYAML(t *testing.T) {
	file := "testdata/network-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)

	masterV3 := ocp3.Master{}
	masterV3.Decode(content)

	networkCR := Transform(masterV3.Config.NetworkConfig)

	networkCRYAML := networkCR.GenYAML()

	expectedYaml, _ := ioutil.ReadFile("testdata/expected-network-cr-master.yaml")
	assert.Equal(t, expectedYaml, networkCRYAML)
}

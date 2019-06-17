package sdn

import (
	"errors"
	"net"

	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	configv1 "github.com/openshift/api/operator/v1"
)

const (
	apiVersion         = "operator.openshift.io/v1"
	kind               = "Network"
	defaultNetworkType = "OpenShiftSDN"
	name               = "cluster"
)

// Translate is called by Transform to do the majority of the work in converting data
func Translate(masterConfig legacyconfigv1.MasterConfig) (*configv1.Network, error) {
	networkConfig := masterConfig.NetworkConfig
	var networkCR configv1.Network

	networkCR.APIVersion = apiVersion
	networkCR.Kind = kind
	networkCR.Name = name
	networkCR.Spec.ServiceNetwork = []string{networkConfig.ServiceNetworkCIDR}
	networkCR.Spec.DefaultNetwork.Type = defaultNetworkType

	// Translate CIDRs and adress size for each node
	translatedClusterNetworks := TranslateClusterNetworks(networkConfig.ClusterNetworks)
	networkCR.Spec.ClusterNetwork = translatedClusterNetworks

	// Translate network plugin name
	selectedNetworkPlugin, err := SelectNetworkPlugin(networkConfig.NetworkPluginName)
	if err != nil {
		return nil, err
	}

	openshiftSDNConfig := &configv1.OpenShiftSDNConfig{
		Mode: configv1.SDNMode(selectedNetworkPlugin),
	}
	networkCR.Spec.DefaultNetwork.OpenShiftSDNConfig = openshiftSDNConfig

	return &networkCR, nil
}

// TranslateClusterNetworks converts Cluster Networks from OCP3 to OCP4
func TranslateClusterNetworks(clusterNeworkEntries []legacyconfigv1.ClusterNetworkEntry) []configv1.ClusterNetworkEntry {
	var translatedClusterNetworks []configv1.ClusterNetworkEntry

	for _, networkConfig := range clusterNeworkEntries {
		var translatedClusterNetwork configv1.ClusterNetworkEntry

		translatedClusterNetwork.CIDR = networkConfig.CIDR
		// host prefix is missing in OCP3 config, default is 23
		translatedClusterNetwork.HostPrefix = 23

		translatedClusterNetworks = append(translatedClusterNetworks, translatedClusterNetwork)
	}

	return translatedClusterNetworks
}

// SelectNetworkPlugin selects the correct plugin for networks
func SelectNetworkPlugin(pluginName string) (string, error) {
	var selectedName string

	switch pluginName {
	case "redhat/openshift-ovs-multitenant":
		selectedName = "Multitenant"
	case "redhat/openshift-ovs-networkpolicy":
		selectedName = "NetworkPolicy"
	case "redhat/openshift-ovs-subnet":
		selectedName = "Subnet"
	default:
		return "", errors.New("Network plugin not supported")
	}

	return selectedName, nil
}

// Validate validate SDN config
func Validate(masterConfig legacyconfigv1.MasterConfig) error {
	networkConfig := masterConfig.NetworkConfig

	if len(networkConfig.ServiceNetworkCIDR) == 0 {
		return errors.New("Service network CIDR can't be empty")
	}

	if _, _, err := net.ParseCIDR(networkConfig.ServiceNetworkCIDR); err != nil {
		return errors.New("Not valid service network CIDR")
	}

	if len(networkConfig.ClusterNetworks) == 0 {
		return errors.New("Cluster network must have at least 1 entry")
	}

	for _, cnet := range networkConfig.ClusterNetworks {
		if len(cnet.CIDR) == 0 {
			return errors.New("Cluster network CIDR can't be empty")
		}

		if _, _, err := net.ParseCIDR(cnet.CIDR); err != nil {
			return errors.New("Not valid cluster network CIDR")
		}
	}

	if len(networkConfig.NetworkPluginName) == 0 {
		return errors.New("Plugin name can't be empty")
	}

	return nil
}

package sdn

import (
	"errors"
	"net"

	configv1 "github.com/openshift/api/legacyconfig/v1"
	"gopkg.in/yaml.v2"
)

// NetworkCR describes Network CR for OCP4
type NetworkCR struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Spec       Spec   `yaml:"spec"`
}

// Spec is a SDN specific spec
type Spec struct {
	ClusterNetworks []ClusterNetwork `yaml:"clusterNetwork"`
	ServiceNetwork  string           `yaml:"serviceNetwork"`
	DefaultNetwork  `yaml:"defaultNetwork"`
}

// ClusterNetwork contains CIDR and address size to assign to each node
type ClusterNetwork struct {
	CIDR       string `yaml:"cidr"`
	HostPrefix int    `yaml:"hostPrefix"`
}

// DefaultNetwork containts network type and SDN plugin name
type DefaultNetwork struct {
	Type               string             `yaml:"type"`
	OpenshiftSDNConfig OpenshiftSDNConfig `yaml:"openshiftSDNConfig"`
}

// OpenshiftSDNConfig is the Openshift SDN Configured Mode
type OpenshiftSDNConfig struct {
	Mode string `yaml:"mode"`
}

const (
	apiVersion         = "operator.openshift.io/v1"
	kind               = "Network"
	defaultNetworkType = "OpenShiftSDN"
)

// Translate is called by Transform to do the majority of the work in converting data
func Translate(masterConfig configv1.MasterConfig) (NetworkCR, error) {
	networkConfig := masterConfig.NetworkConfig
	var networkCR NetworkCR

	networkCR.APIVersion = apiVersion
	networkCR.Kind = kind
	networkCR.Spec.ServiceNetwork = networkConfig.ServiceNetworkCIDR
	networkCR.Spec.DefaultNetwork.Type = defaultNetworkType

	// Translate CIDRs and adress size for each node
	translatedClusterNetworks := TranslateClusterNetworks(networkConfig.ClusterNetworks)
	networkCR.Spec.ClusterNetworks = translatedClusterNetworks

	// Translate network plugin name
	selectedNetworkPlugin, err := SelectNetworkPlugin(networkConfig.NetworkPluginName)
	if err != nil {
		return networkCR, err
	}
	networkCR.Spec.DefaultNetwork.OpenshiftSDNConfig.Mode = selectedNetworkPlugin

	return networkCR, nil
}

// TranslateClusterNetworks converts Cluster Networks from OCP3 to OCP4
func TranslateClusterNetworks(clusterNeworkEntries []configv1.ClusterNetworkEntry) []ClusterNetwork {
	var translatedClusterNetworks []ClusterNetwork

	for _, networkConfig := range clusterNeworkEntries {
		var translatedClusterNetwork ClusterNetwork

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
		err := errors.New("Network plugin not supported")
		return "", err
	}

	return selectedName, nil
}

// Validate validate SDN config
func Validate(masterConfig configv1.MasterConfig) error {
	networkConfig := masterConfig.NetworkConfig

	if len(networkConfig.ServiceNetworkCIDR) == 0 {
		return errors.New("Service network CIDR can't be empty")
	}

	_, _, err := net.ParseCIDR(networkConfig.ServiceNetworkCIDR)
	if err != nil {
		return errors.New("Not valid service network CIDR")
	}

	if len(networkConfig.ClusterNetworks) == 0 {
		return errors.New("Cluster network must have at least 1 entry")
	}

	for _, cnet := range networkConfig.ClusterNetworks {
		if len(cnet.CIDR) == 0 {
			return errors.New("Cluster network CIDR can't be empty")
		}

		_, _, err := net.ParseCIDR(cnet.CIDR)
		if err != nil {
			return errors.New("Not valid cluster network CIDR")
		}
	}

	if len(networkConfig.NetworkPluginName) == 0 {
		return errors.New("Plugin name can't be empty")
	}

	return nil
}

// GenYAML returns a YAML of the SDN CR
func GenYAML(networkCR NetworkCR) ([]byte, error) {
	yamlBytes, err := yaml.Marshal(networkCR)
	if err != nil {
		return nil, err
	}

	return yamlBytes, nil
}

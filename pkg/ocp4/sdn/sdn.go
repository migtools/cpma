package sdn

import (
	"errors"

	configv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// NetworkCR describes Network CR for OCP4
type NetworkCR struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Spec       struct {
		ClusterNetworks []ClusterNetwork `yaml:"clusterNetwork"`
		ServiceNetwork  string           `yaml:"serviceNetwork"`
		DefaultNetwork  `yaml:"defaultNetwork"`
	} `yaml:"spec"`
}

// ClusterNetwork contains CIDR and address size to assign to each node
type ClusterNetwork struct {
	CIDR       string `yaml:"cidr"`
	HostPrefix uint32 `yaml:"hostPrefix"`
}

// DefaultNetwork containts network type and SDN plugin name
type DefaultNetwork struct {
	Type               string `yaml:"type"`
	OpenshiftSDNConfig struct {
		Mode string `yaml:"mode"`
	} `yaml:"openshiftSDNConfig"`
}

const (
	apiVersion         = "operator.openshift.io/v1"
	kind               = "Network"
	defaultNetworkType = "OpenShiftSDN"
)

// Transform converts OCPv3 SDN to OCPv4 SDN Custom Resources
func Transform(networkConfig configv1.MasterNetworkConfig) *NetworkCR {
	var networkCR NetworkCR

	networkCR.APIVersion = apiVersion
	networkCR.Kind = kind
	networkCR.Spec.ServiceNetwork = networkConfig.ServiceNetworkCIDR
	networkCR.Spec.DefaultNetwork.Type = defaultNetworkType

	// Transform CIDRs and adress size for each node
	translatedClusterNetworks := translateClusterNetworks(networkConfig.ClusterNetworks)
	networkCR.Spec.ClusterNetworks = translatedClusterNetworks

	// Transform network plugin name
	selectedNetworkPlugin, err := selectNetworkPlugin(networkConfig.NetworkPluginName)
	if err != nil {
		logrus.Fatal(err)
	}
	networkCR.Spec.DefaultNetwork.OpenshiftSDNConfig.Mode = selectedNetworkPlugin

	return &networkCR
}

func translateClusterNetworks(clusterNeworkEntries []configv1.ClusterNetworkEntry) []ClusterNetwork {
	translatedClusterNetworks := make([]ClusterNetwork, 0)

	for _, networkConfig := range clusterNeworkEntries {
		var translatedClusterNetwork ClusterNetwork

		translatedClusterNetwork.CIDR = networkConfig.CIDR
		translatedClusterNetwork.HostPrefix = networkConfig.HostSubnetLength

		translatedClusterNetworks = append(translatedClusterNetworks, translatedClusterNetwork)
	}

	return translatedClusterNetworks
}

func selectNetworkPlugin(pluginName string) (string, error) {
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

// GenYAML returns a YAML of the OAuthCRD
func (networkCR *NetworkCR) GenYAML() []byte {
	yamlBytes, err := yaml.Marshal(&networkCR)
	if err != nil {
		logrus.Fatal(err)
	}

	return yamlBytes
}

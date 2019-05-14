package transform

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

// SDNExtraction is an SDN specific extraction
type SDNExtraction struct {
	configv1.MasterConfig
}

// NetworkCR describes Network CR for OCP4
type NetworkCR struct {
	APIVersion string  `yaml:"apiVersion"`
	Kind       string  `yaml:"kind"`
	Spec       SDNSpec `yaml:"spec"`
}

// SDNSpec is a SDN specific spec
type SDNSpec struct {
	ClusterNetworks []ClusterNetwork `yaml:"clusterNetwork"`
	ServiceNetwork  string           `yaml:"serviceNetwork"`
	DefaultNetwork  `yaml:"defaultNetwork"`
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

// SDNTransform is an SDN specific transform
type SDNTransform struct {
	Config *Config
}

const (
	apiVersion         = "operator.openshift.io/v1"
	kind               = "Network"
	defaultNetworkType = "OpenShiftSDN"
)

// Transform convers OCP3 data to configuration useful for OCP4
func (e SDNExtraction) Transform() (Output, error) {
	logrus.Info("SDNTransform::Transform")

	var manifests []Manifest

	networkCR := SDNTranslate(e.MasterConfig)
	networkCRYAML := GenYAML(networkCR)

	manifest := Manifest{Name: "100_CPMA-cluster-config-sdn.yaml", CRD: networkCRYAML}
	manifests = append(manifests, manifest)

	return ManifestOutput{
		Manifests: manifests,
	}, nil
}

// SDNTranslate is called by Transform to do the majority of the work in converting data
func SDNTranslate(masterConfig configv1.MasterConfig) NetworkCR {
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
		HandleError(err)
	}
	networkCR.Spec.DefaultNetwork.OpenshiftSDNConfig.Mode = selectedNetworkPlugin

	return networkCR
}

// Extract collects SDN configuration information from an OCP3 cluster
func (e SDNTransform) Extract() Extraction {
	logrus.Info("SDNTransform::Extract")
	content := e.Config.Fetch(e.Config.MasterConfigFile)
	var extraction SDNExtraction

	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	_, _, err := serializer.Decode(content, nil, &extraction.MasterConfig)
	if err != nil {
		HandleError(err)
	}

	return extraction
}

// Validate the data extracted from the OCP3 cluster
func (e SDNExtraction) Validate() error {
	logrus.Warn("SDN Transform Validation Not Implmeneted")
	return nil // Simulate fine
}

// TranslateClusterNetworks converts Cluster Networks from OCP3 to OCP4
func TranslateClusterNetworks(clusterNeworkEntries []configv1.ClusterNetworkEntry) []ClusterNetwork {
	var translatedClusterNetworks []ClusterNetwork

	for _, networkConfig := range clusterNeworkEntries {
		var translatedClusterNetwork ClusterNetwork

		translatedClusterNetwork.CIDR = networkConfig.CIDR
		translatedClusterNetwork.HostPrefix = networkConfig.HostSubnetLength

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

// GenYAML returns a YAML of the OAuthCRD
func GenYAML(networkCR NetworkCR) []byte {
	yamlBytes, err := yaml.Marshal(networkCR)
	if err != nil {
		logrus.Fatal(err)
	}

	return yamlBytes
}

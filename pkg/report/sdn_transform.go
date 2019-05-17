package report

import (
	"errors"

	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/etl"
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
	DefaultNetwork `yaml:"defaultNetwork"`
}

// ClusterNetwork contains CIDR and address size to assign to each node
type ClusterNetwork struct {
	CIDR       string `yaml:"cidr"`
	HostPrefix int    `yaml:"hostPrefix"`
}

// DefaultNetwork containts network type and SDN plugin name
type DefaultNetwork struct {
	Type               string `yaml:"type"`
	OpenshiftSDNConfig struct {
		Mode string `yaml:"mode"`
	} `yaml:"openshiftSDNConfig"`
}

// SDNInstallConfig is a yaml snippet of install config
type SDNInstallConfig struct {
	ClusterNetworks []ClusterNetwork `yaml:"clusterNetwork"`
	ServiceNetwork  []string         `yaml:"serviceNetwork"`
}

// SDNReport is an SDN specific transform
type SDNReport struct {
	Config *etl.Config
}

const (
	apiVersion         = "operator.openshift.io/v1"
	kind               = "Network"
	defaultNetworkType = "OpenShiftSDN"
	readme             = `Migrating IP adresses in a Day 1 operation, it's required to copy values from sdn-install-config.yaml snippet and place them under "networking" section in instal-config.yaml.
In order to migrate network plugin, cluster-config-sdn.yaml should be placed in manifests directory.
`
)

// Transform convers OCP3 data to configuration useful for OCP4
func (e SDNExtraction) Transform() (etl.Output, error) {
	logrus.Info("SDNReport::Transform")

	var reports []etl.Data

	networkCR, sdnInstallConfig, err := SDNTranslate(e.MasterConfig)
	if err != nil {
		return nil, err
	}

	networkCRYAML, err := GenYAML(networkCR)
	if err != nil {
		return nil, err
	}

	report := etl.Data{Name: "SDN-operator-config-sdn.yaml", Type: "reports", File: networkCRYAML}
	reports = append(reports, report)

	sdnInstallConfigYAML, err := GenYAML(sdnInstallConfig)
	if err != nil {
		return nil, err
	}

	report = etl.Data{Name: "SDN-install-config.yaml", Type: "reports", File: sdnInstallConfigYAML}
	reports = append(reports, report)

	readmeByteSlice := []byte(readme)
	report = etl.Data{Name: "SDN-readme.md", Type: "reports", File: readmeByteSlice}
	reports = append(reports, report)

	return etl.DataOutput{
		DataList: reports,
	}, nil
}

// SDNTranslate is called by Transform to do the majority of the work in converting data
func SDNTranslate(masterConfig configv1.MasterConfig) (NetworkCR, SDNInstallConfig, error) {
	networkConfig := masterConfig.NetworkConfig
	var networkCR NetworkCR
	var sdnInstallConfig SDNInstallConfig

	networkCR.APIVersion = apiVersion
	networkCR.Kind = kind
	networkCR.Spec.DefaultNetwork.Type = defaultNetworkType

	// Translate CIDRs and adress size for each node
	translatedClusterNetworks := TranslateClusterNetworks(networkConfig.ClusterNetworks)

	// Translate network plugin name
	selectedNetworkPlugin, err := SelectNetworkPlugin(networkConfig.NetworkPluginName)
	if err != nil {
		return networkCR, sdnInstallConfig, err
	}
	networkCR.Spec.DefaultNetwork.OpenshiftSDNConfig.Mode = selectedNetworkPlugin

	// Put IP adresses to a yaml snippet of install-config
	sdnInstallConfig.ServiceNetwork = []string{networkConfig.ServiceNetworkCIDR}
	sdnInstallConfig.ClusterNetworks = translatedClusterNetworks

	return networkCR, sdnInstallConfig, nil
}

// Extract collects SDN configuration information from an OCP3 cluster
func (e SDNReport) Extract() (etl.Extraction, error) {
	logrus.Info("SDNReport::Extract")
	content, err := e.Config.Fetch(env.Config().GetString("MasterConfigFile"))
	if err != nil {
		return nil, err
	}

	var extraction SDNExtraction

	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	_, _, err = serializer.Decode(content, nil, &extraction.MasterConfig)
	if err != nil {
		return nil, err
	}

	return extraction, nil
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
		// host prefix default value is 23, we should mention this in readme
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

// GenYAML returns a YAML of the OAuthCRD or install config snippet
func GenYAML(parsableStuct interface{}) ([]byte, error) {
	yamlBytes, err := yaml.Marshal(parsableStuct)
	if err != nil {
		return nil, err
	}

	return yamlBytes, nil
}

// Name returns a human readable name for the transform
func (e SDNReport) Name() string {
	return "SDN"
}

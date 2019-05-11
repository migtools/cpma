package transform

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
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

type SDNTransform struct {
	Config *Config
}

const (
	apiVersion         = "operator.openshift.io/v1"
	kind               = "Network"
	defaultNetworkType = "OpenShiftSDN"
)

func (c SDNTransform) Run(content []byte) (TransformOutput, error) {
	logrus.Info("SDNTransform::Run")

	var manifests Manifests

	networkCR := SDNTranslate(content)
	networkCRYAML := GenYAML(networkCR)

	manifest := Manifest{Name: "100_CPMA-cluster-config-sdn.yaml", CRD: networkCRYAML}
	manifests = append(manifests, manifest)

	return ManifestTransformOutput{
		Config:    *c.Config,
		Manifests: manifests,
	}, nil
}

func SDNTranslate(content []byte) NetworkCR {
	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	var masterConfig configv1.MasterConfig

	_, _, err := serializer.Decode(content, nil, &masterConfig)
	if err != nil {
		HandleError(err)
	}
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

func (c SDNTransform) Extract() []byte {
	logrus.Info("SDNTransform::Extract")
	return c.Config.Fetch(c.Config.MasterConfigFile)
}

func (c SDNTransform) Validate() error {
	return nil // Simulate fine
}

func TranslateClusterNetworks(clusterNeworkEntries []configv1.ClusterNetworkEntry) []ClusterNetwork {
	translatedClusterNetworks := make([]ClusterNetwork, 0)

	for _, networkConfig := range clusterNeworkEntries {
		var translatedClusterNetwork ClusterNetwork

		translatedClusterNetwork.CIDR = networkConfig.CIDR
		translatedClusterNetwork.HostPrefix = networkConfig.HostSubnetLength

		translatedClusterNetworks = append(translatedClusterNetworks, translatedClusterNetwork)
	}

	return translatedClusterNetworks
}

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

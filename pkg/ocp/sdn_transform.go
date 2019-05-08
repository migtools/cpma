package ocp

import (
	"errors"
	"fmt"

	"github.com/fusor/cpma/pkg/ocp4"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func (c SDNTransform) Run(content []byte) (TransformOutput, error) {
	fmt.Println("SDNTransform::Run")

	const (
		apiVersion         = "operator.openshift.io/v1"
		kind               = "Network"
		defaultNetworkType = "OpenShiftSDN"
	)

	var manifests ocp4.Manifests

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
	translatedClusterNetworks := translateClusterNetworks(networkConfig.ClusterNetworks)
	networkCR.Spec.ClusterNetworks = translatedClusterNetworks

	// Translate network plugin name
	selectedNetworkPlugin, err := selectNetworkPlugin(networkConfig.NetworkPluginName)
	if err != nil {
		HandleError(err)
	}
	networkCR.Spec.DefaultNetwork.OpenshiftSDNConfig.Mode = selectedNetworkPlugin
	networkCRYAML, err := yaml.Marshal(&networkCR)
	if err != nil {
		HandleError(err)
	}

	manifest := ocp4.Manifest{Name: "100_CPMA-cluster-config-sdn.yaml", CRD: networkCRYAML}
	manifests = append(manifests, manifest)

	return ManifestTransformOutput{
		Config:    *c.Config,
		Manifests: manifests,
	}, nil
}

func (c SDNTransform) Extract() []byte {
	fmt.Println("SDNTransform::Extract")
	return c.Config.Fetch(c.Config.MasterConfigFile)
}

func (c SDNTransform) Validate() error {
	return nil // Simulate fine
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

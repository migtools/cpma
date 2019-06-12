package decode

import (
	"k8s.io/client-go/kubernetes/scheme"

	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func init() {
	legacyconfigv1.InstallLegacy(scheme.Scheme)
}

// MasterConfig unmarshals OCP3 Master
// There is no known differences between OCP3 minor versions of the master config
func MasterConfig(content []byte) (*legacyconfigv1.MasterConfig, error) {
	var masterConfig = new(legacyconfigv1.MasterConfig)

	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	_, _, err := serializer.Decode(content, nil, masterConfig)
	if err != nil {
		return nil, err
	}

	return masterConfig, nil
}

// NodeConfig unmarshals OCP3 Node
// Unknown differences between OCP3 minor versions of the node config
func NodeConfig(content []byte) (*legacyconfigv1.NodeConfig, error) {
	var nodeConfig = new(legacyconfigv1.NodeConfig)
	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	_, _, err := serializer.Decode(content, nil, nodeConfig)
	if err != nil {
		return nil, err
	}

	return nodeConfig, nil
}

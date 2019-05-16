package decode

import (
	configv1 "github.com/openshift/api/legacyconfig/v1"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	configv1.InstallLegacy(scheme.Scheme)
}

// MasterConfig unmarshals OCP3 Master
func MasterConfig(content []byte) (*configv1.MasterConfig, error) {
	var masterConfig = new(configv1.MasterConfig)

	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	_, _, err := serializer.Decode(content, nil, masterConfig)
	if err != nil {
		return nil, err
	}

	return masterConfig, nil
}

// NodeConfig unmarshals OCP3 Node
func NodeConfig(content []byte) (*configv1.NodeConfig, error) {
	var nodeConfig = new(configv1.NodeConfig)
	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	_, _, err := serializer.Decode(content, nil, nodeConfig)
	if err != nil {
		return nil, err
	}

	return nodeConfig, nil
}

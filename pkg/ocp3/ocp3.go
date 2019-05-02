package ocp3

import (
	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

func init() {
	configv1.InstallLegacy(scheme.Scheme)
}

func (cluster *Cluster) Decode(configFile ConfigFile) {
	switch configFile.Type {
	case "master":
		cluster.DecodeOpenshiftMaster(configFile.Content)
	case "node":
		cluster.DecodeOpenshiftNode(configFile.Content)
	}
}

func (cluster *Cluster) DecodeOpenshiftMaster(content []byte) {

	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	_, _, err := serializer.Decode(content, nil, &cluster.MasterConfig)
	if err != nil {
		logrus.Fatal(err)
	}
}

func (cluster *Cluster) DecodeOpenshiftNode(content []byte) {
	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	_, _, err := serializer.Decode(content, nil, &cluster.NodeConfig)
	if err != nil {
		logrus.Fatal(err)
	}
}

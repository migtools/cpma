package ocp3

import (
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

// reference:
// https://docs.openshift.com/container-platform/3.11/install_config/master_node_configuration.html

// TODO: we may want to be OCP3 minor version aware here

type Cluster struct {
	Master Master
}
type Master struct {
	Config configv1.MasterConfig
}

func init() {
	configv1.InstallLegacy(scheme.Scheme)
}

// Decode unmarshals OCP3
func (clusterOCP3 *Cluster) Decode(path string) error {
	filename := filepath.Base(path)
	if strings.Contains(filename, "master") {
		clusterOCP3.Master.Config = ParseMaster(path)
		return nil
	}
	return errors.New("No file decoded")
}

// ParseMaster unmarshals OCP3 master-config
func ParseMaster(file string) configv1.MasterConfig {
	var masterConfig configv1.MasterConfig
	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	f, err := ioutil.ReadFile(file)
	if err != nil {
		logrus.Fatal(err)
	}

	_, _, err = serializer.Decode(f, nil, &masterConfig)
	if err != nil {
		logrus.Fatal(err)
	}

	return masterConfig
}

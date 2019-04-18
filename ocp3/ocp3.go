package ocp3

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fusor/cpma/env"
	"github.com/fusor/cpma/internal/sftpclient"
	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

// reference:
// https://docs.openshift.com/container-platform/3.11/install_config/master_node_configuration.html

// TODO: we may want to be OCP3 minor version aware here

// Config represents OCP3 configuration
type Config struct {
	master  configv1.MasterConfig
	Masterf string

	node  configv1.NodeConfig
	Nodef string
}

// TODO: rework the prototype code below

func init() {
	configv1.InstallLegacy(scheme.Scheme)
}

// ParseMaster unmarshals OCP3 master-config
func (c *Config) ParseMaster() configv1.MasterConfig {
	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	f, err := ioutil.ReadFile(c.Masterf)
	if err != nil {
		logrus.Fatal(err)
	}

	_, _, err = serializer.Decode(f, nil, &c.master)
	if err != nil {
		logrus.Fatal(err)
	}
	return c.master
}

// Fetch checks whether OCP3 configuration is available and retrieves
// it in case it is not.
func (c *Config) Fetch() {
	// TODO: this function must get all the files referred from master-config
	// and node-config as well
	var src, dst string
	var err error
	var client sftpclient.Client

	source := env.Config().GetString("Source")
	outputDir := env.Config().GetString("OutputDir")

	src = filepath.Join(source, c.Masterf)
	if _, err = os.Stat(src); os.IsNotExist(err) {
		goto fetch
	}
	c.Masterf = src

	src = filepath.Join(source, c.Nodef)
	if _, err = os.Stat(src); os.IsNotExist(err) {
		goto fetch
	}
	c.Nodef = src

	logrus.Debug("Local copy of configuration files has been found, skip ssh fetch.")
	goto out

fetch:
	// We weren't successfull in locating configuration in directory
	// given by e.Source, thus we think it is fqdn.
	// TODO: Rework logic or we may want to prompt user here.

	logrus.Debug("Unable to locate configuration, attempt to fetch from ", source)
	client = sftpclient.NewClient()
	defer client.Close()

	dst = filepath.Join(outputDir, c.Masterf)
	client.GetFile(c.Masterf, dst)
	c.Masterf = dst

	dst = filepath.Join(outputDir, c.Nodef)
	client.GetFile(c.Nodef, dst)
	c.Nodef = dst

out:
	return
}

// New instantiate Config structure that represents OCP3 configuration
func New() *Config {
	return &Config{
		Masterf: "/etc/origin/master/master-config.yaml",
		Nodef:   "/etc/origin/node/node-config.yaml",
	}
}

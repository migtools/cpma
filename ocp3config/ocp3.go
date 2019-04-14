package ocp3

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fusor/cpma/env"
	"github.com/fusor/cpma/internal/sftpclient"

	configv1 "github.com/openshift/api/legacyconfig/v1"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"

	log "github.com/sirupsen/logrus"
)

// reference:
// https://docs.openshift.com/container-platform/3.11/install_config/master_node_configuration.html

// TODO: we may want to be OCP3 minor version aware here

// Config represents OCP3 configuration
type Config struct {
	master  configv1.MasterConfig
	masterf string

	node  configv1.NodeConfig
	nodef string
}

// TODO: rework the prototype code below

func init() {
	configv1.InstallLegacy(scheme.Scheme)
}

// ParseMaster unmarshals OCP3 master-config
func (c *Config) ParseMaster() configv1.MasterConfig {
	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	f, err := ioutil.ReadFile(c.masterf)
	if err != nil {
		log.Fatalln(err)
	}

	_, _, err = serializer.Decode(f, nil, &c.master)
	if err != nil {
		log.Fatalln(err)
	}

	return c.master
}

// Fetch checks whether OCP3 configuration is available and retrieves
// it in case it is not.
func (c *Config) Fetch(e *env.Info) {
	// TODO: this function must get all the files referred from master-config
	// and node-config as well

	var sftpclient *sftpclient.Client
	var src, dst string
	var err error

	src = filepath.Join(e.Source, c.masterf)
	if _, err = os.Stat(src); os.IsNotExist(err) {
		goto fetch
	}
	c.masterf = src

	src = filepath.Join(e.Source, c.nodef)
	if _, err = os.Stat(src); os.IsNotExist(err) {
		goto fetch
	}
	c.nodef = src

	goto out

fetch:
	// We weren't successfull in locating configuration in directory
	// given by e.Source, thus we think it is fqdn.
	// TODO: Rework logic or we may want to prompt user here.

	log.Debugln("unable to locate configuration, attempt to fetch from ", e.Source)

	sftpclient, err = e.SSH.NewClient(e.Source)
	if err != nil {
		log.Fatalln(err)
	}
	defer sftpclient.Close()

	dst = filepath.Join(e.OutputDir, c.masterf)
	sftpclient.GetFile(c.masterf, dst)
	c.masterf = dst

	dst = filepath.Join(e.OutputDir, c.nodef)
	sftpclient.GetFile(c.nodef, dst)
	c.nodef = dst

out:
	return
}

// New instantiate Config structure that represents OCP3 configuration
func New() *Config {
	return &Config{
		masterf: "/etc/origin/master/master-config.yaml",
		nodef:   "/etc/origin/node/node-config.yaml",
	}
}

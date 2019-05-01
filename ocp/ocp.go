package ocp

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/fusor/cpma/env"
	"github.com/fusor/cpma/network/sftpclient"
	"github.com/fusor/cpma/ocp3"
	"github.com/fusor/cpma/ocp4"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Hostname string
	Path     string
	Content  []byte
	OCP3     ocp3.Cluster
	OCP4     ocp4.Cluster
}

// Decode unmarshals OCP3
func (config *Config) Decode() {
	config.OCP3.Master.Decode(config.Content)
	// TODO: Keep for when adding Node
	//config.OCP3.DecodeNode(config.Content)
}

// GenYAML returns the list of translated CRDs
func (config *Config) GenYAML() ocp4.Manifests {
	var manifests ocp4.Manifests

	masterManifests, err := config.OCP4.Master.GenYAML()
	if err != nil {
		return nil
	}
	for _, manifest := range masterManifests {
		manifests = append(manifests, manifest)
	}
	return manifests
}

// DumpManifests creates OCDs files
func (config *Config) DumpManifests(outputDir string, manifests []ocp4.Manifest) {
	for _, manifest := range manifests {
		maniftestfile := filepath.Join(outputDir, "manifests", manifest.Name)
		os.MkdirAll(path.Dir(maniftestfile), 0755)
		err := ioutil.WriteFile(maniftestfile, manifest.CRD, 0644)
		logrus.Printf("CR manifest created: %s", maniftestfile)
		if err != nil {
			logrus.Panic(err)
		}
	}
}

// Fetch retrieves file from Host
func (config *Config) Fetch(outputDir string) {
	dst := filepath.Join(outputDir, config.Hostname, config.Path)
	sftpclient.Fetch(config.Hostname, config.Path, dst)

	f, err := ioutil.ReadFile(filepath.Join(outputDir, config.Hostname, config.Path))
	if err != nil {
		logrus.Fatal(err)
	}
	config.Content = f
}

// Translate OCP3 to OCP4
func (config *Config) Translate() {
	config.OCP4.Master.Translate(config.OCP3.Master.Config)
	// TODO: Keep for when adding Node
	//config.OCP4.Node.Translate(config.OCP3.Node.Config)
}

func (config *Config) AddMaster(hostname string) {
	masterf := env.Config().GetString("MasterConfigFile")

	if masterf == "" {
		masterf = "/etc/origin/master/master-config.yaml"
	}
	config.Hostname = hostname
	config.Path = masterf
}

func (config *Config) AddNode(hostname string) {
	nodef := env.Config().GetString("NodeConfigFile")

	if nodef == "" {
		nodef = "/etc/origin/node/node-config.yaml"
	}

	config.Hostname = hostname
	config.Path = nodef
}

package ocp

import (
	"fmt"
	//	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/fusor/cpma/pkg/ocp4"
	"github.com/fusor/cpma/pkg/sftpclient"
	"github.com/sirupsen/logrus"
)

const MasterConfigFile = "/etc/origin/master/master-config.yaml"
const NodeConfigFile = "/etc/origin/node/node-config.yaml"

func (migration *Migration) Decode(configFile ocp3.ConfigFile) {
	migration.OCP3Cluster.Decode(configFile)
}

// DumpManifests creates OCDs files
func (migration *Migration) DumpManifests(manifests []ocp4.Manifest) {
	for _, manifest := range manifests {
		maniftestfile := filepath.Join(migration.OutputDir, "manifests", manifest.Name)
		os.MkdirAll(path.Dir(maniftestfile), 0755)
		err := ioutil.WriteFile(maniftestfile, manifest.CRD, 0644)
		logrus.Printf("CR manifest created: %s", maniftestfile)
		if err != nil {
			logrus.Panic(err)
		}
	}
}

// Fetch retrieves file from Host
func (migration *Migration) Fetch(configFile *ocp3.ConfigFile) {
	dst := filepath.Join(migration.OutputDir, migration.OCP3Cluster.Hostname, configFile.Path)
	sftpclient.Fetch(migration.OCP3Cluster.Hostname, configFile.Path, dst)

	f, err := ioutil.ReadFile(filepath.Join(migration.OutputDir, migration.OCP3Cluster.Hostname, configFile.Path))
	if err != nil {
		logrus.Warning(err)
	}
	configFile.Content = f
}

// Translate OCP3 to OCP4
func (migration *Migration) Translate() {
	migration.OCP4Cluster.Master.Translate(migration.OCP3Cluster)
}

type Transform interface {
	Run() (TransformOutput, error)
	Validate() error
	Extract()
}

type TransformOutput interface {
	Flush() error
}

func (m ManifestTransformOutput) Flush() error {
	fmt.Println("Writing file data:")
	m.Migration.DumpManifests(m.Manifests)
	return nil
}

func NewTransformRunner(config Config) *TransformRunner {
	fmt.Printf("Building TransformRunner with RunnerConfig: %s\n", config.RunnerConfig)
	return &TransformRunner{Config: config.RunnerConfig}
}

func (r TransformRunner) Run(transforms []Transform) error {
	fmt.Println("TransformRunner::Run")

	// For each transform, extract the data, validate it, and run the transform.
	// Handle any errors, and finally flush the output to it's desired destination
	// NOTE: This should be parallelized with channels unless the transforms have
	// some dependency on the outputs of others
	for _, transform := range transforms {
		transform.Extract()

		if err := transform.Validate(); err != nil {
			return HandleError(err)
		}

		output, err := transform.Run()
		if err != nil {
			HandleError(err)
		}

		if err := output.Flush(); err != nil {
			HandleError(err)
		}
	}

	return nil
}

func LoadConfig() Config {
	// Mocking out the details of collecting cli input and file input
	config := Config{
		MasterConfigFile: MasterConfigFile,
		NodeConfigFile:   NodeConfigFile,
		RunnerConfig:     "some_runner_config",
	}

	fmt.Println("Loaded config")
	return config
}

func HandleError(err error) error {
	return fmt.Errorf("An error has occurred: %s", err)
}

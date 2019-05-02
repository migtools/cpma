package ocp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/fusor/cpma/pkg/ocp4"
	"github.com/fusor/cpma/pkg/sftpclient"
	"github.com/sirupsen/logrus"
)

const MasterConfigFile = "/etc/origin/master/master-config.yaml"
const NodeConfigFile = "/etc/origin/node/node-config.yaml"

// DumpManifests creates OCDs files
func (config *Config) DumpManifests(manifests []ocp4.Manifest) {
	for _, manifest := range manifests {
		maniftestfile := filepath.Join(config.OutputDir, "manifests", manifest.Name)
		os.MkdirAll(path.Dir(maniftestfile), 0755)
		err := ioutil.WriteFile(maniftestfile, manifest.CRD, 0644)
		logrus.Printf("CR manifest created: %s", maniftestfile)
		if err != nil {
			logrus.Panic(err)
		}
	}
}

func (config *Config) Fetch(path string) []byte {
	dst := filepath.Join(config.OutputDir, config.Hostname, path)
	sftpclient.Fetch(config.Hostname, path, dst)

	f, err := ioutil.ReadFile(filepath.Join(config.OutputDir, config.Hostname, path))
	if err != nil {
		logrus.Warning(err)
	}
	return f
}

type Transform interface {
	Run([]byte) (TransformOutput, error)
	Validate() error
	Extract() []byte
}

type TransformOutput interface {
	Flush() error
}

func (m ManifestTransformOutput) Flush() error {
	fmt.Println("Writing file data:")
	m.Config.DumpManifests(m.Manifests)
	return nil
}

func NewTransformRunner(config Config) *TransformRunner {
	return &TransformRunner{}
}

func (r TransformRunner) Run(transforms []Transform) error {
	fmt.Println("TransformRunner::Run")

	// For each transform, extract the data, validate it, and run the transform.
	// Handle any errors, and finally flush the output to it's desired destination
	// NOTE: This should be parallelized with channels unless the transforms have
	// some dependency on the outputs of others
	for _, transform := range transforms {
		content := transform.Extract()

		if err := transform.Validate(); err != nil {
			return HandleError(err)
		}

		output, err := transform.Run(content)
		if err != nil {
			HandleError(err)
		}

		if err := output.Flush(); err != nil {
			HandleError(err)
		}
	}

	return nil
}

func HandleError(err error) error {
	return fmt.Errorf("An error has occurred: %s", err)
}

package transform

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/fusor/cpma/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/oauth"
	"github.com/fusor/cpma/pkg/transform/secrets"
	"github.com/sirupsen/logrus"
)

// OCP4InstallMsg message about using generated manifests
const OCP4InstallMsg = `To install OCP4 run the installer as follow in order to add CRDs:
' /openshift-install --dir $INSTALL_DIR create install-config'
'./openshift-install --dir $INSTALL_DIR create manifests'
# Copy generated CRD manifest files  to '$INSTALL_DIR/openshift/'
# Edit them if needed, then run installation:
'./openshift-install --dir $INSTALL_DIR  create cluster'`
const MasterConfigFile = "/etc/origin/master/master-config.yaml"
const NodeConfigFile = "/etc/origin/node/node-config.yaml"
const RegistriesConfigFile = "/etc/containers/registries.conf"

type Cluster struct {
	Master Master
}

type Master struct {
	OAuth   oauth.OAuthCRD
	Secrets []secrets.Secret
}

type Manifests []Manifest

// Manifest holds a CRD object
type Manifest struct {
	Name string
	CRD  []byte
}

type Config struct {
	MasterConfigFile     string
	NodeConfigFile       string
	RegistriesConfigFile string
	OutputDir            string
	Hostname             string
}

type ManifestTransformOutput struct {
	Config    Config
	Manifests []Manifest
}

type TransformRunner struct {
	Config string
}

// GetFile allows to mock file retrieval
var GetFile = io.GetFile

func Start() {
	config := Config{}
	config.OutputDir = env.Config().GetString("OutputDir")
	config.Hostname = env.Config().GetString("Source")
	config.MasterConfigFile = MasterConfigFile
	config.RegistriesConfigFile = RegistriesConfigFile
	transformRunner := NewTransformRunner(config)

	if err := transformRunner.Run([]Transform{
		OAuthTransform{
			Config: &config,
		},
		SDNTransform{
			Config: &config,
		},
		RegistriesTransform{
			Config: &config,
		},
	}); err != nil {
		logrus.WithError(err).Fatalf("%s", err.Error())
	}
}

// DumpManifests creates OCDs files
func (config *Config) DumpManifests(manifests []Manifest) {
	for _, manifest := range manifests {
		maniftestfile := filepath.Join(env.Config().GetString("OutputDir"), "manifests", manifest.Name)
		os.MkdirAll(path.Dir(maniftestfile), 0755)
		err := ioutil.WriteFile(maniftestfile, manifest.CRD, 0644)
		logrus.Printf("CRD:Added: %s", maniftestfile)
		if err != nil {
			logrus.Panic(err)
		}
	}
}

func (config *Config) Fetch(path string) []byte {
	dst := filepath.Join(config.OutputDir, config.Hostname, path)
	f := GetFile(config.Hostname, path, dst)
	logrus.Printf("File:Loaded: %s", dst)

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
	logrus.Info("Writing file data:")
	m.Config.DumpManifests(m.Manifests)
	return nil
}

func NewTransformRunner(config Config) *TransformRunner {
	return &TransformRunner{}
}

func (r TransformRunner) Run(transforms []Transform) error {
	logrus.Info("TransformRunner::Run")

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
	logrus.WithError(err).Fatalf("An error has occurred: %s\n", err)
	return err
}

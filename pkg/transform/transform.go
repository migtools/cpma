package transform

import (
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

type TransformRunner struct {
	Config string
}

type Extraction interface {
	Transform() (TransformOutput, error)
	Validate() error
}

type Transform interface {
	Extract() Extraction
}

type TransformOutput interface {
	Flush() error
}

// GetFile allows to mock file retrieval
var GetFile = io.GetFile

func Start() {
	config := LoadConfig()
	transformRunner := NewTransformRunner(config)

	if err := transformRunner.Transform([]Transform{
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

func LoadConfig() Config {
	logrus.Info("Loaded config")

	config := Config{}
	config.OutputDir = env.Config().GetString("OutputDir")
	config.Hostname = env.Config().GetString("Source")
	config.MasterConfigFile = MasterConfigFile
	config.RegistriesConfigFile = RegistriesConfigFile

	return config
}

func (config *Config) Fetch(path string) []byte {
	dst := filepath.Join(config.OutputDir, config.Hostname, path)
	f := GetFile(config.Hostname, path, dst)
	logrus.Printf("File:Loaded: %s", dst)

	return f
}

func (r TransformRunner) Transform(transforms []Transform) error {
	logrus.Info("TransformRunner::Transform")

	// For each transform, extract the data, validate it, and run the transform.
	// Handle any errors, and finally flush the output to it's desired destination
	// NOTE: This should be parallelized with channels unless the transforms have
	// some dependency on the outputs of others
	for _, transform := range transforms {
		extraction := transform.Extract()

		if err := extraction.Validate(); err != nil {
			HandleError(err)
			continue
		}

		output, err := extraction.Transform()
		if err != nil {
			HandleError(err)
		}

		if err := output.Flush(); err != nil {
			HandleError(err)
		}
	}

	return nil
}

func NewTransformRunner(config Config) *TransformRunner {
	return &TransformRunner{}
}

func HandleError(err error) error {
	logrus.Warnf("%s\n", err)
	return err
}

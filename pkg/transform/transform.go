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

// MasterConfigFile file path to the master-config.yaml
const MasterConfigFile = "/etc/origin/master/master-config.yaml"

// NodeConfigFile file path for node-config.yaml
const NodeConfigFile = "/etc/origin/node/node-config.yaml"

// RegistriesConfigFile file path for registries.conf
const RegistriesConfigFile = "/etc/containers/registries.conf"

// Cluster contains a cluster
type Cluster struct {
	Master Master
}

// Master is a cluster Master
type Master struct {
	OAuth   oauth.CRD
	Secrets []secrets.Secret
}

// Manifest to be exported for use with OCP 4
type Manifest struct {
	Name string
	CRD  []byte
}

// Config contains CPMA configuration information
type Config struct {
	MasterConfigFile     string
	NodeConfigFile       string
	RegistriesConfigFile string
	OutputDir            string
	Hostname             string
}

// Runner a generic transform runner
type Runner struct {
	Config string
}

// Extraction is a generic data extraction
type Extraction interface {
	Transform() (Output, error)
	Validate() error
}

// Transform is a generic transform
type Transform interface {
	Extract() Extraction
}

// Output is a generic output type
type Output interface {
	Flush() error
}

// GetFile allows to mock file retrieval
var GetFile = io.GetFile

//Start generating manifests to be used with Openshift 4
func Start() {
	config := LoadConfig()
	runner := NewRunner(config)

	if err := runner.Transform([]Transform{
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

// LoadConfig collects and stores configuration for CPMA
func LoadConfig() Config {
	logrus.Info("Loaded config")

	config := Config{}
	config.OutputDir = env.Config().GetString("OutputDir")
	config.Hostname = env.Config().GetString("Source")
	config.MasterConfigFile = MasterConfigFile
	config.RegistriesConfigFile = RegistriesConfigFile

	return config
}

// Fetch files from the OCP3 cluster
func (config *Config) Fetch(path string) []byte {
	dst := filepath.Join(config.OutputDir, config.Hostname, path)
	f := GetFile(config.Hostname, path, dst)
	logrus.Printf("File:Loaded: %s", dst)

	return f
}

// Transform is the process run to complete a transform
func (r Runner) Transform(transforms []Transform) error {
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

// NewRunner creates a new Runner
func NewRunner(config Config) *Runner {
	return &Runner{}
}

// HandleError handles errors
func HandleError(err error) error {
	logrus.Warnf("%s\n", err)
	return err
}

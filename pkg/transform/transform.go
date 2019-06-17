package transform

import (
	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/configmaps"
	"github.com/fusor/cpma/pkg/transform/secrets"
	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
		configv1 "github.com/openshift/api/config/v1"

)

const (
	// NoConfidence represents report items we can not migrate
	NoConfidence = iota

	// ModerateConfidence represents report items we can migrate with caveats
	ModerateConfidence

	// HighConfidence represents report items we can migrate without issue
	HighConfidence

	// OCP4InstallMsg message about using generated manifests
	OCP4InstallMsg = `To install OCP4 run the installer as follow in order to add CRDs:
' /openshift-install --dir $INSTALL_DIR create install-config'
'./openshift-install --dir $INSTALL_DIR create manifests'
# Copy generated CRD manifest files  to '$INSTALL_DIR/openshift/'
# Edit them if needed, then run installation:
'./openshift-install --dir $INSTALL_DIR  create cluster'`
)

// Cluster contains a cluster
type Cluster struct {
	Master Master
}

// Master is a cluster Master
type Master struct {
	OAuth      configv1.OAuth
	Secrets    []*secrets.Secret
	ConfigMaps []*configmaps.ConfigMap
}

// Manifest to be exported for use with OCP 4
type Manifest struct {
	Name string
	CRD  []byte
}

// Report of OCP 4 component configuration compatibility
type Report struct {
	Name       string `json:"name"`
	Kind       string `json:"kind"`
	Supported  bool   `json:"supported"`
	Confidence int    `json:"confidence"`
	Comment    string `json:"comment"`
}

// Runner a generic transform runner
type Runner struct {
}

// Extraction is a generic data extraction
type Extraction interface {
	Transform() ([]Output, error)
	Validate() error
}

// Transform is a generic transform
type Transform interface {
	Extract() (Extraction, error)
	Name() string
}

// Output is a generic output type
type Output interface {
	Flush() error
}

//Start generating manifests to be used with Openshift 4
func Start() {
	runner := NewRunner()

	openReports()

	runner.Transform([]Transform{
		APITransform{},
		CrioTransform{},
		DockerTransform{},
		ETCDTransform{},
		OAuthTransform{},
		SDNTransform{},
		ImageTransform{},
	})
}

// Transform is the process run to complete a transform
func (r Runner) Transform(transforms []Transform) {
	logrus.Info("TransformRunner::Transform")

	// For each transform, extract the data, validate it, and run the transform.
	// Handle any errors, and finally flush the output to it's desired destination
	// NOTE: This should be parallelized with channels unless the transforms have
	// some dependency on the outputs of others
	for _, transform := range transforms {
		extraction, err := transform.Extract()
		if err != nil {
			HandleError(err, transform.Name())
			continue
		}

		if err := extraction.Validate(); err != nil {
			HandleError(err, transform.Name())
			continue
		}

		outputs, err := extraction.Transform()
		if err != nil {
			HandleError(err, transform.Name())
			continue
		}
		for _, output := range outputs {
			if err := output.Flush(); err != nil {
				HandleError(err, transform.Name())
				continue
			}
		}
	}
}

// NewRunner creates a new Runner
func NewRunner() *Runner {
	return &Runner{}
}

// HandleError handles errors
func HandleError(err error, transformType string) error {
	logrus.Warnf("Skipping %s - %s\n", transformType, err)
	return err
}

// GenYAML returns a YAML of the CR
func GenYAML(CR interface{}) ([]byte, error) {
	yamlBytes, err := yaml.Marshal(CR)
	if err != nil {
		return nil, err
	}

	return yamlBytes, nil
}

func openReports() {
	jsonfile := "report.json"
	emptyReport := []byte("[]")

	err := io.WriteFile(emptyReport, jsonfile)
	if err != nil {
		logrus.Errorf("unable to open report file: %s", jsonfile)
	}
}

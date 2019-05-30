package transform

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/transform/configmaps"
	"github.com/fusor/cpma/pkg/transform/oauth"
	"github.com/fusor/cpma/pkg/transform/secrets"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	// OCP4InstallMsg message about using generated manifests
	OCP4InstallMsg = `To install OCP4 run the installer as follow in order to add CRDs:
' /openshift-install --dir $INSTALL_DIR create install-config'
'./openshift-install --dir $INSTALL_DIR create manifests'
# Copy generated CRD manifest files  to '$INSTALL_DIR/openshift/'
# Edit them if needed, then run installation:
'./openshift-install --dir $INSTALL_DIR  create cluster'`

	// ReportOutputType identifier string for a report run
	ReportOutputType = "report"

	// ConvertOutputType identifier string for a convert run
	ConvertOutputType = "convert"

	// SDNComponentName is the registry component string
	SDNComponentName = "SDN"

	// OAuthComponentName is the registry component string
	OAuthComponentName = "Oauth"

	// RegistriesComponentName is the registry component string
	RegistriesComponentName = "Registries"
)

// Cluster contains a cluster
type Cluster struct {
	Master Master
}

// Master is a cluster Master
type Master struct {
	OAuth      oauth.CRD
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
	Confidence string `json:"confidence"`
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
		OAuthTransform{},
		SDNTransform{},
		RegistriesTransform{},
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
	logrus.Warnf("Skipping %s, see error below\n", transformType)
	logrus.Warnf("%s\n", err)
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
	jsonfile := filepath.Join(env.Config().GetString("OutputDir"), "report.json")
	os.MkdirAll(path.Dir(jsonfile), 0755)

	err := ioutil.WriteFile(jsonfile, []byte("[]"), 0644)
	if err != nil {
		logrus.Errorf("unable to open report file: %s", jsonfile)
	}
}

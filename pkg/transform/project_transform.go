package transform

import (
	"fmt"

	"github.com/fusor/cpma/pkg/decode"
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/project"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/sirupsen/logrus"
)

// ProjectComponentName is the Project component string
const ProjectComponentName = "Project"

// ProjectExtraction is a Project specific extraction
type ProjectExtraction struct {
	legacyconfigv1.MasterConfig
}

// ProjectTransform is a Project specific transform
type ProjectTransform struct {
}

// Transform converts data collected from an OCP3 into a useful output
func (e ProjectExtraction) Transform() ([]Output, error) {
	logrus.Info("ProjectTransform::Transform")
	manifests, err := e.buildManifestOutput()
	if err != nil {
		return nil, err
	}
	reports, err := e.buildReportOutput()
	if err != nil {
		return nil, err
	}
	outputs := []Output{manifests, reports}
	return outputs, nil
}

func (e ProjectExtraction) buildManifestOutput() (Output, error) {
	var manifests []Manifest

	projectCR, err := project.Translate(e.MasterConfig.ProjectConfig)
	if err != nil {
		return nil, err
	}

	projectCRYAML, err := GenYAML(projectCR)
	if err != nil {
		return nil, err
	}

	manifest := Manifest{Name: "100_CPMA-cluster-config-project.yaml", CRD: projectCRYAML}
	manifests = append(manifests, manifest)

	return ManifestOutput{
		Manifests: manifests,
	}, nil
}

func (e ProjectExtraction) buildReportOutput() (Output, error) {
	reportOutput := ReportOutput{
		Component: ProjectComponentName,
	}

	reportOutput.Reports = append(reportOutput.Reports,
		Report{
			Name:       "ProjectRequestMessage",
			Kind:       "ProjectConfig",
			Supported:  true,
			Confidence: HighConfidence,
			Comment:    "Networks must be configured during installation",
		})

	reportOutput.Reports = append(reportOutput.Reports,
		Report{
			Name:       "ProjectRequestTemplate",
			Kind:       "ProjectConfig",
			Supported:  true,
			Confidence: HighConfidence,
			Comment:    "Networks must be configured during installation",
		})

	reportOutput.Reports = append(reportOutput.Reports,
		Report{
			Name:       "DefaultNodeSelector",
			Kind:       "ProjectConfig",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    fmt.Sprintf("Not supported in OCP4: %s", e.ProjectConfig.DefaultNodeSelector),
		})

	reportOutput.Reports = append(reportOutput.Reports,
		Report{
			Name:       "SecurityAllocator.mcsAllocatorRange",
			Kind:       "ProjectConfig",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    fmt.Sprintf("Not supported in OCP4: %s", e.ProjectConfig.SecurityAllocator.MCSAllocatorRange),
		})

	reportOutput.Reports = append(reportOutput.Reports,
		Report{
			Name:       "SecurityAllocator.mcsLabelsPerProject",
			Kind:       "ProjectConfig",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    fmt.Sprintf("Not supported in OCP4: %d", e.ProjectConfig.SecurityAllocator.MCSLabelsPerProject),
		})

	reportOutput.Reports = append(reportOutput.Reports,
		Report{
			Name:       "SecurityAllocator.uidAllocatorRange",
			Kind:       "ProjectConfig",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    fmt.Sprintf("Not supported in OCP4: %s", e.ProjectConfig.SecurityAllocator.UIDAllocatorRange),
		})

	return reportOutput, nil
}

// Extract collects Project configuration information from an OCP3 cluster
func (e ProjectTransform) Extract() (Extraction, error) {
	logrus.Info("ProjectTransform::Extract")

	content, err := io.FetchFile(env.Config().GetString("MasterConfigFile"))
	if err != nil {
		return nil, err
	}

	masterConfig, err := decode.MasterConfig(content)
	if err != nil {
		return nil, err
	}

	var extraction ProjectExtraction
	extraction.MasterConfig = *masterConfig

	return extraction, nil
}

// Validate the data extracted from the OCP3 cluster
func (e ProjectExtraction) Validate() error {
	err := project.Validate(e.MasterConfig.ProjectConfig)
	if err != nil {
		return err
	}

	return nil
}

// Name returns a human readable name for the transform
func (e ProjectTransform) Name() string {
	return ProjectComponentName
}

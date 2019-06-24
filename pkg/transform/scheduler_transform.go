package transform

import (
	"github.com/fusor/cpma/pkg/decode"
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/scheduler"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/sirupsen/logrus"
)

// SchedulerComponentName is the Scheduler component string
const SchedulerComponentName = "Scheduler"

// SchedulerExtraction is a Scheduler specific extraction
type SchedulerExtraction struct {
	legacyconfigv1.MasterConfig
}

// SchedulerTransform is a Scheduler specific transform
type SchedulerTransform struct {
}

// Transform converts data collected from an OCP3 into a useful output
func (e SchedulerExtraction) Transform() ([]Output, error) {
	logrus.Info("SchedulerTransform::Transform")
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

func (e SchedulerExtraction) buildManifestOutput() (Output, error) {
	var manifests []Manifest

	schedulerCR, err := scheduler.Translate(e.MasterConfig)
	if err != nil {
		return nil, err
	}

	schedulerCRYAML, err := GenYAML(schedulerCR)
	if err != nil {
		return nil, err
	}

	manifest := Manifest{Name: "100_CPMA-cluster-config-scheduler.yaml", CRD: schedulerCRYAML}
	manifests = append(manifests, manifest)

	return ManifestOutput{
		Manifests: manifests,
	}, nil
}

func (e SchedulerExtraction) buildReportOutput() (Output, error) {
	componentReport := ComponentReport{
		Component: SchedulerComponentName,
	}

	componentReport.Reports = append(componentReport.Reports,
		Report{
			Name:       "DefaultNodeSelector",
			Kind:       "ProjectConfig",
			Supported:  true,
			Confidence: HighConfidence,
			Comment:    "",
		})

	reportOutput := ReportOutput{
		ComponentReports: []ComponentReport{componentReport},
	}

	return reportOutput, nil
}

// Extract collects Scheduler configuration information from an OCP3 cluster
func (e SchedulerTransform) Extract() (Extraction, error) {
	logrus.Info("SchedulerTransform::Extract")

	content, err := io.FetchFile(env.Config().GetString("MasterConfigFile"))
	if err != nil {
		return nil, err
	}

	masterConfig, err := decode.MasterConfig(content)
	if err != nil {
		return nil, err
	}

	var extraction SchedulerExtraction
	extraction.MasterConfig = *masterConfig

	return extraction, nil
}

// Validate the data extracted from the OCP3 cluster
func (e SchedulerExtraction) Validate() error {
	err := scheduler.Validate(e.MasterConfig)
	if err != nil {
		return err
	}

	return nil
}

// Name returns a human readable name for the transform
func (e SchedulerTransform) Name() string {
	return SchedulerComponentName
}

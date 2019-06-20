package transform

import (
	"github.com/BurntSushi/toml"
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// CrioComponentName is the name of the Crio component
const CrioComponentName = "Crio"

// CrioExtraction holds Crio data extracted from OCP3
type CrioExtraction struct {
	Crio Crio `toml:"crio"`
}

// Crio holds transformable/reportable content from crio.conf
type Crio struct {
	PidsLimit  int64  `toml:"pidsLimit"`
	LogLevel   string `toml:"logLevel"`
	LogSizeMax int64  `toml:"logSizeMax"`
	InfraImage string `toml:"infraImage"`
}

// CrioCR is a is a Crio Cluster Resource
type CrioCR struct {
	APIVersion string       `json:"apiVersion"`
	Kind       string       `json:"kind"`
	Metadata   CrioMetadata `json:"metadata"`
	Spec       CrioSpec     `json:"spec"`
}

// CrioMetadata is the Metadata for a Crio CR
type CrioMetadata struct {
	Name string `json:"name"`
}

//CrioSpec is the Spec for a Crio CR
type CrioSpec struct {
	MachineConfigPoolSelector MachineConfigPoolSelector `json:"machineConfigPoolSelector"`
	ContainerRuntimeConfig    ContainerRuntimeConfig    `json:"containerRuntimeConfig"`
}

// MachineConfigPoolSelector is the Pool Selector for a Machine Config
type MachineConfigPoolSelector struct {
	MatchLabels MatchLabels `json:"matchLabels"`
}

// MatchLabels matches the labels for a Pool Selector
type MatchLabels struct {
	CustomCrio string `json:"custom-crio"`
}

// ContainerRuntimeConfig contains a Crio Runtime Machine Config
type ContainerRuntimeConfig struct {
	PidsLimit  int64  `json:"pidsLimit,omitempty"`
	LogLevel   string `json:"logLevel,omitempty"`
	LogSizeMax int64  `json:"logSizeMax,omitempty"`
	InfraImage string `json:"infraImage,omitempty"`
}

// CrioTransform is an Crio specific transform
type CrioTransform struct {
}

// Transform converts data collected from an OCP3 into a useful output
func (e CrioExtraction) Transform() ([]Output, error) {
	logrus.Info("CrioTransform::Transform")
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

func (e CrioExtraction) buildManifestOutput() (Output, error) {
	var manifests []Manifest

	const (
		apiVersion = "machineconfiguration.openshift.io/v1"
		kind       = "ContainerRuntimeConfig"
		name       = "set-log-and-pid"
		annokey    = "release.openshift.io/create-only"
		annoval    = "true"
	)

	var crioCR CrioCR
	crioCR.APIVersion = apiVersion
	crioCR.Kind = kind
	crioCR.Metadata.Name = name
	crioCR.Spec.MachineConfigPoolSelector.MatchLabels.CustomCrio = name
	crioCR.Spec.ContainerRuntimeConfig.PidsLimit = e.Crio.PidsLimit
	crioCR.Spec.ContainerRuntimeConfig.LogLevel = e.Crio.LogLevel
	crioCR.Spec.ContainerRuntimeConfig.LogSizeMax = e.Crio.LogSizeMax
	crioCR.Spec.ContainerRuntimeConfig.InfraImage = e.Crio.InfraImage

	crioCRYAML, err := yaml.Marshal(&crioCR)
	if err != nil {
		return nil, err
	}

	manifest := Manifest{Name: "100_CPMA-crio-config.yaml", CRD: crioCRYAML}
	manifests = append(manifests, manifest)

	return ManifestOutput{
		Manifests: manifests,
	}, nil
}

func (e CrioExtraction) buildReportOutput() (Output, error) {
	componentReport := ComponentReport{
		Component: CrioComponentName,
	}

	var confidence = HighConfidence
	var supported = true

	if e.Crio.PidsLimit != 0 {
		componentReport.Reports = append(componentReport.Reports,
			Report{
				Name:       "pidsLimit",
				Kind:       "Configuration",
				Supported:  supported,
				Confidence: confidence,
			})
	}
	if e.Crio.LogLevel != "" {
		componentReport.Reports = append(componentReport.Reports,
			Report{
				Name:       "logLevel",
				Kind:       "Configuration",
				Supported:  supported,
				Confidence: confidence,
			})
	}
	if e.Crio.LogSizeMax != 0 {
		componentReport.Reports = append(componentReport.Reports,
			Report{
				Name:       "logSizeMax",
				Kind:       "Configuration",
				Supported:  supported,
				Confidence: confidence,
			})
	}
	if e.Crio.InfraImage != "" {
		componentReport.Reports = append(componentReport.Reports,
			Report{
				Name:       "infrImage",
				Kind:       "Configuration",
				Supported:  supported,
				Confidence: confidence,
			})
	}

	reportOutput := ReportOutput{
		ComponentReports: []ComponentReport{componentReport},
	}

	return reportOutput, nil
}

// Extract collects Crio configuration from an OCP3 cluster
func (e CrioTransform) Extract() (Extraction, error) {
	logrus.Info("CrioTransform::Extract")
	content, err := io.FetchFile(env.Config().GetString("CrioConfigFile"))
	if err != nil {
		return nil, err
	}

	var extraction CrioExtraction
	if _, err := toml.Decode(string(content), &extraction); err != nil {
		return nil, errors.Wrap(err, "Failed to decode crio, see error")
	}
	return extraction, nil
}

// Validate confirms we have recieved good Crio configuration data during Extract
func (e CrioExtraction) Validate() error {
	if e.Crio.PidsLimit == 0 &&
		e.Crio.LogSizeMax == 0 &&
		e.Crio.LogLevel == "" &&
		e.Crio.InfraImage == "" {
		return errors.New("no supported crio configuration found")
	}

	return nil
}

// Name returns a human readable name for the transform
func (e CrioTransform) Name() string {
	return CrioComponentName
}

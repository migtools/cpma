package transform_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/transform"
	"github.com/fusor/cpma/pkg/transform/reportoutput"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var loadCrioExtraction = func() transform.CrioExtraction {
	file := "testdata/crio.conf"
	content, _ := ioutil.ReadFile(file)
	var config transform.Crios
	_, err := toml.Decode(string(content), &config)
	if err != nil {
		fmt.Printf("Error decoding file: %s\n", file)
	}

	var extraction transform.CrioExtraction
	extraction.Runtime = config["crio"].Runtime
	return extraction
}()

func TestCrioExtractionTransform(t *testing.T) {
	var expectedManifests []transform.Manifest

	var expectedCrd transform.CrioCR
	expectedCrd.APIVersion = "machineconfiguration.openshift.io/v1"
	expectedCrd.Kind = "ContainerRuntimeConfig"
	expectedCrd.Metadata.Name = "set-log-and-pid"
	expectedCrd.Spec.MachineConfigPoolSelector.MatchLabels.CustomCrio = "set-log-and-pid"
	expectedCrd.Spec.ContainerRuntimeConfig.PidsLimit = 2048
	expectedCrd.Spec.ContainerRuntimeConfig.LogLevel = "debug"
	expectedCrd.Spec.ContainerRuntimeConfig.LogSizeMax = 100000

	crioCRYAML, err := yaml.Marshal(&expectedCrd)
	require.NoError(t, err)

	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-crio-config.yaml", CRD: crioCRYAML})

	expectedReport := reportoutput.ComponentReport{
		Component: "Crio",
	}

	expectedReport.Reports = append(expectedReport.Reports,
		reportoutput.Report{
			Name:       "pidsLimit",
			Kind:       "Configuration",
			Supported:  true,
			Confidence: 2,
		})
	expectedReport.Reports = append(expectedReport.Reports,
		reportoutput.Report{
			Name:       "logLevel",
			Kind:       "Configuration",
			Supported:  true,
			Confidence: 2,
		})
	expectedReport.Reports = append(expectedReport.Reports,
		reportoutput.Report{
			Name:       "logSizeMax",
			Kind:       "Configuration",
			Supported:  true,
			Confidence: 2,
		})

	expectedReportOutput := reportoutput.ReportOutput{
		ComponentReports: []reportoutput.ComponentReport{expectedReport},
	}

	testCases := []struct {
		name              string
		expectedManifests []transform.Manifest
		expectedReports   reportoutput.ReportOutput
	}{
		{
			name:              "transform crio extraction",
			expectedManifests: expectedManifests,
			expectedReports:   expectedReportOutput,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualManifestsChan := make(chan []transform.Manifest)
			transform.FinalReportOutput = transform.Report{}

			// Override flush method
			transform.ManifestOutputFlush = func(manifests []transform.Manifest) error {
				actualManifestsChan <- manifests
				return nil
			}
			transform.ReportOutputFlush = func(reports transform.Report) error {
				return nil
			}

			testExtraction := loadCrioExtraction

			go func() {
				env.Config().Set("Reporting", true)
				env.Config().Set("Manifests", true)

				transformOutput, err := testExtraction.Transform()
				if err != nil {
					t.Error(err)
				}
				for _, output := range transformOutput {
					output.Flush()
				}
				transform.FinalReportOutput.Flush()
			}()

			actualManifests := <-actualManifestsChan
			assert.Equal(t, actualManifests, tc.expectedManifests)

			assert.Equal(t, transform.FinalReportOutput.Report, tc.expectedReports)
		})

	}
}

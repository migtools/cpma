package transform_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/fusor/cpma/pkg/transform"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var loadCrioExtraction = func() transform.CrioExtraction {
	file := "testdata/crio.conf"
	content, _ := ioutil.ReadFile(file)
	var extraction transform.CrioExtraction
	_, err := toml.Decode(string(content), &extraction)
	if err != nil {
		fmt.Printf("Error decoding file: %s\n", file)
	}

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
	expectedCrd.Spec.ContainerRuntimeConfig.InfraImage = "image/infraImage:1"

	crioCRYAML, err := yaml.Marshal(&expectedCrd)
	require.NoError(t, err)

	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-crio-config.yaml", CRD: crioCRYAML})

	expectedReport := transform.ComponentReport{
		Component: "Crio",
	}

	expectedReport.Reports = append(expectedReport.Reports,
		transform.Report{
			Name:       "pidsLimit",
			Kind:       "Configuration",
			Supported:  true,
			Confidence: 2,
		})
	expectedReport.Reports = append(expectedReport.Reports,
		transform.Report{
			Name:       "logLevel",
			Kind:       "Configuration",
			Supported:  true,
			Confidence: 2,
		})
	expectedReport.Reports = append(expectedReport.Reports,
		transform.Report{
			Name:       "logSizeMax",
			Kind:       "Configuration",
			Supported:  true,
			Confidence: 2,
		})
	expectedReport.Reports = append(expectedReport.Reports,
		transform.Report{
			Name:       "infrImage",
			Kind:       "Configuration",
			Supported:  true,
			Confidence: 2,
		})

	expectedReportOutput := transform.ReportOutput{
		ComponentReports: []transform.ComponentReport{expectedReport},
	}

	testCases := []struct {
		name              string
		expectedManifests []transform.Manifest
		expectedReports   transform.ReportOutput
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
			actualReportsChan := make(chan transform.ReportOutput)

			// Override flush method
			transform.ManifestOutputFlush = func(manifests []transform.Manifest) error {
				actualManifestsChan <- manifests
				return nil
			}
			transform.ReportOutputFlush = func(reports transform.ReportOutput) error {
				actualReportsChan <- reports
				return nil
			}

			testExtraction := loadCrioExtraction

			go func() {
				transformOutput, err := testExtraction.Transform()
				if err != nil {
					t.Error(err)
				}
				for _, output := range transformOutput {
					output.Flush()
				}
			}()

			actualManifests := <-actualManifestsChan
			assert.Equal(t, actualManifests, tc.expectedManifests)
			actualReports := <-actualReportsChan
			assert.Equal(t, actualReports, tc.expectedReports)
		})

	}
}

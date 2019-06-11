package transform_test

import (
	"io/ioutil"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/fusor/cpma/pkg/transform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func loadCrioExtraction() (transform.CrioExtraction, error) {
	// TODO: Something is broken here in a way that it's causing the translaters
	// to fail. Need some help with creating test identiy providers in a way
	// that won't crash the translator

	// Build example identity providers, this is straight copy pasted from
	// oauth test, IMO this loading of example identity providers should be
	// some shared test helper
	file := "testdata/crio.conf" // File copied into transform pkg testdata
	content, _ := ioutil.ReadFile(file)
	var extraction transform.CrioExtraction
	_, err := toml.Decode(string(content), &extraction)

	return extraction, err
}

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

	expectedReport := transform.ReportOutput{
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

	testCases := []struct {
		name              string
		expectedManifests []transform.Manifest
		expectedReports   transform.ReportOutput
	}{
		{
			name:              "transform crio extraction",
			expectedManifests: expectedManifests,
			expectedReports:   expectedReport,
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

			testExtraction, err := loadCrioExtraction()
			require.NoError(t, err)

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

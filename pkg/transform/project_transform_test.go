package transform_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/konveyor/cpma/pkg/decode"
	"github.com/konveyor/cpma/pkg/env"
	"github.com/konveyor/cpma/pkg/io"
	"github.com/konveyor/cpma/pkg/transform"
	"github.com/konveyor/cpma/pkg/transform/reportoutput"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadProjectExtraction() (transform.ProjectExtraction, error) {
	var extraction transform.ProjectExtraction

	masterConfigContent, _ := ioutil.ReadFile("testdata/master_config-project.yaml")
	masterConfig, err := decode.MasterConfig(masterConfigContent)
	extraction.MasterConfig = *masterConfig

	return extraction, err
}

func TestProjectExtractionTransform(t *testing.T) {
	var expectedManifests []transform.Manifest

	expectedProjectCRYAML, err := ioutil.ReadFile("testdata/expected-CR-project.yaml")
	require.NoError(t, err)

	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-project.yaml", CRD: expectedProjectCRYAML})

	expectedReport := reportoutput.ReportOutput{}
	jsonData, err := io.ReadFile("testdata/expected-report-project.json")
	require.NoError(t, err)

	err = json.Unmarshal(jsonData, &expectedReport)
	require.NoError(t, err)

	testCases := []struct {
		name              string
		expectedManifests []transform.Manifest
		expectedReports   reportoutput.ReportOutput
	}{
		{
			name:              "transform project extraction",
			expectedManifests: expectedManifests,
			expectedReports:   expectedReport,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualManifestsChan := make(chan []transform.Manifest)
			actualReportsChan := make(chan reportoutput.ReportOutput)
			transform.FinalReportOutput = transform.Report{}

			// Override flush method
			transform.ManifestOutputFlush = func(manifests []transform.Manifest) error {
				actualManifestsChan <- manifests
				return nil
			}
			transform.ReportOutputFlush = func(reports transform.Report) error {
				actualReportsChan <- reports.Report
				return nil
			}

			testExtraction, err := loadProjectExtraction()
			require.NoError(t, err)

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
			actualReports := <-actualReportsChan
			assert.Equal(t, actualReports.ComponentReports, tc.expectedReports.ComponentReports)
		})
	}
}

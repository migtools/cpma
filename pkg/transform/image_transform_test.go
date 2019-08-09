package transform_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/fusor/cpma/pkg/decode"
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform"
	"github.com/fusor/cpma/pkg/transform/reportoutput"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadImageExtraction() (transform.ImageExtraction, error) {
	var extraction transform.ImageExtraction

	registriesContent, _ := ioutil.ReadFile("testdata/registries.conf")
	_, err := toml.Decode(string(registriesContent), &extraction.RegistriesConfig)

	masterConfigContent, _ := ioutil.ReadFile("testdata/master_config-image.yaml")
	masterConfig, err := decode.MasterConfig(masterConfigContent)
	extraction.MasterConfig = *masterConfig

	return extraction, err
}

func TestImageExtractionTransform(t *testing.T) {
	var expectedManifests []transform.Manifest

	expectedImageCRYAML, err := ioutil.ReadFile("testdata/expected-CR-image.yaml")
	require.NoError(t, err)

	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-image.yaml", CRD: expectedImageCRYAML})

	expectedReport := reportoutput.ReportOutput{}
	jsonData, err := io.ReadFile("testdata/expected-report-image.json")
	require.NoError(t, err)

	err = json.Unmarshal(jsonData, &expectedReport)
	require.NoError(t, err)

	testCases := []struct {
		name              string
		expectedManifests []transform.Manifest
		expectedReports   reportoutput.ReportOutput
	}{
		{
			name:              "transform image extraction",
			expectedManifests: expectedManifests,
			expectedReports:   expectedReport,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualManifestsChan := make(chan []transform.Manifest)
			actualReportsChan := make(chan reportoutput.ReportOutput)
			transform.FinalReportOutput = transform.Report{}

			// Override flush methods
			transform.ManifestOutputFlush = func(manifests []transform.Manifest) error {
				actualManifestsChan <- manifests
				return nil
			}
			transform.ReportOutputFlush = func(reports transform.Report) error {
				actualReportsChan <- reports.Report
				return nil
			}

			testExtraction, err := loadImageExtraction()
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

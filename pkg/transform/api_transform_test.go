package transform_test

import (
	"encoding/json"
	"fmt"
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

var loadAPIExtraction = func() transform.APIExtraction {
	file := "testdata/master_config-api.yaml"
	content, _ := ioutil.ReadFile(file)
	masterConfig, err := decode.MasterConfig(content)
	if err != nil {
		fmt.Printf("Error decoding file: %s\n", file)
	}
	var extraction transform.APIExtraction
	extraction.ServingInfo.BindAddress = masterConfig.ServingInfo.BindAddress
	extraction.ServingInfo.CertInfo = masterConfig.ServingInfo.CertInfo

	return extraction
}()

func TestAPIExtractionTransform(t *testing.T) {
	var expectedManifests []transform.Manifest

	expectedAPISecretCRYAML, err := ioutil.ReadFile("testdata/expected-CR-API-cert-secret.yaml")
	require.NoError(t, err)

	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-API-certificate-secret.yaml", CRD: expectedAPISecretCRYAML})

	expectedReport := reportoutput.ReportOutput{}
	jsonData, err := io.ReadFile("testdata/expected-report-api.json")
	require.NoError(t, err)

	err = json.Unmarshal(jsonData, &expectedReport)
	require.NoError(t, err)

	testCases := []struct {
		name              string
		expectedManifests []transform.Manifest
		expectedReports   reportoutput.ReportOutput
	}{
		{
			name:              "transform API extraction",
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

			testExtraction := loadAPIExtraction

			go func() {
				env.Config().Set("Manifests", true)
				env.Config().Set("Reporting", true)

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

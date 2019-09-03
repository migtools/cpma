package transform_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/decode"
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/transform"
	"github.com/fusor/cpma/pkg/transform/reportoutput"
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

	expectedAPISecretCRYAML, err := ioutil.ReadFile("testdata/expected-CR-APISecret.yaml")
	require.NoError(t, err)

	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-APISecret.yaml", CRD: expectedAPISecretCRYAML})

	expectedReport := reportoutput.ComponentReport{
		Component: "API",
	}

	expectedReport.Reports = append(expectedReport.Reports,
		reportoutput.Report{
			Name:       "API",
			Kind:       "Port",
			Supported:  false,
			Confidence: 0,
			Comment:    "The API Port for Openshift 4 is 6443 and is non-configurable. Your OCP 3 cluster is currently configured to use port 8443",
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
			name:              "transform API extraction",
			expectedManifests: expectedManifests,
			expectedReports:   expectedReportOutput,
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

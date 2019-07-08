package transform_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/decode"
	"github.com/fusor/cpma/pkg/transform"
	"github.com/stretchr/testify/assert"
)

var loadAPIExtraction = func() transform.APIExtraction {
	file := "testdata/master_config-api.yaml"
	content, _ := ioutil.ReadFile(file)
	masterConfig, err := decode.MasterConfig(content)
	if err != nil {
		fmt.Printf("Error decoding file: %s\n", file)
	}
	var extraction transform.APIExtraction
	extraction.HTTPServingInfo.BindAddress = masterConfig.ServingInfo.BindAddress

	return extraction
}()

func TestAPIExtractionTransform(t *testing.T) {
	expectedReport := transform.ComponentReport{
		Component: "API",
	}

	expectedReport.Reports = append(expectedReport.Reports,
		transform.Report{
			Name:       "API",
			Kind:       "Port",
			Supported:  false,
			Confidence: 0,
			Comment:    "The API Port for Openshift 4 is 6443 and is non-configurable. Your OCP 3 cluster is currently configured to use port 8443",
		})

	expectedReportOutput := transform.ReportOutput{
		ComponentReports: []transform.ComponentReport{expectedReport},
	}

	testCases := []struct {
		name            string
		expectedReports transform.ReportOutput
	}{
		{
			name:            "transform API extraction",
			expectedReports: expectedReportOutput,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualReportsChan := make(chan transform.ReportOutput)

			// Override flush method
			transform.ReportOutputFlush = func(reports transform.ReportOutput) error {
				actualReportsChan <- reports
				return nil
			}

			testExtraction := loadAPIExtraction

			go func() {
				transformOutput, err := testExtraction.Transform()
				if err != nil {
					t.Error(err)
				}
				for _, output := range transformOutput {
					output.Flush()
				}
			}()

			actualReports := <-actualReportsChan
			assert.Equal(t, actualReports, tc.expectedReports)
		})

	}
}

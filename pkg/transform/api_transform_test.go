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
		name            string
		expectedReports reportoutput.ReportOutput
	}{
		{
			name:            "transform API extraction",
			expectedReports: expectedReportOutput,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualReportsChan := make(chan reportoutput.ReportOutput)
			transform.FinalReportOutput = transform.Report{}

			// Override flush method
			transform.ReportOutputFlush = func(reports transform.Report) error {
				actualReportsChan <- reports.Report
				return nil
			}

			testExtraction := loadAPIExtraction

			go func() {
				env.Config().Set("Reporting", true)
				_, err := testExtraction.Transform()
				if err != nil {
					t.Error(err)
				}
				transform.FinalReportOutput.Flush()
			}()

			actualReports := <-actualReportsChan
			assert.Equal(t, actualReports.ComponentReports, tc.expectedReports.ComponentReports)
		})

	}
}

package transform_test

import (
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/decode"
	"github.com/fusor/cpma/pkg/transform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadAPIExtraction() (transform.APIExtraction, error) {
	// TODO: Something is broken here in a way that it's causing the translaters
	// to fail. Need some help with creating test identiy providers in a way
	// that won't crash the translator

	// Build example identity providers, this is straight copy pasted from
	// oauth test, IMO this loading of example identity providers should be
	// some shared test helper
	file := "testdata/master_config-api.yaml" // File copied into transform pkg testdata
	content, _ := ioutil.ReadFile(file)

	masterConfig, err := decode.MasterConfig(content)

	var extraction transform.APIExtraction
	extraction.HTTPServingInfo.BindAddress = masterConfig.ServingInfo.BindAddress

	return extraction, err
}

func TestAPIExtractionTransform(t *testing.T) {
	t.Parallel()

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

			testExtraction, err := loadAPIExtraction()
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

			actualReports := <-actualReportsChan
			assert.Equal(t, actualReports, tc.expectedReports)
		})

	}
}

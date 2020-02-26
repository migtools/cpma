package transform_test

import (
	"testing"

	"github.com/konveyor/cpma/pkg/env"
	"github.com/konveyor/cpma/pkg/transform"
	"github.com/konveyor/cpma/pkg/transform/reportoutput"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadDockerExtraction() (transform.DockerExtraction, error) {
	var extraction transform.DockerExtraction
	return extraction, nil
}

func TestDockerExtractionTransform(t *testing.T) {
	expectedReport := reportoutput.ComponentReport{
		Component: "Docker",
	}

	expectedReport.Reports = append(expectedReport.Reports,
		reportoutput.Report{
			Name:       "Docker",
			Kind:       "Container Runtime",
			Supported:  false,
			Confidence: 0,
			Comment:    "The Docker runtime has been replaced with CRI-O",
		})

	expectedReportOutput := reportoutput.ReportOutput{
		ComponentReports: []reportoutput.ComponentReport{expectedReport},
	}

	testCases := []struct {
		name            string
		expectedReports reportoutput.ReportOutput
	}{
		{
			name:            "transform crio extraction",
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

			testExtraction, err := loadDockerExtraction()
			require.NoError(t, err)

			go func() {
				env.Config().Set("Reporting", true)
				env.Config().Set("Manifests", true)
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

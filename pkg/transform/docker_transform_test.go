package transform_test

import (
	"testing"

	"github.com/fusor/cpma/pkg/transform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadDockerExtraction() (transform.DockerExtraction, error) {
	var extraction transform.DockerExtraction
	return extraction, nil
}

func TestDockerExtractionTransform(t *testing.T) {
	expectedReport := transform.ComponentReport{
		Component: "Docker",
	}

	expectedReport.Reports = append(expectedReport.Reports,
		transform.Report{
			Name:       "Docker",
			Kind:       "Container Runtime",
			Supported:  false,
			Confidence: 0,
			Comment:    "The Docker runtime has been replaced with CRI-O",
		})

	expectedReportOutput := transform.ReportOutput{
		ComponentReports: []transform.ComponentReport{expectedReport},
	}

	testCases := []struct {
		name            string
		expectedReports transform.ReportOutput
	}{
		{
			name:            "transform crio extraction",
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

			testExtraction, err := loadDockerExtraction()
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

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

func loadRegistriesExtraction() (transform.RegistriesExtraction, error) {
	// TODO: Something is broken here in a way that it's causing the translaters
	// to fail. Need some help with creating test identiy providers in a way
	// that won't crash the translator

	// Build example identity providers, this is straight copy pasted from
	// oauth test, IMO this loading of example identity providers should be
	// some shared test helper
	file := "testdata/registries.conf" // File copied into transform pkg testdata
	content, _ := ioutil.ReadFile(file)
	var extraction transform.RegistriesExtraction
	_, err := toml.Decode(string(content), &extraction)

	return extraction, err
}

func TestRegistriesExtractionTransform(t *testing.T) {
	var expectedManifests []transform.Manifest

	var expectedCrd transform.ImageCR
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "Image"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.Annotations = map[string]string{"release.openshift.io/create-only": "true"}
	expectedCrd.Spec.RegistrySources.BlockedRegistries = []string{"bad.guy"}
	expectedCrd.Spec.RegistrySources.InsecureRegistries = []string{"insecure.guy"}

	imageCRYAML, err := yaml.Marshal(&expectedCrd)
	require.NoError(t, err)

	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-registries.yaml", CRD: imageCRYAML})

	expectedReport := transform.ReportOutput{
		Component: "Registries",
	}
	expectedReport.Reports = append(expectedReport.Reports,
		transform.Report{
			Name:       "bad.guy",
			Kind:       "Blocked",
			Supported:  true,
			Confidence: 2,
		})
	expectedReport.Reports = append(expectedReport.Reports,
		transform.Report{
			Name:       "insecure.guy",
			Kind:       "Insecure",
			Supported:  true,
			Confidence: 2,
		})
	expectedReport.Reports = append(expectedReport.Reports,
		transform.Report{
			Name:       "search.guy",
			Kind:       "Search",
			Supported:  false,
			Confidence: 0,
			Comment:    "Search registries can not be configured in OCP 4",
		})

	testCases := []struct {
		name              string
		expectedManifests []transform.Manifest
		expectedReports   transform.ReportOutput
	}{
		{
			name:              "transform registries extraction",
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

			testExtraction, err := loadRegistriesExtraction()
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

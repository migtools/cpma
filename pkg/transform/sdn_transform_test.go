package transform_test

import (
	"testing"

	"github.com/fusor/cpma/pkg/transform"
	"github.com/fusor/cpma/pkg/transform/sdn"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestSDNExtractionTransform(t *testing.T) {
	var expectedManifests []transform.Manifest

	var expectedCrd sdn.NetworkCR
	expectedCrd.APIVersion = "operator.openshift.io/v1"
	expectedCrd.Kind = "Network"
	expectedCrd.Spec.ClusterNetworks = []sdn.ClusterNetwork{{HostPrefix: 23, CIDR: "10.128.0.0/14"}}
	expectedCrd.Spec.ServiceNetwork = "172.30.0.0/16"
	expectedCrd.Spec.DefaultNetwork.Type = "OpenShiftSDN"
	expectedCrd.Spec.DefaultNetwork.OpenshiftSDNConfig.Mode = "Subnet"

	networkCRYAML, err := yaml.Marshal(&expectedCrd)
	require.NoError(t, err)

	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-sdn.yaml", CRD: networkCRYAML})

	expectedReport := transform.ReportOutput{
		Component: "SDN",
	}
	expectedReport.Reports = append(expectedReport.Reports,
		transform.Report{
			Name:       "CIDR",
			Kind:       "ClusterNetwork",
			Supported:  true,
			Confidence: 1,
			Comment:    "Networks must be configured during installation, it's possible to use 10.128.0.0/14",
		})

	expectedReport.Reports = append(expectedReport.Reports,
		transform.Report{
			Name:       "HostSubnetLength",
			Kind:       "ClusterNetwork",
			Supported:  false,
			Confidence: 0,
			Comment:    "Networks must be configured during installation,\n hostSubnetLength was replaced with hostPrefix in OCP4, default value was set to 23",
		})
	expectedReport.Reports = append(expectedReport.Reports,
		transform.Report{
			Name:       "172.30.0.0/16",
			Kind:       "ServiceNetwork",
			Supported:  true,
			Confidence: 1,
			Comment:    "Networks must be configured during installation",
		})
	expectedReport.Reports = append(expectedReport.Reports,
		transform.Report{
			Name:       "",
			Kind:       "ExternalIPNetworkCIDRs",
			Supported:  false,
			Confidence: 0,
			Comment:    "Configuration of ExternalIPNetworkCIDRs is not supported in OCP4",
		})
	expectedReport.Reports = append(expectedReport.Reports,
		transform.Report{
			Name:       "",
			Kind:       "IngressIPNetworkCIDR",
			Supported:  false,
			Confidence: 0,
			Comment:    "Translation of this configuration is not supported, refer to ingress operator configuration for more information",
		})

	testCases := []struct {
		name              string
		expectedManifests []transform.Manifest
		expectedReports   transform.ReportOutput
	}{
		{
			name:              "transform sdn extraction",
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

			testExtraction, err := cpmatest.LoadSDNExtraction("testdata/sdn-test-master-config.yaml")
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

package sdn_test

import (
	"errors"

	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/transform/sdn"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	configv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfig(t *testing.T) {
	testExtraction, err := cpmatest.LoadSDNExtraction("testdata/test-master-config.yaml")
	require.NoError(t, err)

	testCases := []struct {
		name                           string
		expectedAPIVersion             string
		expectedKind                   string
		expectedCIDR                   string
		expectedHostPrefix             int
		expectedServiceNetwork         string
		expectedDefaultNetwork         string
		expectedOpenshiftSDNConfigMode string
	}{
		{
			expectedAPIVersion:             "operator.openshift.io/v1",
			expectedKind:                   "Network",
			expectedCIDR:                   "10.128.0.0/14",
			expectedHostPrefix:             23,
			expectedServiceNetwork:         "172.30.0.0/16",
			expectedDefaultNetwork:         "OpenShiftSDN",
			expectedOpenshiftSDNConfigMode: "Subnet",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			networkCR, err := sdn.Translate(testExtraction.MasterConfig)
			require.NoError(t, err)
			// Check if network CR was translated correctly
			assert.Equal(t, networkCR.APIVersion, "operator.openshift.io/v1")
			assert.Equal(t, networkCR.Kind, "Network")
			assert.Equal(t, networkCR.Spec.ClusterNetworks[0].CIDR, "10.128.0.0/14")
			assert.Equal(t, networkCR.Spec.ClusterNetworks[0].HostPrefix, 23)
			assert.Equal(t, networkCR.Spec.ServiceNetwork, "172.30.0.0/16")
			assert.Equal(t, networkCR.Spec.DefaultNetwork.Type, "OpenShiftSDN")
			assert.Equal(t, networkCR.Spec.DefaultNetwork.OpenshiftSDNConfig.Mode, "Subnet")

		})
	}
}

func TestSelectNetworkPlugin(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		output      string
		expectederr bool
	}{
		{
			name:        "translate multitenant",
			input:       "redhat/openshift-ovs-multitenant",
			output:      "Multitenant",
			expectederr: false,
		},
		{
			name:        "translate networkpolicy",
			input:       "redhat/openshift-ovs-networkpolicy",
			output:      "NetworkPolicy",
			expectederr: false,
		},
		{
			name:        "translate subnet",
			input:       "redhat/openshift-ovs-subnet",
			output:      "Subnet",
			expectederr: false,
		},
		{
			name:        "error on invalid plugin",
			input:       "123",
			output:      "error",
			expectederr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resPluginName, err := sdn.SelectNetworkPlugin(tc.input)

			if tc.expectederr {
				err := errors.New("Network plugin not supported")
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.output, resPluginName)
			}
		})
	}
}

func TestTransformClusterNetworks(t *testing.T) {
	testCases := []struct {
		name   string
		input  []configv1.ClusterNetworkEntry
		output []sdn.ClusterNetwork
	}{
		{
			name: "transform cluster networks",
			input: []configv1.ClusterNetworkEntry{
				configv1.ClusterNetworkEntry{CIDR: "10.128.0.0/14",
					HostSubnetLength: uint32(9),
				},
				configv1.ClusterNetworkEntry{CIDR: "10.127.0.0/14",
					HostSubnetLength: uint32(9),
				},
			},
			output: []sdn.ClusterNetwork{
				sdn.ClusterNetwork{
					CIDR:       "10.128.0.0/14",
					HostPrefix: 23,
				},
				sdn.ClusterNetwork{
					CIDR:       "10.127.0.0/14",
					HostPrefix: 23,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			translatedClusterNetworks := sdn.TranslateClusterNetworks(tc.input)
			assert.Equal(t, tc.output, translatedClusterNetworks)
		})
	}
}

func TestGenYAML(t *testing.T) {
	testExtraction, err := cpmatest.LoadSDNExtraction("testdata/test-master-config.yaml")
	require.NoError(t, err)

	networkCR, err := sdn.Translate(testExtraction.MasterConfig)
	require.NoError(t, err)

	expectedYaml, err := ioutil.ReadFile("testdata/expected-network-cr-master.yaml")
	require.NoError(t, err)

	testCases := []struct {
		name      string
		networkCR sdn.NetworkCR
		output    []byte
	}{
		{
			name:      "generate yaml for sdn",
			networkCR: networkCR,
			output:    expectedYaml,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			networkCRYAML, err := sdn.GenYAML(tc.networkCR)
			require.NoError(t, err)
			assert.Equal(t, tc.output, networkCRYAML)
		})
	}
}

func TestSDNValidation(t *testing.T) {
	testCases := []struct {
		name         string
		requireError bool
		inputFile    string
		expectedErr  error
	}{
		{
			name:         "validate sdn provider",
			requireError: false,
			inputFile:    "testdata/test-master-config.yaml",
		},
		{
			name:         "fail on empty service network CIDR in sdn provider",
			requireError: true,
			inputFile:    "testdata/test-empty-service-cidr-config.yaml",
			expectedErr:  errors.New("Service network CIDR can't be empty"),
		},
		{
			name:         "fail on invalid service network CIDR in sdn provider",
			requireError: true,
			inputFile:    "testdata/test-invalid-service-cidr-config.yaml",
			expectedErr:  errors.New("Not valid service network CIDR"),
		},
		{
			name:         "fail on empty cluster network in sdn provider",
			requireError: true,
			inputFile:    "testdata/test-empty-cluster-config.yaml",
			expectedErr:  errors.New("Cluster network must have at least 1 entry"),
		},
		{
			name:         "fail on empty cluster network CIDR in sdn provider",
			requireError: true,
			inputFile:    "testdata/test-empty-cluster-cidr-config.yaml",
			expectedErr:  errors.New("Cluster network CIDR can't be empty"),
		},
		{
			name:         "fail on invalid cluster network CIDR in sdn provider",
			requireError: true,
			inputFile:    "testdata/test-invalid-cluster-cidr-config.yaml",
			expectedErr:  errors.New("Not valid cluster network CIDR"),
		},
		{
			name:         "fail on empty plugin name in sdn provider",
			requireError: true,
			inputFile:    "testdata/test-empty-plugin-config.yaml",
			expectedErr:  errors.New("Plugin name can't be empty"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testExtraction, err := cpmatest.LoadSDNExtraction(tc.inputFile)
			require.NoError(t, err)

			err = testExtraction.Validate()

			if tc.requireError {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

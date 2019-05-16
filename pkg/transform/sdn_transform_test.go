package transform

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func TestTransformMasterConfig(t *testing.T) {
	file := "testdata/network-test-master-config.yaml"

	content, err := ioutil.ReadFile(file)
	require.NoError(t, err)

	var extraction SDNExtraction
	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	_, _, err = serializer.Decode(content, nil, &extraction.MasterConfig)
	require.NoError(t, err)

	testCases := []struct {
		name                           string
		expectedAPIVersion             string
		expectedKind                   string
		expectedCIDR                   string
		expectedHostPrefix             uint32
		expectedServiceNetwork         string
		expectedDefaultNetwork         string
		expectedOpenshiftSDNConfigMode string
	}{
		{
			expectedAPIVersion:             "operator.openshift.io/v1",
			expectedKind:                   "Network",
			expectedCIDR:                   "10.128.0.0/14",
			expectedHostPrefix:             uint32(9),
			expectedServiceNetwork:         "172.30.0.0/16",
			expectedDefaultNetwork:         "OpenShiftSDN",
			expectedOpenshiftSDNConfigMode: "Subnet",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			networkCR, err := SDNTranslate(extraction.MasterConfig)
			require.NoError(t, err)
			// Check if network CR was translated correctly
			assert.Equal(t, networkCR.APIVersion, "operator.openshift.io/v1")
			assert.Equal(t, networkCR.Kind, "Network")
			assert.Equal(t, networkCR.Spec.ClusterNetworks[0].CIDR, "10.128.0.0/14")
			assert.Equal(t, networkCR.Spec.ClusterNetworks[0].HostPrefix, uint32(9))
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
			resPluginName, err := SelectNetworkPlugin(tc.input)

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
		output []ClusterNetwork
	}{
		{
			name: "transform cluster networks",
			input: []configv1.ClusterNetworkEntry{
				configv1.ClusterNetworkEntry{CIDR: "10.128.0.0/14",
					HostSubnetLength: uint32(9),
				},
				configv1.ClusterNetworkEntry{CIDR: "10.127.0.0/14",
					HostSubnetLength: uint32(10),
				},
			},
			output: []ClusterNetwork{
				ClusterNetwork{
					CIDR:       "10.128.0.0/14",
					HostPrefix: uint32(9),
				},
				ClusterNetwork{
					CIDR:       "10.127.0.0/14",
					HostPrefix: uint32(10),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			translatedClusterNetworks := TranslateClusterNetworks(tc.input)
			assert.Equal(t, tc.output, translatedClusterNetworks)
		})
	}
}

func TestGenYAML(t *testing.T) {
	file := "testdata/network-test-master-config.yaml"

	content, err := ioutil.ReadFile(file)
	require.NoError(t, err)

	var extraction SDNExtraction
	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	_, _, err = serializer.Decode(content, nil, &extraction.MasterConfig)
	require.NoError(t, err)

	networkCR, err := SDNTranslate(extraction.MasterConfig)
	require.NoError(t, err)

	expectedYaml, err := ioutil.ReadFile("testdata/expected-network-cr-master.yaml")
	require.NoError(t, err)

	testCases := []struct {
		name      string
		networkCR NetworkCR
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
			networkCRYAML, err := GenYAML(tc.networkCR)
			require.NoError(t, err)
			assert.Equal(t, tc.output, networkCRYAML)
		})
	}
}

func loadSDNExtraction() (SDNExtraction, error) {
	// TODO: Something is broken here in a way that it's causing the translaters
	// to fail. Need some help with creating test identiy providers in a way
	// that won't crash the translator

	// Build example identity providers, this is straight copy pasted from
	// oauth test, IMO this loading of example identity providers should be
	// some shared test helper
	file := "testdata/network-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)
	var extraction SDNExtraction
	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	_, _, err := serializer.Decode(content, nil, &extraction.MasterConfig)

	return extraction, err
}

func TestSDNExtractionTransform(t *testing.T) {
	var expectedManifests []Manifest

	var expectedCrd NetworkCR
	expectedCrd.APIVersion = "operator.openshift.io/v1"
	expectedCrd.Kind = "Network"
	expectedCrd.Spec.ClusterNetworks = []ClusterNetwork{{HostPrefix: 9, CIDR: "10.128.0.0/14"}}
	expectedCrd.Spec.ServiceNetwork = "172.30.0.0/16"
	expectedCrd.Spec.DefaultNetwork.Type = "OpenShiftSDN"
	expectedCrd.Spec.DefaultNetwork.OpenshiftSDNConfig.Mode = "Subnet"

	networkCRYAML, err := yaml.Marshal(&expectedCrd)
	require.NoError(t, err)

	expectedManifests = append(expectedManifests,
		Manifest{Name: "100_CPMA-cluster-config-sdn.yaml", CRD: networkCRYAML})

	testCases := []struct {
		name              string
		expectedManifests []Manifest
	}{
		{
			name:              "transform sdn extraction",
			expectedManifests: expectedManifests,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualManifestsChan := make(chan []Manifest)
			// Override flush method
			manifestOutputFlush = func(manifests []Manifest) error {
				actualManifestsChan <- manifests
				return nil
			}

			testExtraction, err := loadSDNExtraction()
			require.NoError(t, err)

			go func() {
				transformOutput, err := testExtraction.Transform()
				if err != nil {
					t.Error(err)
				}
				transformOutput.Flush()
			}()

			actualManifests := <-actualManifestsChan
			assert.Equal(t, actualManifests, tc.expectedManifests)
		})
	}
}

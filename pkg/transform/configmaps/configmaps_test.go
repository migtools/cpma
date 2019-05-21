package configmaps

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenConfigMap(t *testing.T) {
	testCases := []struct {
		name          string
		configMapname string
		CAData        []byte
		namespace     string
		expected      ConfigMap
	}{
		{
			name:          "generate configmap",
			configMapname: "testname",
			CAData:        []byte("testdata"),
			namespace:     "openshift-config",
			expected: ConfigMap{
				APIVersion: APIVersion,
				Data: Data{
					CAData: "testdata",
				},
				Kind: Kind,
				Metadata: MetaData{
					Name:      "testname",
					Namespace: "openshift-config",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resConfigMap := GenConfigMap(tc.configMapname, tc.namespace, tc.CAData)
			assert.Equal(t, &tc.expected, resConfigMap)
		})
	}
}

func TestGenYaml(t *testing.T) {
	expectedYaml, err := ioutil.ReadFile("testdata/expected-configmap.yaml")
	require.NoError(t, err)

	testCases := []struct {
		name           string
		inputConfigMap ConfigMap
		expectedYaml   []byte
	}{
		{
			name: "generate yaml from configmap",
			inputConfigMap: ConfigMap{
				APIVersion: APIVersion,
				Data: Data{
					CAData: "testval: 123",
				},
				Kind: Kind,
				Metadata: MetaData{
					Name:      "testname",
					Namespace: "openshift-config",
				},
			},
			expectedYaml: expectedYaml,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			manifest, err := tc.inputConfigMap.GenYAML()
			require.NoError(t, err)
			assert.Equal(t, expectedYaml, manifest)
		})
	}
}

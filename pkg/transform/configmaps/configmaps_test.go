package configmaps_test

import (
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/transform"
	"github.com/fusor/cpma/pkg/transform/configmaps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenConfigMap(t *testing.T) {
	testCases := []struct {
		name          string
		configMapname string
		CAData        []byte
		namespace     string
		expected      configmaps.ConfigMap
	}{
		{
			name:          "generate configmap",
			configMapname: "testname",
			CAData:        []byte("testdata"),
			namespace:     "openshift-config",
			expected: configmaps.ConfigMap{
				APIVersion: configmaps.APIVersion,
				Data: configmaps.Data{
					CAData: "testdata",
				},
				Kind: configmaps.Kind,
				Metadata: configmaps.MetaData{
					Name:      "testname",
					Namespace: "openshift-config",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resConfigMap := configmaps.GenConfigMap(tc.configMapname, tc.namespace, tc.CAData)
			assert.Equal(t, &tc.expected, resConfigMap)
		})
	}
}

func TestGenYaml(t *testing.T) {
	expectedYaml, err := ioutil.ReadFile("testdata/expected-configmap.yaml")
	require.NoError(t, err)

	testCases := []struct {
		name           string
		inputConfigMap configmaps.ConfigMap
		expectedYaml   []byte
	}{
		{
			name: "generate yaml from configmap",
			inputConfigMap: configmaps.ConfigMap{
				APIVersion: configmaps.APIVersion,
				Data: configmaps.Data{
					CAData: "testval: 123",
				},
				Kind: configmaps.Kind,
				Metadata: configmaps.MetaData{
					Name:      "testname",
					Namespace: "openshift-config",
				},
			},
			expectedYaml: expectedYaml,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			manifest, err := transform.GenYAML(tc.inputConfigMap)
			require.NoError(t, err)
			assert.Equal(t, expectedYaml, manifest)
		})
	}
}

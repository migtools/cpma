package transform

import (
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/transform/configmaps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			manifest, err := GenYAML(tc.inputConfigMap)
			require.NoError(t, err)
			assert.Equal(t, expectedYaml, manifest)
		})
	}
}

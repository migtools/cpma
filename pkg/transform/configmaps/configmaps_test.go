package configmaps

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

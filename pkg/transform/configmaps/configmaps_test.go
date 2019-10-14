package configmaps

import (
	"testing"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGenConfigMap(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		configMapname string
		CAData        []byte
		namespace     string
		expected      corev1.ConfigMap
	}{
		{
			name:          "generate configmap",
			configMapname: "testname",
			CAData:        []byte("testdata"),
			namespace:     "openshift-config",
			expected: corev1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "ConfigMap",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testname",
					Namespace: "openshift-config",
				},
				Data: map[string]string{
					"ca.crt": "testdata",
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

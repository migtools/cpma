package secrets

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGenSecret(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name            string
		inputSecretName string
		inputSecretFile string
		inputSecretType string
		expected        corev1.Secret
		expectederr     bool
	}{
		{
			name:            "generate htpasswd secret",
			inputSecretName: "htpasswd-test",
			inputSecretFile: "testfile1",
			inputSecretType: "htpasswd",
			expected: corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Data: map[string][]byte{
					"htpasswd": []byte("testfile1"),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "htpasswd-test",
					Namespace: "openshift-config",
				},
				Type: "Opaque",
			},
			expectederr: false,
		},
		{
			name:            "generate keystone secret",
			inputSecretName: "keystone-test",
			inputSecretFile: "testfile2",
			inputSecretType: "keystone",
			expected: corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Data: map[string][]byte{
					"keystone": []byte("testfile2"),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "keystone-test",
					Namespace: "openshift-config",
				},
				Type: "Opaque",
			},
			expectederr: false,
		},
		{
			name:            "generate basic auth secret",
			inputSecretName: "basicauth-test",
			inputSecretFile: "testfile3",
			inputSecretType: "testname",
			expected: corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Data: map[string][]byte{
					"testname": []byte("testfile3"),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "basicauth-test",
					Namespace: "openshift-config",
				},
				Type: "Opaque",
			},
			expectederr: false,
		},
		{
			name:            "generate litetal secret",
			inputSecretName: "literal-secret",
			inputSecretFile: "some-value",
			inputSecretType: "clientSecret",
			expected: corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Data: map[string][]byte{
					"clientSecret": []byte("some-value"),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "literal-secret",
					Namespace: "openshift-config",
				},
				Type: "Opaque",
			},
			expectederr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resSecret, err := Opaque(tc.inputSecretName, []byte(tc.inputSecretFile), "openshift-config", tc.inputSecretType)
			if tc.expectederr {
				err := errors.New("Not valid secret type " + "notvalidtype")
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, &tc.expected, resSecret)
			}
		})
	}
}

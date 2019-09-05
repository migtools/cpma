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
		inputSecretType SecretType
		expected        corev1.Secret
		expectederr     bool
	}{
		{
			name:            "generate htpasswd secret",
			inputSecretName: "htpasswd-test",
			inputSecretFile: "testfile1",
			inputSecretType: HtpasswdSecretType,
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
			inputSecretType: KeystoneSecretType,
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
			inputSecretType: BasicAuthSecretType,
			expected: corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Data: map[string][]byte{
					"basicAuth": []byte("testfile3"),
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
			inputSecretType: LiteralSecretType,
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
		{
			name:            "fail generating invalid secret",
			inputSecretName: "notvalid-secret",
			inputSecretFile: "some-value",
			inputSecretType: 42, // Unknown secret type value
			expectederr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resSecret, err := GenSecret(tc.inputSecretName, tc.inputSecretFile, "openshift-config", tc.inputSecretType)
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

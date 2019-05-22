package secrets

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenSecret(t *testing.T) {
	testCases := []struct {
		name            string
		inputSecretName string
		inputSecretFile string
		inputSecretType SecretType
		expected        Secret
		expectederr     bool
	}{
		{
			name:            "generate htpasswd secret",
			inputSecretName: "htpasswd-test",
			inputSecretFile: "testfile1",
			inputSecretType: HtpasswdSecretType,
			expected: Secret{
				APIVersion: APIVersion,
				Data:       HTPasswdFileSecret{HTPasswd: "testfile1"},
				Kind:       "Secret",
				Type:       "Opaque",
				Metadata: MetaData{
					Name:      "htpasswd-test",
					Namespace: "openshift-config",
				},
			},
			expectederr: false,
		},
		{
			name:            "generate keystone secret",
			inputSecretName: "keystone-test",
			inputSecretFile: "testfile2",
			inputSecretType: KeystoneSecretType,
			expected: Secret{
				APIVersion: APIVersion,
				Data:       KeystoneFileSecret{Keystone: "testfile2"},
				Kind:       "Secret",
				Type:       "Opaque",
				Metadata: MetaData{
					Name:      "keystone-test",
					Namespace: "openshift-config",
				},
			},
			expectederr: false,
		},
		{
			name:            "generate basic auth secret",
			inputSecretName: "basicauth-test",
			inputSecretFile: "testfile3",
			inputSecretType: BasicAuthSecretType,
			expected: Secret{
				APIVersion: APIVersion,
				Data:       BasicAuthFileSecret{BasicAuth: "testfile3"},
				Kind:       "Secret",
				Type:       "Opaque",
				Metadata: MetaData{
					Name:      "basicauth-test",
					Namespace: "openshift-config",
				},
			},
			expectederr: false,
		},
		{
			name:            "generate litetal secret",
			inputSecretName: "literal-secret",
			inputSecretFile: "some-value",
			inputSecretType: LiteralSecretType,
			expected: Secret{
				APIVersion: APIVersion,
				Data:       LiteralSecret{ClientSecret: "some-value"},
				Kind:       "Secret",
				Type:       "Opaque",
				Metadata: MetaData{
					Name:      "literal-secret",
					Namespace: "openshift-config",
				},
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

func TestGenYaml(t *testing.T) {
	expectedYaml, err := ioutil.ReadFile("testdata/expected-secret.yaml")
	require.NoError(t, err)

	testCases := []struct {
		name         string
		inputSecret  Secret
		expectedYaml []byte
	}{
		{
			name: "generate yaml from secret",
			inputSecret: Secret{
				APIVersion: APIVersion,
				Data:       LiteralSecret{ClientSecret: "some-value"},
				Kind:       "Secret",
				Type:       "Opaque",
				Metadata: MetaData{
					Name:      "literal-secret",
					Namespace: "openshift-config",
				},
			},
			expectedYaml: expectedYaml,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			manifest, err := tc.inputSecret.GenYAML()
			require.NoError(t, err)
			assert.Equal(t, expectedYaml, manifest)
		})
	}
}

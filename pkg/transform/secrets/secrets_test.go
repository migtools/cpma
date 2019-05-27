package secrets_test

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/transform/secrets"

	"github.com/fusor/cpma/pkg/transform"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenSecret(t *testing.T) {
	testCases := []struct {
		name            string
		inputSecretName string
		inputSecretFile string
		inputSecretType secrets.SecretType
		expected        secrets.Secret
		expectederr     bool
	}{
		{
			name:            "generate htpasswd secret",
			inputSecretName: "htpasswd-test",
			inputSecretFile: "testfile1",
			inputSecretType: secrets.HtpasswdSecretType,
			expected: secrets.Secret{
				APIVersion: secrets.APIVersion,
				Data:       secrets.HTPasswdFileSecret{HTPasswd: "testfile1"},
				Kind:       "Secret",
				Type:       "Opaque",
				Metadata: secrets.MetaData{
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
			inputSecretType: secrets.KeystoneSecretType,
			expected: secrets.Secret{
				APIVersion: secrets.APIVersion,
				Data:       secrets.KeystoneFileSecret{Keystone: "testfile2"},
				Kind:       "Secret",
				Type:       "Opaque",
				Metadata: secrets.MetaData{
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
			inputSecretType: secrets.BasicAuthSecretType,
			expected: secrets.Secret{
				APIVersion: secrets.APIVersion,
				Data:       secrets.BasicAuthFileSecret{BasicAuth: "testfile3"},
				Kind:       "Secret",
				Type:       "Opaque",
				Metadata: secrets.MetaData{
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
			inputSecretType: secrets.LiteralSecretType,
			expected: secrets.Secret{
				APIVersion: secrets.APIVersion,
				Data:       secrets.LiteralSecret{ClientSecret: "some-value"},
				Kind:       "Secret",
				Type:       "Opaque",
				Metadata: secrets.MetaData{
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
			resSecret, err := secrets.GenSecret(tc.inputSecretName, tc.inputSecretFile, "openshift-config", tc.inputSecretType)
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
		inputSecret  secrets.Secret
		expectedYaml []byte
	}{
		{
			name: "generate yaml from secret",
			inputSecret: secrets.Secret{
				APIVersion: secrets.APIVersion,
				Data:       secrets.LiteralSecret{ClientSecret: "some-value"},
				Kind:       "Secret",
				Type:       "Opaque",
				Metadata: secrets.MetaData{
					Name:      "literal-secret",
					Namespace: "openshift-config",
				},
			},
			expectedYaml: expectedYaml,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			manifest, err := transform.GenYAML(tc.inputSecret)
			require.NoError(t, err)
			assert.Equal(t, expectedYaml, manifest)
		})
	}
}

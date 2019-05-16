package secrets

import (
	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenSecretFileHtpasswd(t *testing.T) {
	htpasswdFile := "testfile1"
	encoded := base64.StdEncoding.EncodeToString([]byte(htpasswdFile))
	resSecret, err := GenSecret("htpasswd-test", encoded, "openshift-config", "htpasswd")
	require.NoError(t, err)

	var data = HTPasswdFileSecret{HTPasswd: encoded}
	var expectedSecret = Secret{
		APIVersion: APIVersion,
		Data:       data,
		Kind:       "Secret",
		Type:       "Opaque",
		Metadata: MetaData{
			Name:      "htpasswd-test",
			Namespace: "openshift-config",
		},
	}

	assert.Equal(t, &expectedSecret, resSecret)
}
func TestGenSecretFileKeystone(t *testing.T) {
	keystoneFile := "testfile2"
	encoded := base64.StdEncoding.EncodeToString([]byte(keystoneFile))
	resSecret, err := GenSecret("keystone-test", encoded, "openshift-config", "keystone")
	require.NoError(t, err)

	var data = KeystoneFileSecret{Keystone: encoded}
	var expectedSecret = Secret{
		APIVersion: APIVersion,
		Data:       data,
		Kind:       "Secret",
		Type:       "Opaque",
		Metadata: MetaData{
			Name:      "keystone-test",
			Namespace: "openshift-config",
		},
	}

	assert.Equal(t, &expectedSecret, resSecret)
}
func TestGenSecretFileBasicAuth(t *testing.T) {
	basicAuth := "testfile2"
	encoded := base64.StdEncoding.EncodeToString([]byte(basicAuth))
	resSecret, err := GenSecret("keystone-test", encoded, "openshift-config", "basicauth")
	require.NoError(t, err)

	var data = BasicAuthFileSecret{BasicAuth: encoded}
	var expectedSecret = Secret{
		APIVersion: APIVersion,
		Data:       data,
		Kind:       "Secret",
		Type:       "Opaque",
		Metadata: MetaData{
			Name:      "keystone-test",
			Namespace: "openshift-config",
		},
	}

	assert.Equal(t, &expectedSecret, resSecret)
}
func TestGenSecretLiteral(t *testing.T) {
	resSecret, err := GenSecret("literal-secret", "some-value", "openshift-config", "literal")
	require.NoError(t, err)

	var data = LiteralSecret{ClientSecret: "some-value"}
	var expectedSecret = Secret{
		APIVersion: APIVersion,
		Data:       data,
		Kind:       "Secret",
		Type:       "Opaque",
		Metadata: MetaData{
			Name:      "literal-secret",
			Namespace: "openshift-config",
		},
	}

	assert.Equal(t, &expectedSecret, resSecret)
}

func TestGenYaml(t *testing.T) {
	var data = LiteralSecret{ClientSecret: "some-value"}
	var secret = Secret{
		APIVersion: APIVersion,
		Data:       data,
		Kind:       "Secret",
		Type:       "Opaque",
		Metadata: MetaData{
			Name:      "literal-secret",
			Namespace: "openshift-config",
		},
	}

	manifest, err := secret.GenYAML()
	require.NoError(t, err)

	expectedYaml, err := ioutil.ReadFile("testdata/expected-secret.yaml")
	require.NoError(t, err)

	assert.Equal(t, expectedYaml, manifest)
}

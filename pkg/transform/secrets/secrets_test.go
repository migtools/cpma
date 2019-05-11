package secrets

import (
	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenSecretFileHtpasswd(t *testing.T) {
	htpasswdFile := "testfile1"
	encoded := base64.StdEncoding.EncodeToString([]byte(htpasswdFile))
	resSecret := GenSecret("htpasswd-test", encoded, "openshift-config", "htpasswd")

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
	resSecret := GenSecret("keystone-test", encoded, "openshift-config", "keystone")

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
	resSecret := GenSecret("keystone-test", encoded, "openshift-config", "basicauth")

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
	resSecret := GenSecret("literal-secret", "some-value", "openshift-config", "literal")

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

	manifest := secret.GenYAML()
	expectedYaml, _ := ioutil.ReadFile("testdata/expected-secret.yaml")
	assert.Equal(t, expectedYaml, manifest)
}

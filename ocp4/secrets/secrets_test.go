package secrets

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenSecretFileHtpasswd(t *testing.T) {
	htpasswdFile := "testfile1"
	encoded := base64.StdEncoding.EncodeToString([]byte(htpasswdFile))
	resSecret := GenSecretFile("htpasswd-test", encoded, "openshift-config", "htpasswd")

	var data = HTPasswdFileSecret{HTPasswd: encoded}
	var expectedSecret = Secret{
		APIVersion: APIVersion,
		Data:       data,
		Kind:       "Secret",
		Type:       "Opaque",
		MetaData: MetaData{
			Name:      "htpasswd-test",
			Namespace: "openshift-config",
		},
	}

	assert.Equal(t, &expectedSecret, resSecret)
}
func TestGenSecretFileKeystone(t *testing.T) {
	keystoneFile := "testfile2"
	encoded := base64.StdEncoding.EncodeToString([]byte(keystoneFile))
	resSecret := GenSecretFile("keystone-test", encoded, "openshift-config", "keystone")

	var data = KeystoneFileSecret{Keystone: encoded}
	var expectedSecret = Secret{
		APIVersion: APIVersion,
		Data:       data,
		Kind:       "Secret",
		Type:       "Opaque",
		MetaData: MetaData{
			Name:      "keystone-test",
			Namespace: "openshift-config",
		},
	}

	assert.Equal(t, &expectedSecret, resSecret)
}
func TestGenSecretLiteral(t *testing.T) {
	resSecret := GenSecretLiteral("literal-secret", "some-value", "openshift-config")

	var data = LiteralSecret{ClientSecret: "some-value"}
	var expectedSecret = Secret{
		APIVersion: APIVersion,
		Data:       data,
		Kind:       "Secret",
		Type:       "Opaque",
		MetaData: MetaData{
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
		MetaData: MetaData{
			Name:      "literal-secret",
			Namespace: "openshift-config",
		},
	}

	yaml := secret.GenYAML()
	expectedYaml := "apiVersion: v1\nkind: Secret\ntype: Opaque\nmetaData:\n  name: literal-secret\n  namespace: openshift-config\ndata:\n  clientSecret: some-value\n"

	assert.Equal(t, expectedYaml, yaml)
}

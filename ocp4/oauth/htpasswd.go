package oauth

import (
	"encoding/base64"

	"github.com/fusor/cpma/ocp4/secrets"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

type identityProviderHTPasswd struct {
	identityProviderCommon `yaml:",inline"`
	HTPasswd               struct {
		FileData struct {
			Name string `yaml:"name"`
		} `yaml:"fileData"`
	} `yaml:"htpasswd"`
}

func buildHTPasswdIP(serializer *json.Serializer, p configv1.IdentityProvider) (identityProviderHTPasswd, secrets.Secret) {
	var idP identityProviderHTPasswd
	var htpasswd configv1.HTPasswdPasswordIdentityProvider
	_, _, _ = serializer.Decode(p.Provider.Raw, nil, &htpasswd)

	idP.Name = p.Name
	idP.Type = "HTPasswd"
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.HTPasswd.FileData.Name = htpasswd.File

	secretName := p.Name + "-secret"
	idP.HTPasswd.FileData.Name = secretName
	// Retrieve file
	//htpasswdFile := Fetch_File(htpasswd.File)
	htpasswdFile := "This is pretend content"
	encoded := base64.StdEncoding.EncodeToString([]byte(htpasswdFile))
	secret := secrets.GenSecretFile(secretName, encoded, "openshift-config", "htpasswd")

	return idP, *secret
}

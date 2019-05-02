package oauth

import (
	"encoding/base64"

	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/fusor/cpma/pkg/ocp4/secrets"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

type identityProviderHTPasswd struct {
	identityProviderCommon `yaml:",inline"`
	HTPasswd               struct {
		FileData struct {
			Name string `yaml:"name"`
		} `yaml:"fileData"`
	} `yaml:"htpasswd"`
}

func buildHTPasswdIP(serializer *json.Serializer, p ocp3.IdentityProvider) (identityProviderHTPasswd, secrets.Secret) {
	var idP identityProviderHTPasswd

	idP.Name = p.Name
	idP.Type = "HTPasswd"
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.HTPasswd.FileData.Name = p.HTFileName

	secretName := p.Name + "-secret"
	idP.HTPasswd.FileData.Name = secretName
	encoded := base64.StdEncoding.EncodeToString(p.HTFileData)
	secret := secrets.GenSecret(secretName, encoded, "openshift-config", "htpasswd")

	return idP, *secret
}

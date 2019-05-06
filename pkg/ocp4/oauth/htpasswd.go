package oauth

import (
	"encoding/base64"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	"github.com/fusor/cpma/env"
	"github.com/fusor/cpma/pkg/ocp4/secrets"
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

	host := env.Config().GetString("Source")
	src := filepath.Join(htpasswd.File)
	dst := filepath.Join(env.Config().GetString("OutputDir"), host, htpasswd.File)
	f := GetFile(host, src, dst)

	encoded := base64.StdEncoding.EncodeToString(f)
	secret := secrets.GenSecret(secretName, encoded, "openshift-config", "htpasswd")

	return idP, *secret
}

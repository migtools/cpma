package oauth

import (
	"encoding/base64"
	"path/filepath"

	"github.com/fusor/cpma/env"
	"github.com/fusor/cpma/ocp3"
	"github.com/fusor/cpma/ocp4/secrets"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

type identityProviderHTPasswd struct {
	Name          string `yaml:"name"`
	Challenge     bool   `yaml:"challenge"`
	Login         bool   `yaml:"login"`
	MappingMethod string `yaml:"mappingMethod"`
	Type          string `yaml:"type"`
	HTPasswd      struct {
		FileData struct {
			Name string `yaml:"name"`
		} `yaml:"fileData"`
	} `yaml:"htpasswd"`
}

func buildHTPasswdIP(serializer *json.Serializer, p configv1.IdentityProvider) (identityProviderHTPasswd, secrets.Secret) {
	var idP identityProviderHTPasswd
	var htpasswd configv1.HTPasswdPasswordIdentityProvider
	_, _, _ = serializer.Decode(p.Provider.Raw, nil, &htpasswd)

	idP.Type = "HTPasswd"
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.HTPasswd.FileData.Name = htpasswd.File

	secretName := p.Name + "-secret"
	idP.HTPasswd.FileData.Name = secretName

	outputDir := env.Config().GetString("OutputDir")
	srcPath := filepath.Join(htpasswd.File)
	dstPath := filepath.Join(outputDir, htpasswd.File)
	ocp3.FetchFile(srcPath, dstPath)

	encoded := base64.StdEncoding.EncodeToString([]byte(dstPath))
	secret := secrets.GenSecretFile(secretName, encoded, "openshift-config")

	return idP, *secret
}

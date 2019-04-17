package oauth

import (
	configv1 "github.com/openshift/api/legacyconfig/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
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

func buildHTPasswdIP(serializer *json.Serializer, p configv1.IdentityProvider) identityProviderHTPasswd {
	var idP identityProviderHTPasswd
	var htpasswd configv1.HTPasswdPasswordIdentityProvider
	_, _, _ = serializer.Decode(p.Provider.Raw, nil, &htpasswd)

	idP.Type = "HTPasswd"
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.HTPasswd.FileData.Name = htpasswd.File
	return idP
}

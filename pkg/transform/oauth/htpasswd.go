package oauth

import (
	"encoding/base64"

	"github.com/fusor/cpma/pkg/transform/secrets"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

// IdentityProviderHTPasswd is a htpasswd specific identity provider
type IdentityProviderHTPasswd struct {
	identityProviderCommon `json:",inline"`
	HTPasswd               `json:"htpasswd"`
}

// HTPasswd contains htpasswd FileData
type HTPasswd struct {
	FileData FileData `json:"fileData"`
}

// FileData from htpasswd file
type FileData struct {
	Name string `json:"name"`
}

func buildHTPasswdIP(serializer *json.Serializer, p IdentityProvider) (*IdentityProviderHTPasswd, *secrets.Secret, error) {
	var (
		err      error
		idP      = &IdentityProviderHTPasswd{}
		secret   *secrets.Secret
		htpasswd legacyconfigv1.HTPasswdPasswordIdentityProvider
	)

	if _, _, err = serializer.Decode(p.Provider.Raw, nil, &htpasswd); err != nil {
		return nil, nil, errors.Wrap(err, "Failed to decode htpasswd, see error")
	}

	idP.Name = p.Name
	idP.Type = "HTPasswd"
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.HTPasswd.FileData.Name = htpasswd.File

	secretName := p.Name + "-secret"
	idP.HTPasswd.FileData.Name = secretName

	encoded := base64.StdEncoding.EncodeToString(p.HTFileData)

	secret, err = secrets.GenSecret(secretName, encoded, OAuthNamespace, secrets.HtpasswdSecretType)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to generate secret for htpasswd, see error")
	}

	return idP, secret, nil
}

func validateHTPasswdProvider(serializer *json.Serializer, p IdentityProvider) error {
	var htpasswd legacyconfigv1.HTPasswdPasswordIdentityProvider

	if _, _, err := serializer.Decode(p.Provider.Raw, nil, &htpasswd); err != nil {
		return errors.Wrap(err, "Failed to decode htpasswd, see error")
	}

	if p.Name == "" {
		return errors.New("Name can't be empty")
	}

	if err := validateMappingMethod(p.MappingMethod); err != nil {
		return err
	}

	if htpasswd.File == "" {
		return errors.New("File can't be empty")
	}

	return nil
}

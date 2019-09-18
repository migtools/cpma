package oauth

import (
	"github.com/fusor/cpma/pkg/transform/secrets"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func buildHTPasswdIP(serializer *json.Serializer, p IdentityProvider) (*ProviderResources, error) {
	var (
		err             error
		idP             = &configv1.IdentityProvider{}
		providerSecrets []*corev1.Secret
		htpasswd        legacyconfigv1.HTPasswdPasswordIdentityProvider
	)

	if _, _, err = serializer.Decode(p.Provider.Raw, nil, &htpasswd); err != nil {
		return nil, errors.Wrap(err, "Failed to decode htpasswd, see error")
	}

	idP.Name = p.Name
	idP.Type = "HTPasswd"
	idP.MappingMethod = configv1.MappingMethodType(p.MappingMethod)

	secretName := "htpasswd-secret"
	idP.HTPasswd = &configv1.HTPasswdIdentityProvider{}
	idP.HTPasswd.FileData.Name = secretName

	secret, err := secrets.GenSecret(secretName, p.HTFileData, OAuthNamespace, secrets.HtpasswdSecretType)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate secret for htpasswd, see error")
	}
	providerSecrets = append(providerSecrets, secret)

	return &ProviderResources{
		IDP:     idP,
		Secrets: providerSecrets,
	}, nil
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

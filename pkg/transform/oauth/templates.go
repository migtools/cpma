package oauth

import (
	"encoding/base64"

	"github.com/fusor/cpma/pkg/transform/secrets"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
)

const (
	loginSecret             = "templates-login-secret"
	errorSecret             = "templates-error-secret"
	providerSelectionSecret = "templates-providerselect-secret"
)

func translateTemplates(templates legacyconfigv1.OAuthTemplates) (*configv1.OAuthTemplates, []*secrets.Secret, error) {
	var templateSecrets []*secrets.Secret

	translatedTemplates := &configv1.OAuthTemplates{}

	if templates.Login != "" {
		encoded := base64.StdEncoding.EncodeToString([]byte(templates.Login))
		secret, err := secrets.GenSecret(loginSecret, encoded, OAuthNamespace, secrets.LiteralSecretType)
		if err != nil {
			return nil, nil, err
		}
		templateSecrets = append(templateSecrets, secret)
		translatedTemplates.Login = configv1.SecretNameReference{Name: loginSecret}
	}

	if templates.Error != "" {
		encoded := base64.StdEncoding.EncodeToString([]byte(templates.Error))
		secret, err := secrets.GenSecret(errorSecret, encoded, OAuthNamespace, secrets.LiteralSecretType)
		if err != nil {
			return nil, nil, err
		}
		templateSecrets = append(templateSecrets, secret)
		translatedTemplates.Error = configv1.SecretNameReference{Name: errorSecret}
	}

	if templates.ProviderSelection != "" {
		encoded := base64.StdEncoding.EncodeToString([]byte(templates.ProviderSelection))
		secret, err := secrets.GenSecret(providerSelectionSecret, encoded, OAuthNamespace, secrets.LiteralSecretType)
		if err != nil {
			return nil, nil, err
		}
		templateSecrets = append(templateSecrets, secret)
		translatedTemplates.ProviderSelection = configv1.SecretNameReference{Name: providerSelectionSecret}
	}

	return translatedTemplates, templateSecrets, nil
}

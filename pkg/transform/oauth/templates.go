package oauth

import (
	"github.com/fusor/cpma/pkg/transform/secrets"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"

	corev1 "k8s.io/api/core/v1"
)

const (
	loginSecret             = "templates-login-secret"
	errorSecret             = "templates-error-secret"
	providerSelectionSecret = "templates-providerselect-secret"
)

func translateTemplates(templates legacyconfigv1.OAuthTemplates) (*configv1.OAuthTemplates, []*corev1.Secret, error) {
	var templateSecrets []*corev1.Secret

	translatedTemplates := &configv1.OAuthTemplates{}

	if templates.Login != "" {
		secret, err := secrets.GenSecret(loginSecret, []byte(templates.Login), OAuthNamespace, secrets.LiteralSecretType)
		if err != nil {
			return nil, nil, err
		}
		templateSecrets = append(templateSecrets, secret)
		translatedTemplates.Login = configv1.SecretNameReference{Name: loginSecret}
	}

	if templates.Error != "" {
		secret, err := secrets.GenSecret(errorSecret, []byte(templates.Error), OAuthNamespace, secrets.LiteralSecretType)
		if err != nil {
			return nil, nil, err
		}
		templateSecrets = append(templateSecrets, secret)
		translatedTemplates.Error = configv1.SecretNameReference{Name: errorSecret}
	}

	if templates.ProviderSelection != "" {
		secret, err := secrets.GenSecret(providerSelectionSecret, []byte(templates.ProviderSelection), OAuthNamespace, secrets.LiteralSecretType)
		if err != nil {
			return nil, nil, err
		}
		templateSecrets = append(templateSecrets, secret)
		translatedTemplates.ProviderSelection = configv1.SecretNameReference{Name: providerSelectionSecret}
	}

	return translatedTemplates, templateSecrets, nil
}

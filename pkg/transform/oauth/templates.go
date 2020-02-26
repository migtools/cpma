package oauth

import (
	"github.com/konveyor/cpma/pkg/transform/secrets"
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
		secret, err := secrets.Opaque(loginSecret, []byte(templates.Login), OAuthNamespace, "clientSecret")
		if err != nil {
			return nil, nil, err
		}
		templateSecrets = append(templateSecrets, secret)
		translatedTemplates.Login = configv1.SecretNameReference{Name: loginSecret}
	}

	if templates.Error != "" {
		secret, err := secrets.Opaque(errorSecret, []byte(templates.Error), OAuthNamespace, "clientSecret")
		if err != nil {
			return nil, nil, err
		}
		templateSecrets = append(templateSecrets, secret)
		translatedTemplates.Error = configv1.SecretNameReference{Name: errorSecret}
	}

	if templates.ProviderSelection != "" {
		secret, err := secrets.Opaque(providerSelectionSecret, []byte(templates.ProviderSelection), OAuthNamespace, "clientSecret")
		if err != nil {
			return nil, nil, err
		}
		templateSecrets = append(templateSecrets, secret)
		translatedTemplates.ProviderSelection = configv1.SecretNameReference{Name: providerSelectionSecret}
	}

	return translatedTemplates, templateSecrets, nil
}

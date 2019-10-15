package oauth

import (
	"path/filepath"

	configv1 "github.com/openshift/api/config/v1"

	"github.com/fusor/cpma/pkg/transform/configmaps"
	"github.com/fusor/cpma/pkg/transform/secrets"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func buildBasicAuthIP(serializer *json.Serializer, p IdentityProvider) (*ProviderResources, error) {
	var (
		err                error
		idP                = &configv1.IdentityProvider{}
		basicAuth          legacyconfigv1.BasicAuthPasswordIdentityProvider
		providerSecrets    []*corev1.Secret
		providerConfigMaps []*corev1.ConfigMap
	)

	if _, _, err = serializer.Decode(p.Provider.Raw, nil, &basicAuth); err != nil {
		return nil, errors.Wrap(err, "Failed to decode basic auth, see error")
	}

	idP.Type = "BasicAuth"
	idP.Name = p.Name
	idP.MappingMethod = configv1.MappingMethodType(p.MappingMethod)
	idP.BasicAuth = &configv1.BasicAuthIdentityProvider{}
	idP.BasicAuth.URL = basicAuth.URL
	if basicAuth.CA != "" {
		caConfigmap := configmaps.GenConfigMap("basicauth-configmap", OAuthNamespace, p.CAData)
		idP.BasicAuth.CA = configv1.ConfigMapNameReference{Name: caConfigmap.ObjectMeta.Name}
		providerConfigMaps = append(providerConfigMaps, caConfigmap)
	}

	if basicAuth.CertFile != "" {
		certSecretName := "basicauth-client-cert-secret"
		idP.BasicAuth.TLSClientCert.Name = certSecretName

		certSecret, err := secrets.Opaque(certSecretName, p.CrtData, OAuthNamespace, filepath.Base(basicAuth.CertFile))
		if err != nil {
			return nil, errors.Wrap(err, "Failed to generate cert secret for basic auth, see error")
		}
		providerSecrets = append(providerSecrets, certSecret)

		keySecretName := "basicauth-client-key-secret"
		idP.BasicAuth.TLSClientKey.Name = keySecretName

		keySecret, err := secrets.Opaque(keySecretName, p.KeyData, OAuthNamespace, filepath.Base(basicAuth.KeyFile))
		if err != nil {
			return nil, errors.Wrap(err, "Failed to generate key secret for basic auth, see error")
		}
		providerSecrets = append(providerSecrets, keySecret)
	}

	return &ProviderResources{
		IDP:        idP,
		Secrets:    providerSecrets,
		ConfigMaps: providerConfigMaps,
	}, nil
}

func validateBasicAuthProvider(serializer *json.Serializer, p IdentityProvider) error {
	var basicAuth legacyconfigv1.BasicAuthPasswordIdentityProvider

	if _, _, err := serializer.Decode(p.Provider.Raw, nil, &basicAuth); err != nil {
		return errors.Wrap(err, "Failed to decode basic auth, see error")
	}

	if p.Name == "" {
		return errors.New("Name can't be empty")
	}

	if err := validateMappingMethod(p.MappingMethod); err != nil {
		return err
	}

	if basicAuth.URL == "" {
		return errors.New("URL can't be empty")
	}

	if basicAuth.CertFile != "" && basicAuth.KeyFile == "" {
		return errors.New("Key file can't be empty if cert file is specified")
	}

	return nil
}

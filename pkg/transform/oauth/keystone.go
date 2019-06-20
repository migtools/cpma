package oauth

import (
	"encoding/base64"

	"github.com/fusor/cpma/pkg/transform/configmaps"
	"github.com/fusor/cpma/pkg/transform/secrets"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func buildKeystoneIP(serializer *json.Serializer, p IdentityProvider) (*configv1.IdentityProvider, *secrets.Secret, *secrets.Secret, *configmaps.ConfigMap, error) {
	var (
		idP = &configv1.IdentityProvider{}

		certSecret, keySecret *secrets.Secret
		caConfigmap           *configmaps.ConfigMap
		err                   error
		keystone              legacyconfigv1.KeystonePasswordIdentityProvider
	)

	if _, _, err = serializer.Decode(p.Provider.Raw, nil, &keystone); err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "Failed to decode keystone, see error")
	}

	idP.Type = "Keystone"
	idP.Name = p.Name
	idP.MappingMethod = configv1.MappingMethodType(p.MappingMethod)
	idP.Keystone = &configv1.KeystoneIdentityProvider{}
	idP.Keystone.DomainName = keystone.DomainName
	idP.Keystone.URL = keystone.URL

	if keystone.CA != "" {
		caConfigmap = configmaps.GenConfigMap("keystone-configmap", OAuthNamespace, p.CAData)
		idP.Keystone.CA = configv1.ConfigMapNameReference{Name: caConfigmap.Metadata.Name}
	}

	if keystone.UseKeystoneIdentity {
		logrus.Warn("Keystone useKeystoneIdentity value is not supported in OCP4")
	}

	if keystone.CertFile != "" {
		certSecretName := "keystone-client-cert-secret"
		idP.Keystone.TLSClientCert.Name = certSecretName
		encoded := base64.StdEncoding.EncodeToString(p.CrtData)
		if certSecret, err = secrets.GenSecret(certSecretName, encoded, OAuthNamespace, secrets.KeystoneSecretType); err != nil {
			return nil, nil, nil, nil, errors.Wrap(err, "Failed to generate cert secret for keystone, see error")
		}

		keySecretName := "keystone-client-key-secret"
		idP.Keystone.TLSClientKey.Name = keySecretName
		encoded = base64.StdEncoding.EncodeToString(p.KeyData)
		if keySecret, err = secrets.GenSecret(keySecretName, encoded, OAuthNamespace, secrets.KeystoneSecretType); err != nil {
			return nil, nil, nil, nil, errors.Wrap(err, "Failed to generate key secret for keystone, see error")
		}
	}

	return idP, certSecret, keySecret, caConfigmap, nil
}

func validateKeystoneProvider(serializer *json.Serializer, p IdentityProvider) error {
	var keystone legacyconfigv1.KeystonePasswordIdentityProvider

	if _, _, err := serializer.Decode(p.Provider.Raw, nil, &keystone); err != nil {
		return errors.Wrap(err, "Failed to decode keystone, see error")
	}

	if p.Name == "" {
		return errors.New("Name can't be empty")
	}

	if err := validateMappingMethod(p.MappingMethod); err != nil {
		return err
	}

	if keystone.DomainName == "" {
		return errors.New("Domain name can't be empty")
	}

	if keystone.URL == "" {
		return errors.New("URL can't be empty")
	}

	if keystone.CertFile != "" && keystone.KeyFile == "" {
		return errors.New("Key file can't be empty if cert file is specified")
	}

	return nil
}

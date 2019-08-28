package oauth

import (
	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/configmaps"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func buildLdapIP(serializer *json.Serializer, p IdentityProvider) (*ProviderResources, error) {
	var (
		err                error
		idP                = &configv1.IdentityProvider{}
		providerConfigMaps []*configmaps.ConfigMap
		ldap               legacyconfigv1.LDAPPasswordIdentityProvider
	)

	if _, _, err = serializer.Decode(p.Provider.Raw, nil, &ldap); err != nil {
		return nil, errors.Wrap(err, "Failed to decode ldap. Check ssh user permissions and file content, see error")
	}

	idP.Type = "LDAP"
	idP.Name = p.Name
	idP.MappingMethod = configv1.MappingMethodType(p.MappingMethod)
	idP.LDAP = &configv1.LDAPIdentityProvider{}
	idP.LDAP.Attributes.ID = ldap.Attributes.ID
	idP.LDAP.Attributes.Email = ldap.Attributes.Email
	idP.LDAP.Attributes.Name = ldap.Attributes.Name
	idP.LDAP.Attributes.PreferredUsername = ldap.Attributes.PreferredUsername
	idP.LDAP.BindDN = ldap.BindDN

	if ldap.BindPassword.Value != "" || ldap.BindPassword.File != "" || ldap.BindPassword.Env != "" {
		bindPassword, err := io.FetchStringSource(ldap.BindPassword)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to fetch bind password for ldap")
		}

		idP.LDAP.BindPassword.Name = bindPassword
	}

	if ldap.CA != "" {
		caConfigmap := configmaps.GenConfigMap("ldap-configmap", OAuthNamespace, p.CAData)
		idP.LDAP.CA = configv1.ConfigMapNameReference{Name: caConfigmap.Metadata.Name}
		providerConfigMaps = append(providerConfigMaps, caConfigmap)
	}

	idP.LDAP.Insecure = ldap.Insecure
	idP.LDAP.URL = ldap.URL

	return &ProviderResources{
		IDP:        idP,
		ConfigMaps: providerConfigMaps,
	}, nil
}

func validateLDAPProvider(serializer *json.Serializer, p IdentityProvider) error {
	var ldap legacyconfigv1.LDAPPasswordIdentityProvider

	if _, _, err := serializer.Decode(p.Provider.Raw, nil, &ldap); err != nil {
		return errors.Wrap(err, "Failed to decode ldap. Check ssh user permissions and file content, see error")
	}

	if p.Name == "" {
		return errors.New("Name can't be empty")
	}

	if err := validateMappingMethod(p.MappingMethod); err != nil {
		return err
	}

	if len(ldap.Attributes.ID) == 0 {
		return errors.New("ID can't be empty")
	}

	if len(ldap.Attributes.Email) == 0 {
		return errors.New("Email can't be empty")
	}

	if len(ldap.Attributes.Name) == 0 {
		return errors.New("Name can't be empty")
	}

	if len(ldap.Attributes.PreferredUsername) == 0 {
		return errors.New("Preferred username can't be empty")
	}

	if ldap.URL == "" {
		return errors.New("URL can't be empty")
	}

	if ldap.BindPassword.KeyFile != "" {
		return errors.New("Usage of encrypted files as bind password value is not supported")
	}

	return nil
}

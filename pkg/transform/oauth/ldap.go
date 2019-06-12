package oauth

import (
	"errors"

	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/configmaps"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

// IdentityProviderLDAP is a LDAP specific identity provider
type IdentityProviderLDAP struct {
	identityProviderCommon `yaml:",inline"`
	LDAP                   LDAP `yaml:"ldap"`
}

// LDAP provider specific data
type LDAP struct {
	Attributes   LDAPAttributes `yaml:"attributes"`
	BindDN       string         `yaml:"bindDN,omitempty"`
	BindPassword string         `yaml:"bindPassword,omitempty"`
	CA           *CA            `yaml:"ca,omitempty"`
	Insecure     bool           `yaml:"insecure"`
	URL          string         `yaml:"url"`
}

// LDAPAttributes for an LDAP provider
type LDAPAttributes struct {
	ID                []string `yaml:"id"`
	Email             []string `yaml:"email"`
	Name              []string `yaml:"name"`
	PreferredUsername []string `yaml:"preferredUsername"`
}

func buildLdapIP(serializer *json.Serializer, p IdentityProvider) (*IdentityProviderLDAP, *configmaps.ConfigMap, error) {
	var (
		err         error
		idP         = &IdentityProviderLDAP{}
		caConfigmap *configmaps.ConfigMap
		ldap        legacyconfigv1.LDAPPasswordIdentityProvider
	)
	_, _, err = serializer.Decode(p.Provider.Raw, nil, &ldap)
	if err != nil {
		return nil, nil, err
	}

	idP.Type = "LDAP"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.LDAP.Attributes.ID = ldap.Attributes.ID
	idP.LDAP.Attributes.Email = ldap.Attributes.Email
	idP.LDAP.Attributes.Name = ldap.Attributes.Name
	idP.LDAP.Attributes.PreferredUsername = ldap.Attributes.PreferredUsername
	idP.LDAP.BindDN = ldap.BindDN

	if ldap.BindPassword.Value != "" || ldap.BindPassword.File != "" || ldap.BindPassword.Env != "" {
		bindPassword, err := io.FetchStringSource(ldap.BindPassword)
		if err != nil {
			return nil, nil, err
		}

		idP.LDAP.BindPassword = bindPassword
	}

	if ldap.CA != "" {
		caConfigmap = configmaps.GenConfigMap("ldap-configmap", OAuthNamespace, p.CAData)
		idP.LDAP.CA = &CA{Name: caConfigmap.Metadata.Name}
	}

	idP.LDAP.Insecure = ldap.Insecure
	idP.LDAP.URL = ldap.URL

	return idP, caConfigmap, nil
}

func validateLDAPProvider(serializer *json.Serializer, p IdentityProvider) error {
	var ldap legacyconfigv1.LDAPPasswordIdentityProvider

	_, _, err := serializer.Decode(p.Provider.Raw, nil, &ldap)
	if err != nil {
		return err
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

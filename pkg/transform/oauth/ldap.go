package oauth

import (
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

// IdentityProviderLDAP is a LDAP specific identity provider
type IdentityProviderLDAP struct {
	identityProviderCommon `yaml:",inline"`
	LDAP                   LDAP `yaml:"ldap"`
}

// LDAP provider specific data
type LDAP struct {
	Attributes   LDAPAttributes `yaml:"attributes"`
	BindDN       string         `yaml:"bindDN"`
	BindPassword string         `yaml:"bindPassword"`
	CA           CA             `yaml:"ca"`
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

func buildLdapIP(serializer *json.Serializer, p IdentityProvider) (IdentityProviderLDAP, error) {
	var (
		err  error
		idP  IdentityProviderLDAP
		ldap configv1.LDAPPasswordIdentityProvider
	)
	_, _, err = serializer.Decode(p.Provider.Raw, nil, &ldap)
	if err != nil {
		return idP, err
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
	idP.LDAP.CA.Name = ldap.CA
	idP.LDAP.Insecure = ldap.Insecure
	idP.LDAP.URL = ldap.URL

	return idP, nil
}

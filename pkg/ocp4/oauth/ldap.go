package oauth

import (
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

type IdentityProviderLDAP struct {
	identityProviderCommon `yaml:",inline"`
	LDAP                   struct {
		Attributes struct {
			ID                []string `yaml:"id"`
			Email             []string `yaml:"email"`
			Name              []string `yaml:"name"`
			PreferredUsername []string `yaml:"preferredUsername"`
		} `yaml:"attributes"`
		BindDN       string `yaml:"bindDN"`
		BindPassword string `yaml:"bindPassword"`
		CA           struct {
			Name string `yaml:"name"`
		} `yaml:"ca"`
		Insecure bool   `yaml:"insecure"`
		URL      string `yaml:"url"`
	} `yaml:"ldap"`
}

func buildLdapIP(serializer *json.Serializer, p configv1.IdentityProvider) IdentityProviderLDAP {
	var idP IdentityProviderLDAP
	var ldap configv1.LDAPPasswordIdentityProvider
	_, _, _ = serializer.Decode(p.Provider.Raw, nil, &ldap)

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

	return idP
}

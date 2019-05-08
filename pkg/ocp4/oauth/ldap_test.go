package oauth

import (
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigLDAP(t *testing.T) {
	file := "testdata/ldap-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)

	masterV3 := ocp3.Master{}
	masterV3.Decode(content)

	var expectedCrd OAuthCRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = "openshift-config"

	var ldapIDP identityProviderLDAP
	ldapIDP.Name = "my_ldap_provider"
	ldapIDP.Type = "LDAP"
	ldapIDP.Challenge = true
	ldapIDP.Login = true
	ldapIDP.MappingMethod = "claim"
	ldapIDP.LDAP.Attributes.ID = []string{"dn"}
	ldapIDP.LDAP.Attributes.Email = []string{"mail"}
	ldapIDP.LDAP.Attributes.Name = []string{"cn"}
	ldapIDP.LDAP.Attributes.PreferredUsername = []string{"uid"}
	ldapIDP.LDAP.BindDN = "123"
	ldapIDP.LDAP.CA.Name = "my-ldap-ca-bundle.crt"
	ldapIDP.LDAP.Insecure = false
	ldapIDP.LDAP.URL = "ldap://ldap.example.com/ou=users,dc=acme,dc=com?uid"

	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, ldapIDP)

	resCrd, _, err := Transform(masterV3.Config.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, &expectedCrd, resCrd)
}

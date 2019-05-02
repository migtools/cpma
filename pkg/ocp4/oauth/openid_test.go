package oauth

import (
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTranslateMasterConfigOpenID(t *testing.T) {
	file := "testdata/openid-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)

	masterV3 := ocp3.Master{}
	masterV3.Decode(content)

	var expectedCrd OAuthCRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = "openshift-config"

	var openidIDP identityProviderOpenID
	openidIDP.Type = "OpenID"
	openidIDP.Challenge = false
	openidIDP.Login = true
	openidIDP.MappingMethod = "claim"
	openidIDP.Name = "my_openid_connect"
	openidIDP.OpenID.ClientID = "testid"
	openidIDP.OpenID.Claims.PreferredUsername = []string{"preferred_username", "email"}
	openidIDP.OpenID.Claims.Name = []string{"nickname", "given_name", "name"}
	openidIDP.OpenID.Claims.Email = []string{"custom_email_claim", "email"}
	openidIDP.OpenID.URLs.Authorize = "https://myidp.example.com/oauth2/authorize"
	openidIDP.OpenID.URLs.Token = "https://myidp.example.com/oauth2/token"
	openidIDP.OpenID.ClientSecret.Name = "my_openid_connect-secret"

	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, openidIDP)

	resCrd, _, err := Translate(masterV3.Config.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, &expectedCrd, resCrd)
}

package oauth

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fusor/cpma/pkg/ocp3"
)

func TestTranslateMasterConfigGoogle(t *testing.T) {
	file := "testdata/google-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)

	masterV3 := ocp3.Master{}
	masterV3.Decode(content)

	var expectedCrd OAuthCRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = "openshift-config"

	var googleIDP identityProviderGoogle
	googleIDP.Type = "Google"
	googleIDP.Challenge = false
	googleIDP.Login = true
	googleIDP.MappingMethod = "claim"
	googleIDP.Name = "google123456789123456789"
	googleIDP.Google.ClientID = "82342890327-tf5lqn4eikdf4cb4edfm85jiqotvurpq.apps.googleusercontent.com"
	googleIDP.Google.ClientSecret.Name = "google123456789123456789-secret"
	googleIDP.Google.HostedDomain = "test.example.com"
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, googleIDP)

	resCrd, _, err := Translate(masterV3.Config.OAuthConfig)
	require.NoError(t, err)
	assert.Equal(t, &expectedCrd, resCrd)
}

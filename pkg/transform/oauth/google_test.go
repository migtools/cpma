package oauth_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/transform/oauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func TestTransformMasterConfigGoogle(t *testing.T) {
	file := "testdata/google-test-master-config.yaml"
	content, _ := ioutil.ReadFile(file)
	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	var masterV3 configv1.MasterConfig
	_, _, _ = serializer.Decode(content, nil, &masterV3)

	var identityProviders []oauth.IdentityProvider
	for _, identityProvider := range masterV3.OAuthConfig.IdentityProviders {
		providerJSON, _ := identityProvider.Provider.MarshalJSON()
		provider := oauth.Provider{}
		json.Unmarshal(providerJSON, &provider)

		identityProviders = append(identityProviders,
			oauth.IdentityProvider{
				provider.Kind,
				provider.APIVersion,
				identityProvider.MappingMethod,
				identityProvider.Name,
				identityProvider.Provider,
				provider.File,
				nil,
				nil,
				nil,
				identityProvider.UseAsChallenger,
				identityProvider.UseAsLogin,
			})
	}

	var expectedCrd oauth.CRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = "openshift-config"

	var googleIDP oauth.IdentityProviderGoogle
	googleIDP.Type = "Google"
	googleIDP.Challenge = false
	googleIDP.Login = true
	googleIDP.MappingMethod = "claim"
	googleIDP.Name = "google123456789123456789"
	googleIDP.Google.ClientID = "82342890327-tf5lqn4eikdf4cb4edfm85jiqotvurpq.apps.googleusercontent.com"
	googleIDP.Google.ClientSecret.Name = "google123456789123456789-secret"
	googleIDP.Google.HostedDomain = "test.example.com"
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, googleIDP)

	resCrd, _, err := oauth.Translate(identityProviders)
	require.NoError(t, err)
	assert.Equal(t, &expectedCrd, resCrd)
}

package oauth_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/transform/oauth"
	configv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestTransformMasterConfig(t *testing.T) {
	file := "testdata/bulk-test-master-config.yaml"

	content, err := ioutil.ReadFile(file)
	require.NoError(t, err)

	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	var masterV3 configv1.MasterConfig

	_, _, err = serializer.Decode(content, nil, &masterV3)
	require.NoError(t, err)

	var identityProviders []oauth.IdentityProvider
	for _, identityProvider := range masterV3.OAuthConfig.IdentityProviders {
		providerJSON, err := identityProvider.Provider.MarshalJSON()
		require.NoError(t, err)

		provider := oauth.Provider{}

		err = json.Unmarshal(providerJSON, &provider)
		require.NoError(t, err)

		identityProviders = append(identityProviders,
			oauth.IdentityProvider{
				Kind:            provider.Kind,
				APIVersion:      provider.APIVersion,
				MappingMethod:   identityProvider.MappingMethod,
				Name:            identityProvider.Name,
				Provider:        identityProvider.Provider,
				HTFileName:      provider.File,
				UseAsChallenger: identityProvider.UseAsChallenger,
				UseAsLogin:      identityProvider.UseAsLogin,
			})
	}

	testCases := []struct {
		name string
	}{
		{
			name: "transform master config",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oauthResources, err := oauth.Translate(identityProviders)
			require.NoError(t, err)
			assert.Equal(t, len(oauthResources.OAuthCRD.Spec.IdentityProviders), 9)
			assert.Equal(t, oauthResources.OAuthCRD.Spec.IdentityProviders[0].(*oauth.IdentityProviderBasicAuth).Type, "BasicAuth")
			assert.Equal(t, oauthResources.OAuthCRD.Spec.IdentityProviders[1].(*oauth.IdentityProviderGitHub).Type, "GitHub")
			assert.Equal(t, oauthResources.OAuthCRD.Spec.IdentityProviders[2].(*oauth.IdentityProviderGitLab).Type, "GitLab")
			assert.Equal(t, oauthResources.OAuthCRD.Spec.IdentityProviders[3].(*oauth.IdentityProviderGoogle).Type, "Google")
			assert.Equal(t, oauthResources.OAuthCRD.Spec.IdentityProviders[4].(*oauth.IdentityProviderHTPasswd).Type, "HTPasswd")
			assert.Equal(t, oauthResources.OAuthCRD.Spec.IdentityProviders[5].(*oauth.IdentityProviderKeystone).Type, "Keystone")
			assert.Equal(t, oauthResources.OAuthCRD.Spec.IdentityProviders[6].(*oauth.IdentityProviderLDAP).Type, "LDAP")
			assert.Equal(t, oauthResources.OAuthCRD.Spec.IdentityProviders[7].(*oauth.IdentityProviderRequestHeader).Type, "RequestHeader")
			assert.Equal(t, oauthResources.OAuthCRD.Spec.IdentityProviders[8].(*oauth.IdentityProviderOpenID).Type, "OpenID")
		})
	}

}

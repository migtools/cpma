package oauth_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/transform/oauth"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func loadExpectedOAuth(file string) (*configv1.OAuth, error) {
	expectedContent, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	expectedOAuth := new(configv1.OAuth)
	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	_, _, err = serializer.Decode(expectedContent, nil, expectedOAuth)
	if err != nil {
		return nil, err
	}

	return expectedOAuth, nil
}

func TestTransformMasterConfig(t *testing.T) {
	file := "testdata/master_config-bulk.yaml"

	content, err := ioutil.ReadFile(file)
	require.NoError(t, err)

	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	var masterV3 legacyconfigv1.MasterConfig

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
				Kind:          provider.Kind,
				APIVersion:    provider.APIVersion,
				MappingMethod: identityProvider.MappingMethod,
				Name:          identityProvider.Name,
				Provider:      identityProvider.Provider,
				HTFileName:    provider.File,
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
			oauthResources, err := oauth.Translate(identityProviders, oauth.TokenConfig{AccessTokenMaxAgeSeconds: 42000, AuthorizeTokenMaxAgeSeconds: 42000})
			require.NoError(t, err)
			assert.Equal(t, 9, len(oauthResources.OAuthCRD.Spec.IdentityProviders))
			assert.Equal(t, configv1.IdentityProviderType("BasicAuth"), oauthResources.OAuthCRD.Spec.IdentityProviders[0].Type)
			assert.Equal(t, configv1.IdentityProviderType("GitHub"), oauthResources.OAuthCRD.Spec.IdentityProviders[1].Type)
			assert.Equal(t, configv1.IdentityProviderType("GitLab"), oauthResources.OAuthCRD.Spec.IdentityProviders[2].Type)
			assert.Equal(t, configv1.IdentityProviderType("Google"), oauthResources.OAuthCRD.Spec.IdentityProviders[3].Type)
			assert.Equal(t, configv1.IdentityProviderType("HTPasswd"), oauthResources.OAuthCRD.Spec.IdentityProviders[4].Type)
			assert.Equal(t, configv1.IdentityProviderType("Keystone"), oauthResources.OAuthCRD.Spec.IdentityProviders[5].Type)
			assert.Equal(t, configv1.IdentityProviderType("LDAP"), oauthResources.OAuthCRD.Spec.IdentityProviders[6].Type)
			assert.Equal(t, configv1.IdentityProviderType("RequestHeader"), oauthResources.OAuthCRD.Spec.IdentityProviders[7].Type)
			assert.Equal(t, configv1.IdentityProviderType("OpenID"), oauthResources.OAuthCRD.Spec.IdentityProviders[8].Type)

			assert.Equal(t, int32(42000), oauthResources.OAuthCRD.Spec.TokenConfig.AccessTokenMaxAgeSeconds)
		})
	}

}

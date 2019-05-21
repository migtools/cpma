package oauth_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fusor/cpma/pkg/transform/oauth"
	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func TestTransformMasterConfigGithub(t *testing.T) {
	file := "testdata/github-test-master-config.yaml"

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

	var expectedCrd oauth.CRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = oauth.OAuthNamespace

	var githubIDP oauth.IdentityProviderGitHub
	githubIDP.Type = "GitHub"
	githubIDP.Challenge = false
	githubIDP.Login = true
	githubIDP.MappingMethod = "claim"
	githubIDP.Name = "github123456789"
	githubIDP.GitHub.HostName = "test.example.com"
	githubIDP.GitHub.CA.Name = "github-configmap"
	githubIDP.GitHub.ClientID = "2d85ea3f45d6777bffd7"
	githubIDP.GitHub.Organizations = []string{"myorganization1", "myorganization2"}
	githubIDP.GitHub.Teams = []string{"myorganization1/team-a", "myorganization2/team-b"}
	githubIDP.GitHub.ClientSecret.Name = "github123456789-secret"
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, githubIDP)

	testCases := []struct {
		name        string
		expectedCrd *oauth.CRD
	}{
		{
			name:        "build github provider",
			expectedCrd: &expectedCrd,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resCrd, _, _, err := oauth.Translate(identityProviders)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedCrd, resCrd)
		})
	}
}

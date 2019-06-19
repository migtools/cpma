package transform

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/transform/configmaps"
	"github.com/fusor/cpma/pkg/transform/oauth"
	"github.com/fusor/cpma/pkg/transform/secrets"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestOauthGenYAML(t *testing.T) {
	testCases := []struct {
		name                    string
		inputConfigfile         string
		expectedYaml            string
		expectedSecretsLength   int
		expectedConfigMapsength int
	}{
		{
			name:                    "generate yaml for oauth providers",
			inputConfigfile:         "testdata/master_config-bulk.yaml",
			expectedYaml:            "testdata/expected-master_config-oauth-bulk.yaml",
			expectedSecretsLength:   9,
			expectedConfigMapsength: 6,
		},
		{
			name:                    "generate yaml for oauth providers and omit empty values",
			inputConfigfile:         "testdata/master_config-omit-empty-values.yaml",
			expectedYaml:            "testdata/expected-master_config-oauth-omit-empty-values.yaml",
			expectedSecretsLength:   5,
			expectedConfigMapsength: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectedYaml, err := ioutil.ReadFile(tc.expectedYaml)
			require.NoError(t, err)

			content, err := ioutil.ReadFile(tc.inputConfigfile)
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

			oauthResources, err := oauth.Translate(identityProviders, oauth.TokenConfig{
				AccessTokenMaxAgeSeconds:    int32(86400),
				AuthorizeTokenMaxAgeSeconds: int32(500),
			})
			require.NoError(t, err)

			CRD, err := GenYAML(oauthResources.OAuthCRD)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedSecretsLength, len(oauthResources.Secrets))
			assert.Equal(t, tc.expectedConfigMapsength, len(oauthResources.ConfigMaps))
			assert.Equal(t, expectedYaml, CRD)
		})
	}
}

func TestAllOtherCRGenYaml(t *testing.T) {
	expectedConfigMapYaml, err := ioutil.ReadFile("testdata/expected-CR-configmap.yaml")
	require.NoError(t, err)

	expectedSecretYaml, err := ioutil.ReadFile("testdata/expected-CR-secret.yaml")
	require.NoError(t, err)

	testCases := []struct {
		name         string
		inputCR      interface{}
		expectedYaml []byte
	}{
		{
			name: "generate yaml from configmap",
			inputCR: configmaps.ConfigMap{
				APIVersion: configmaps.APIVersion,
				Data: configmaps.Data{
					CAData: "testval: 123",
				},
				Kind: configmaps.Kind,
				Metadata: configmaps.MetaData{
					Name:      "testname",
					Namespace: "openshift-config",
				},
			},
			expectedYaml: expectedConfigMapYaml,
		},
		{
			name: "generate yaml from secret",
			inputCR: secrets.Secret{
				APIVersion: secrets.APIVersion,
				Data:       secrets.LiteralSecret{ClientSecret: "some-value"},
				Kind:       "Secret",
				Type:       "Opaque",
				Metadata: secrets.MetaData{
					Name:      "literal-secret",
					Namespace: "openshift-config",
				},
			},
			expectedYaml: expectedSecretYaml,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			manifest, err := GenYAML(tc.inputCR)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedYaml, manifest)
		})
	}
}

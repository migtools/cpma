package transform

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/konveyor/cpma/pkg/decode"
	"github.com/konveyor/cpma/pkg/env"
	"github.com/konveyor/cpma/pkg/transform/oauth"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestOauthGenYAML(t *testing.T) {
	env.Config().Set("Manifests", true)
	env.Config().Set("Reporting", true)

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

			masterV3, err := decode.MasterConfig(content)
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
			}, legacyconfigv1.OAuthTemplates{})
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
			inputCR: corev1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "ConfigMap",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testname",
					Namespace: "openshift-config",
				},
				Data: map[string]string{
					"ca.crt": "testval: 123",
				},
			},
			expectedYaml: expectedConfigMapYaml,
		},
		{
			name: "generate yaml from secret",
			inputCR: corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				Data: map[string][]byte{
					"clientSecret": []byte("some-value"),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "literal-secret",
					Namespace: "openshift-config",
				},
				Type: "Opaque",
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

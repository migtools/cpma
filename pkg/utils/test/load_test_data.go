package test

import (
	"encoding/json"
	"io/ioutil"

	"github.com/fusor/cpma/pkg/transform/oauth"

	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

// LoadIPTestData load identity providers from file
func LoadIPTestData(file string) ([]oauth.IdentityProvider, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	var masterV3 configv1.MasterConfig

	_, _, err = serializer.Decode(content, nil, &masterV3)
	if err != nil {
		return nil, err
	}

	var identityProviders []oauth.IdentityProvider
	for _, identityProvider := range masterV3.OAuthConfig.IdentityProviders {
		providerJSON, err := identityProvider.Provider.MarshalJSON()
		if err != nil {
			return nil, err
		}

		provider := oauth.Provider{}

		err = json.Unmarshal(providerJSON, &provider)
		if err != nil {
			return nil, err
		}

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

	return identityProviders, nil
}

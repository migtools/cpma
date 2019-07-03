package oauth

import (
	"github.com/fusor/cpma/pkg/transform/configmaps"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func buildRequestHeaderIP(serializer *json.Serializer, p IdentityProvider) (*ProviderResources, error) {
	var (
		err                error
		idP                = &configv1.IdentityProvider{}
		providerConfigMaps []*configmaps.ConfigMap
		requestHeader      legacyconfigv1.RequestHeaderIdentityProvider
	)

	if _, _, err = serializer.Decode(p.Provider.Raw, nil, &requestHeader); err != nil {
		return nil, errors.Wrap(err, "Failed to decode request header, see error")
	}

	idP.Type = "RequestHeader"
	idP.Name = p.Name
	idP.MappingMethod = configv1.MappingMethodType(p.MappingMethod)
	idP.RequestHeader = &configv1.RequestHeaderIdentityProvider{}
	idP.RequestHeader.ChallengeURL = requestHeader.ChallengeURL
	idP.RequestHeader.LoginURL = requestHeader.LoginURL

	if requestHeader.ClientCA != "" {
		caConfigmap := configmaps.GenConfigMap("requestheader-configmap", OAuthNamespace, p.CAData)
		idP.RequestHeader.ClientCA = configv1.ConfigMapNameReference{Name: caConfigmap.Metadata.Name}
		providerConfigMaps = append(providerConfigMaps, caConfigmap)
	}

	idP.RequestHeader.ClientCommonNames = requestHeader.ClientCommonNames
	idP.RequestHeader.Headers = requestHeader.Headers
	idP.RequestHeader.EmailHeaders = requestHeader.EmailHeaders
	idP.RequestHeader.NameHeaders = requestHeader.NameHeaders
	idP.RequestHeader.PreferredUsernameHeaders = requestHeader.PreferredUsernameHeaders

	return &ProviderResources{
		IDP:        idP,
		ConfigMaps: providerConfigMaps,
	}, nil
}

func validateRequestHeaderProvider(serializer *json.Serializer, p IdentityProvider) error {
	var requestHeader legacyconfigv1.RequestHeaderIdentityProvider

	if _, _, err := serializer.Decode(p.Provider.Raw, nil, &requestHeader); err != nil {
		return errors.Wrap(err, "Failed to decode request header, see error")
	}

	if p.Name == "" {
		return errors.New("Name can't be empty")
	}

	if err := validateMappingMethod(p.MappingMethod); err != nil {
		return err
	}

	if len(requestHeader.Headers) == 0 {
		return errors.New("Headers can't be empty")
	}

	return nil
}

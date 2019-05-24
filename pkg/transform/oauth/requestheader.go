package oauth

import (
	"errors"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	"github.com/fusor/cpma/pkg/transform/configmaps"
	configv1 "github.com/openshift/api/legacyconfig/v1"
)

// IdentityProviderRequestHeader is a request header specific identity provider
type IdentityProviderRequestHeader struct {
	identityProviderCommon `yaml:",inline"`
	RequestHeader          RequestHeader `yaml:"requestHeader"`
}

// RequestHeader provider specific data
type RequestHeader struct {
	ChallengeURL             string   `yaml:"challengeURL,omitempty"`
	LoginURL                 string   `yaml:"loginURL,omitempty"`
	CA                       *CA      `yaml:"ca,omitempty"`
	ClientCommonNames        []string `yaml:"сlientCommonNames,omitempty"`
	Headers                  []string `yaml:"headers"`
	EmailHeaders             []string `yaml:"emailHeaders,omitempty"`
	NameHeaders              []string `yaml:"nameHeaders,omitempty"`
	PreferredUsernameHeaders []string `yaml:"preferredUsernameHeaders,omitempty"`
}

func buildRequestHeaderIP(serializer *json.Serializer, p IdentityProvider) (*IdentityProviderRequestHeader, *configmaps.ConfigMap, error) {
	var (
		err           error
		idP           = &IdentityProviderRequestHeader{}
		caConfigmap   *configmaps.ConfigMap
		requestHeader configv1.RequestHeaderIdentityProvider
	)
	_, _, err = serializer.Decode(p.Provider.Raw, nil, &requestHeader)
	if err != nil {
		return nil, nil, err
	}

	idP.Type = "RequestHeader"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.RequestHeader.ChallengeURL = requestHeader.ChallengeURL
	idP.RequestHeader.LoginURL = requestHeader.LoginURL

	if requestHeader.ClientCA != "" {
		caConfigmap = configmaps.GenConfigMap("requestheader-configmap", OAuthNamespace, p.CAData)
		idP.RequestHeader.CA = &CA{Name: caConfigmap.Metadata.Name}
	}

	idP.RequestHeader.ClientCommonNames = requestHeader.ClientCommonNames
	idP.RequestHeader.Headers = requestHeader.Headers
	idP.RequestHeader.EmailHeaders = requestHeader.EmailHeaders
	idP.RequestHeader.NameHeaders = requestHeader.NameHeaders
	idP.RequestHeader.PreferredUsernameHeaders = requestHeader.PreferredUsernameHeaders

	return idP, caConfigmap, nil
}

func validateRequestHeaderProvider(serializer *json.Serializer, p IdentityProvider) error {
	var requestHeader configv1.RequestHeaderIdentityProvider

	_, _, err := serializer.Decode(p.Provider.Raw, nil, &requestHeader)
	if err != nil {
		return err
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

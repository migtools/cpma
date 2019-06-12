package oauth

import (
	"errors"

	"github.com/fusor/cpma/pkg/transform/configmaps"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

// IdentityProviderRequestHeader is a request header specific identity provider
type IdentityProviderRequestHeader struct {
	identityProviderCommon `json:",inline"`
	RequestHeader          RequestHeader `json:"requestHeader"`
}

// RequestHeader provider specific data
type RequestHeader struct {
	ChallengeURL             string   `json:"challengeURL,omitempty"`
	LoginURL                 string   `json:"loginURL,omitempty"`
	CA                       *CA      `json:"ca,omitempty"`
	ClientCommonNames        []string `json:"сlientCommonNames,omitempty"`
	Headers                  []string `json:"headers"`
	EmailHeaders             []string `json:"emailHeaders,omitempty"`
	NameHeaders              []string `json:"nameHeaders,omitempty"`
	PreferredUsernameHeaders []string `json:"preferredUsernameHeaders,omitempty"`
}

func buildRequestHeaderIP(serializer *json.Serializer, p IdentityProvider) (*IdentityProviderRequestHeader, *configmaps.ConfigMap, error) {
	var (
		err           error
		idP           = &IdentityProviderRequestHeader{}
		caConfigmap   *configmaps.ConfigMap
		requestHeader legacyconfigv1.RequestHeaderIdentityProvider
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
	var requestHeader legacyconfigv1.RequestHeaderIdentityProvider

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

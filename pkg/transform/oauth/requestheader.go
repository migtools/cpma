package oauth

import (
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

// IdentityProviderRequestHeader is a request header specific identity provider
type IdentityProviderRequestHeader struct {
	identityProviderCommon `yaml:",inline"`
	RequestHeader          struct {
		ChallengeURL string `yaml:"challengeURL"`
		LoginURL     string `yaml:"loginURL"`
		CA           struct {
			Name string `yaml:"name"`
		} `yaml:"ca"`
		ClientCommonNames        []string `yaml:"сlientCommonNames"`
		Headers                  []string `yaml:"headers"`
		EmailHeaders             []string `yaml:"emailHeaders"`
		NameHeaders              []string `yaml:"nameHeaders"`
		PreferredUsernameHeaders []string `yaml:"preferredUsernameHeaders"`
	} `yaml:"requestHeader"`
}

func buildRequestHeaderIP(serializer *json.Serializer, p IdentityProvider) IdentityProviderRequestHeader {
	var idP IdentityProviderRequestHeader
	var requestHeader configv1.RequestHeaderIdentityProvider
	_, _, _ = serializer.Decode(p.Provider.Raw, nil, &requestHeader)

	idP.Type = "RequestHeader"
	idP.Name = p.Name
	idP.Challenge = p.UseAsChallenger
	idP.Login = p.UseAsLogin
	idP.MappingMethod = p.MappingMethod
	idP.RequestHeader.ChallengeURL = requestHeader.ChallengeURL
	idP.RequestHeader.LoginURL = requestHeader.LoginURL
	idP.RequestHeader.CA.Name = requestHeader.ClientCA
	idP.RequestHeader.ClientCommonNames = requestHeader.ClientCommonNames
	idP.RequestHeader.Headers = requestHeader.Headers
	idP.RequestHeader.EmailHeaders = requestHeader.EmailHeaders
	idP.RequestHeader.NameHeaders = requestHeader.NameHeaders
	idP.RequestHeader.PreferredUsernameHeaders = requestHeader.PreferredUsernameHeaders

	return idP
}

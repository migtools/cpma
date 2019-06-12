package oauth

import (
	"errors"

	"github.com/fusor/cpma/pkg/transform/configmaps"
	"github.com/fusor/cpma/pkg/transform/secrets"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	oauthv1 "github.com/openshift/api/oauth/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	oauthv1.Install(scheme.Scheme)
	legacyconfigv1.InstallLegacy(scheme.Scheme)
}

// reference:
//   [v3] OCPv3:
//   - [1] https://docs.openshift.com/container-platform/3.11/install_config/configuring_authentication.html#identity_providers_master_config
//   [v4] OCPv4:
//   - [2] htpasswd: https://docs.openshift.com/container-platform/4.0/authentication/understanding-identity-provider.html
//   - [3] github: https://docs.openshift.com/container-platform/4.0/authentication/identity_providers/configuring-github-identity-provider.html

// CRD Shared CRD part, present in all types of OAuth CRDs
type CRD struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   MetaData `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

// Spec is a CRD Spec
type Spec struct {
	IdentityProviders []interface{}         `yaml:"identityProviders,omitempty"`
	TokenConfig       TranslatedTokenConfig `yaml:"tokenConfig,omitempty"`
}

// TranslatedTokenConfig holds lifetime of access tokens
type TranslatedTokenConfig struct {
	AccessTokenMaxAgeSeconds int32 `yaml:"accessTokenMaxAgeSeconds"`
}

type identityProviderCommon struct {
	Name          string `yaml:"name"`
	Challenge     bool   `yaml:"challenge"`
	Login         bool   `yaml:"login"`
	MappingMethod string `yaml:"mappingMethod"`
	Type          string `yaml:"type"`
}

// MetaData contains CRD Metadata
type MetaData struct {
	Name      string `yaml:"name"`
	NameSpace string `yaml:"namespace"`
}

// Provider contains an identity providers type specific provider data
type Provider struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	File       string `json:"file"`
	CA         string `json:"ca"`
	CertFile   string `json:"certFile"`
	KeyFile    string `json:"keyFile"`
}

// IdentityProvider stores an identity provider
type IdentityProvider struct {
	Kind            string
	APIVersion      string
	MappingMethod   string
	Name            string
	Provider        runtime.RawExtension
	HTFileName      string
	HTFileData      []byte
	CAData          []byte
	CrtData         []byte
	KeyData         []byte
	UseAsChallenger bool
	UseAsLogin      bool
}

// Resources stores all oAuth config parts
type Resources struct {
	OAuthCRD   *CRD
	Secrets    []*secrets.Secret
	ConfigMaps []*configmaps.ConfigMap
}

// TokenConfig store internal OAuth tokens duration
type TokenConfig struct {
	AuthorizeTokenMaxAgeSeconds int32
	AccessTokenMaxAgeSeconds    int32
}

const (
	// APIVersion is the apiVersion string
	APIVersion = "config.openshift.io/v1"
	// OAuthNamespace is namespace for oauth manifests
	OAuthNamespace = "openshift-config"
)

// Translate converts OCPv3 OAuth to OCPv4 OAuth Custom Resources
func Translate(identityProviders []IdentityProvider, tokenConfig TokenConfig) (*Resources, error) {
	var err error
	var idP interface{}
	var secretsSlice []*secrets.Secret
	var сonfigMapSlice []*configmaps.ConfigMap

	// Translate configuration of diffent oAuth providers to CRD, secrets and config maps
	var oauthCrd CRD
	oauthCrd.APIVersion = APIVersion
	oauthCrd.Kind = "OAuth"
	oauthCrd.Metadata.Name = "cluster"
	oauthCrd.Metadata.NameSpace = OAuthNamespace
	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	for _, p := range identityProviders {
		var secret, certSecret, keySecret *secrets.Secret
		var caConfigMap *configmaps.ConfigMap

		p.Provider.Object, _, err = serializer.Decode(p.Provider.Raw, nil, nil)
		if err != nil {
			return nil, err
		}

		kind := p.Kind

		switch kind {
		case "GitHubIdentityProvider":
			idP, secret, caConfigMap, err = buildGitHubIP(serializer, p)
		case "GitLabIdentityProvider":
			idP, secret, caConfigMap, err = buildGitLabIP(serializer, p)
		case "GoogleIdentityProvider":
			idP, secret, err = buildGoogleIP(serializer, p)
		case "HTPasswdPasswordIdentityProvider":
			idP, secret, err = buildHTPasswdIP(serializer, p)
		case "OpenIDIdentityProvider":
			idP, secret, err = buildOpenIDIP(serializer, p)
		case "RequestHeaderIdentityProvider":
			idP, caConfigMap, err = buildRequestHeaderIP(serializer, p)
		case "LDAPPasswordIdentityProvider":
			idP, caConfigMap, err = buildLdapIP(serializer, p)
		case "KeystonePasswordIdentityProvider":
			idP, certSecret, keySecret, caConfigMap, err = buildKeystoneIP(serializer, p)
		case "BasicAuthPasswordIdentityProvider":
			idP, certSecret, keySecret, caConfigMap, err = buildBasicAuthIP(serializer, p)
		default:
			logrus.Infof("Can't handle %s OAuth kind", kind)
			continue
		}

		// Skip OAuth provider if error was returned
		if err != nil {
			logrus.Error("Can't handle ", kind, " skipping.. error:", err)
			continue
		}

		// Check if secret is not empty
		if secret != nil {
			secretsSlice = append(secretsSlice, secret)
		}

		// Check if certSecret is not empty
		if certSecret != nil {
			secretsSlice = append(secretsSlice, certSecret)
			secretsSlice = append(secretsSlice, keySecret)
		}

		// Check if config map is not empty
		if caConfigMap != nil {
			сonfigMapSlice = append(сonfigMapSlice, caConfigMap)
		}

		oauthCrd.Spec.IdentityProviders = append(oauthCrd.Spec.IdentityProviders, idP)
	}

	// Translate lifetime of access tokens
	oauthCrd.Spec.TokenConfig.AccessTokenMaxAgeSeconds = tokenConfig.AccessTokenMaxAgeSeconds

	return &Resources{
		OAuthCRD:   &oauthCrd,
		Secrets:    secretsSlice,
		ConfigMaps: сonfigMapSlice,
	}, nil
}

// Validate validate oauth providers
func Validate(identityProviders []IdentityProvider) error {
	var err error
	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	for _, identityProvider := range identityProviders {
		switch identityProvider.Kind {
		case "BasicAuthPasswordIdentityProvider":
			err = validateBasicAuthProvider(serializer, identityProvider)
		case "GitHubIdentityProvider":
			err = validateGithubProvider(serializer, identityProvider)
		case "GitLabIdentityProvider":
			err = validateGitLabProvider(serializer, identityProvider)
		case "GoogleIdentityProvider":
			err = validateGoogleProvider(serializer, identityProvider)
		case "HTPasswdPasswordIdentityProvider":
			err = validateHTPasswdProvider(serializer, identityProvider)
		case "KeystonePasswordIdentityProvider":
			err = validateKeystoneProvider(serializer, identityProvider)
		case "LDAPPasswordIdentityProvider":
			err = validateLDAPProvider(serializer, identityProvider)
		case "OpenIDIdentityProvider":
			err = validateOpenIDProvider(serializer, identityProvider)
		case "RequestHeaderIdentityProvider":
			err = validateRequestHeaderProvider(serializer, identityProvider)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func validateMappingMethod(method string) error {
	switch method {
	case "claim", "lookup", "generate", "add":
		return nil
	}

	return errors.New("Not valid mapping method")
}

func validateClientData(clientID string, clientSecret legacyconfigv1.StringSource) error {
	if clientID == "" {
		return errors.New("Client ID can't be empty")
	}

	if clientSecret.Value == "" && clientSecret.Env == "" && clientSecret.File == "" {
		return errors.New("Client Secret can't be empty")
	}

	return nil
}

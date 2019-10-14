package oauth

import (
	"errors"

	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	oauthv1 "github.com/openshift/api/oauth/v1"
	"github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
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
	Kind          string
	APIVersion    string
	MappingMethod string
	Name          string
	Provider      runtime.RawExtension
	HTFileName    string
	HTFileData    []byte
	CAData        []byte
	CrtData       []byte
	KeyData       []byte
}

// ResultResources stores all oAuth config parts
type ResultResources struct {
	OAuthCRD   *configv1.OAuth
	Secrets    []*corev1.Secret
	ConfigMaps []*corev1.ConfigMap
}

// ProviderResources stores all resources related to one provider
type ProviderResources struct {
	IDP        *configv1.IdentityProvider
	Secrets    []*corev1.Secret
	ConfigMaps []*corev1.ConfigMap
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
func Translate(identityProviders []IdentityProvider, tokenConfig TokenConfig, templates legacyconfigv1.OAuthTemplates) (*ResultResources, error) {
	var err error
	var secretsSlice []*corev1.Secret
	var сonfigMapSlice []*corev1.ConfigMap
	var providerResources *ProviderResources

	// Translate configuration of diffent oAuth providers to CRD, secrets and config maps
	var oauthCrd configv1.OAuth
	oauthCrd.APIVersion = APIVersion
	oauthCrd.Kind = "OAuth"
	oauthCrd.Name = "cluster"
	oauthCrd.Namespace = OAuthNamespace
	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	for _, p := range identityProviders {
		p.Provider.Object, _, err = serializer.Decode(p.Provider.Raw, nil, nil)
		if err != nil {
			return nil, err
		}

		kind := p.Kind

		switch kind {
		case "GitHubIdentityProvider":
			providerResources, err = buildGitHubIP(serializer, p)
		case "GitLabIdentityProvider":
			providerResources, err = buildGitLabIP(serializer, p)
		case "GoogleIdentityProvider":
			providerResources, err = buildGoogleIP(serializer, p)
		case "HTPasswdPasswordIdentityProvider":
			providerResources, err = buildHTPasswdIP(serializer, p)
		case "OpenIDIdentityProvider":
			providerResources, err = buildOpenIDIP(serializer, p)
		case "RequestHeaderIdentityProvider":
			providerResources, err = buildRequestHeaderIP(serializer, p)
		case "LDAPPasswordIdentityProvider":
			providerResources, err = buildLdapIP(serializer, p)
		case "KeystonePasswordIdentityProvider":
			providerResources, err = buildKeystoneIP(serializer, p)
		case "BasicAuthPasswordIdentityProvider":
			providerResources, err = buildBasicAuthIP(serializer, p)
		default:
			logrus.Warnf("Can't handle %s OAuth kind", kind)
			continue
		}

		// Skip OAuth provider if error was returned
		if err != nil {
			logrus.Error("Can't handle ", kind, " skipping.. error:", err)
			continue
		}

		// Check if provider has secrets
		if len(providerResources.Secrets) != 0 {
			secretsSlice = append(secretsSlice, providerResources.Secrets...)
		}

		// Check if provider has configmaps
		if len(providerResources.ConfigMaps) != 0 {
			сonfigMapSlice = append(сonfigMapSlice, providerResources.ConfigMaps...)
		}

		oauthCrd.Spec.IdentityProviders = append(oauthCrd.Spec.IdentityProviders, *providerResources.IDP)
	}

	// Translate lifetime of access tokens
	oauthCrd.Spec.TokenConfig.AccessTokenMaxAgeSeconds = tokenConfig.AccessTokenMaxAgeSeconds

	// Translate templates that allow to customize the login page
	translatedTemplates, templateSecrets, err := translateTemplates(templates)
	if err != nil {
		return nil, err
	}

	oauthCrd.Spec.Templates = *translatedTemplates
	secretsSlice = append(secretsSlice, templateSecrets...)

	return &ResultResources{
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

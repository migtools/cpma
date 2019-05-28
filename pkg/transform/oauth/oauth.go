package oauth

import (
	"errors"
	"strings"

	"github.com/fusor/cpma/pkg/config"
	"github.com/fusor/cpma/pkg/transform/configmaps"
	"github.com/fusor/cpma/pkg/transform/secrets"
	configv1 "github.com/openshift/api/legacyconfig/v1"
	oauthv1 "github.com/openshift/api/oauth/v1"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	oauthv1.Install(scheme.Scheme)
	configv1.InstallLegacy(scheme.Scheme)
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
	IdentityProviders []interface{} `yaml:"identityProviders"`
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

// IdentityProvider stroes an identity provider
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

const (
	// APIVersion is the apiVersion string
	APIVersion = "config.openshift.io/v1"
	// OAuthNamespace is namespace for oauth manifests
	OAuthNamespace = "openshift-config"
)

// Translate converts OCPv3 OAuth to OCPv4 OAuth Custom Resources
func Translate(identityProviders []IdentityProvider, config *config.Config) (*CRD, []*secrets.Secret, []*configmaps.ConfigMap, error) {
	var err error
	var idP interface{}
	var secretsSlice []*secrets.Secret
	var сonfigMapSlice []*configmaps.ConfigMap

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
			return nil, nil, nil, err
		}

		kind := p.Kind

		switch kind {
		case "GitHubIdentityProvider":
			idP, secret, caConfigMap, err = buildGitHubIP(serializer, p, config)
		case "GitLabIdentityProvider":
			idP, secret, caConfigMap, err = buildGitLabIP(serializer, p, config)
		case "GoogleIdentityProvider":
			idP, secret, err = buildGoogleIP(serializer, p, config)
		case "HTPasswdPasswordIdentityProvider":
			idP, secret, err = buildHTPasswdIP(serializer, p)
		case "OpenIDIdentityProvider":
			idP, secret, err = buildOpenIDIP(serializer, p, config)
		case "RequestHeaderIdentityProvider":
			idP, caConfigMap, err = buildRequestHeaderIP(serializer, p)
		case "LDAPPasswordIdentityProvider":
			idP, caConfigMap, err = buildLdapIP(serializer, p, config)
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

	return &oauthCrd, secretsSlice, сonfigMapSlice, nil
}

// GenYAML returns a YAML of the CRD
func (oauth *CRD) GenYAML() ([]byte, error) {
	yamlBytes, err := yaml.Marshal(&oauth)
	if err != nil {
		logrus.Debugf("Error in OAuth CRD, OAuth CRD - %+v", yamlBytes)
		return nil, err
	}

	return yamlBytes, nil
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

func validateClientData(clientID string, clientSecret configv1.StringSource) error {
	if clientID == "" {
		return errors.New("Client ID can't be empty")
	}

	if clientSecret.Value == "" && clientSecret.Env == "" && clientSecret.File == "" {
		return errors.New("Client Secret can't be empty")
	}

	return nil
}

func fetchStringSource(stringSource configv1.StringSource, config *config.Config) (string, error) {
	if stringSource.Value != "" {
		return stringSource.Value, nil
	}

	if stringSource.File != "" {
		fileContent, err := config.Fetch(stringSource.File)
		if err != nil {
			return "", nil
		}

		fileString := strings.TrimSuffix(string(fileContent), "\n")
		return fileString, nil
	}

	if stringSource.Env != "" {
		env, err := config.FetchEnv(stringSource.Env)
		if err != nil {
			return "", nil
		}

		return env, nil
	}

	return "", nil
}

package oauth

import (
	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/secrets"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
	oauthv1 "github.com/openshift/api/oauth/v1"
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
	CrtData         []byte
	KeyData         []byte
	UseAsChallenger bool
	UseAsLogin      bool
}

var (
	// APIVersion is the apiVersion string
	APIVersion = "config.openshift.io/v1"
	// GetFile allows to mock file retrieval
	GetFile = io.GetFile
)

const (
	// httpdAuthSecret is the htpasswd auth secret string
	httpdAuthSecret = "htpasswd_auth-secret"
	// httpProviderKind is the htpasswd idp type string
	httpProviderKind = "HTPasswdPasswordIdentityProvider"
)

// Translate converts OCPv3 OAuth to OCPv4 OAuth Custom Resources
func Translate(identityProviders []IdentityProvider) (*CRD, []secrets.Secret, error) {
	var err error
	var idP interface{}
	var secretsSlice []secrets.Secret
	var secret, certSecret, keySecret secrets.Secret

	var oauthCrd CRD
	oauthCrd.APIVersion = APIVersion
	oauthCrd.Kind = "OAuth"
	oauthCrd.Metadata.Name = "cluster"
	oauthCrd.Metadata.NameSpace = "openshift-config"
	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	for _, p := range identityProviders {
		p.Provider.Object, _, err = serializer.Decode(p.Provider.Raw, nil, nil)
		if err != nil {
			return nil, nil, err
		}

		switch kind := p.Kind; kind {
		case "GitHubIdentityProvider":
			idP, secret = buildGitHubIP(serializer, p)
		case "GitLabIdentityProvider":
			idP, secret = buildGitLabIP(serializer, p)
		case "GoogleIdentityProvider":
			idP, secret = buildGoogleIP(serializer, p)
		case "HTPasswdPasswordIdentityProvider":
			idP, secret = buildHTPasswdIP(serializer, p)
		case "OpenIDIdentityProvider":
			idP, secret = buildOpenIDIP(serializer, p)
		case "RequestHeaderIdentityProvider":
			idP = buildRequestHeaderIP(serializer, p)
		case "LDAPPasswordIdentityProvider":
			idP = buildLdapIP(serializer, p)
		case "KeystonePasswordIdentityProvider":
			idP, certSecret, keySecret = buildKeystoneIP(serializer, p)
			if certSecret != (secrets.Secret{}) {
				secretsSlice = append(secretsSlice, certSecret)
				secretsSlice = append(secretsSlice, keySecret)
			}
		case "BasicAuthPasswordIdentityProvider":
			idP, certSecret, keySecret = buildBasicAuthIP(serializer, p)
			if certSecret != (secrets.Secret{}) {
				secretsSlice = append(secretsSlice, certSecret)
				secretsSlice = append(secretsSlice, keySecret)
			}
		default:
			logrus.Infof("Can't handle %s OAuth kind", kind)
			continue
		}
		oauthCrd.Spec.IdentityProviders = append(oauthCrd.Spec.IdentityProviders, idP)

		if secret.Metadata.Name != httpdAuthSecret || p.Kind == httpProviderKind {
			secretsSlice = append(secretsSlice, secret)
		}
	}

	return &oauthCrd, secretsSlice, nil
}

// GenYAML returns a YAML of the CRD
func (oauth *CRD) GenYAML() []byte {
	yamlBytes, err := yaml.Marshal(&oauth)
	if err != nil {
		logrus.Fatal(err)
	}

	return yamlBytes
}

package oauth

import (
	"github.com/fusor/cpma/ocp4/secrets"
	configv1 "github.com/openshift/api/legacyconfig/v1"
	oauthv1 "github.com/openshift/api/oauth/v1"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	// TODO: Is this line needed at all? It may be superflous to
	// ocp3.go/init()/configv1.InstallLegacy(scheme.Scheme)
	oauthv1.Install(scheme.Scheme)
}

// reference:
//   [v3] OCPv3:
//   - [1] https://docs.openshift.com/container-platform/3.11/install_config/configuring_authentication.html#identity_providers_master_config
//   [v4] OCPv4:
//   - [2] htpasswd: https://docs.openshift.com/container-platform/4.0/authentication/understanding-identity-provider.html
//   - [3] github: https://docs.openshift.com/container-platform/4.0/authentication/identity_providers/configuring-github-identity-provider.html

// Shared CRD part, present in all types of OAuth CRDs
type OAuthCRD struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	MetaData   MetaData `yaml:"metaData"`
	Spec       struct {
		IdentityProviders []interface{} `yaml:"identityProviders"`
	} `yaml:"spec"`
}

type MetaData struct {
	Name      string `yaml:"name"`
	NameSpace string `yaml:"namespace"`
}

var APIVersion = "config.openshift.io/v1"

// Translate converts OCPv3 OAuth to OCPv4 OAuth Custom Resources
func Translate(oauthconfig *configv1.OAuthConfig) (*OAuthCRD, []secrets.Secret, error) {
	var auth = oauthconfig.DeepCopy()
	var err error

	var oauthCrd OAuthCRD
	oauthCrd.APIVersion = APIVersion
	oauthCrd.Kind = "OAuth"
	oauthCrd.MetaData.Name = "cluster"
	oauthCrd.MetaData.NameSpace = "openshift-config"
	var secrets []secrets.Secret

	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	for _, p := range auth.IdentityProviders {
		p.Provider.Object, _, err = serializer.Decode(p.Provider.Raw, nil, nil)
		if err != nil {
			return nil, nil, err
		}

		switch kind := p.Provider.Object.GetObjectKind().GroupVersionKind().Kind; kind {
		case "GitHubIdentityProvider":
			idP, secret := buildGitHubIP(serializer, p)
			oauthCrd.Spec.IdentityProviders = append(oauthCrd.Spec.IdentityProviders, idP)
			secrets = append(secrets, secret)
		case "GitLabIdentityProvider":
			idP, secret := buildGitLabIP(serializer, p)
			oauthCrd.Spec.IdentityProviders = append(oauthCrd.Spec.IdentityProviders, idP)
			secrets = append(secrets, secret)
		case "HTPasswdPasswordIdentityProvider":
			idP, secret := buildHTPasswdIP(serializer, p)
			oauthCrd.Spec.IdentityProviders = append(oauthCrd.Spec.IdentityProviders, idP)
			secrets = append(secrets, secret)
		default:
			logrus.Infof("Can't handle %s OAuth kind", kind)
		}
	}

	return &oauthCrd, secrets, nil
}

// GenYAML returns a YAML of the OAuthCRD
func (oauth *OAuthCRD) GenYAML() string {
	yamlBytes, err := yaml.Marshal(&oauth)
	if err != nil {
		logrus.Fatal(err)
	}

	return string(yamlBytes)
}

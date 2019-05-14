package transform

import (
	"encoding/json"

	"github.com/fusor/cpma/pkg/transform/oauth"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

// OAuthExtraction holds OAuth data extracted from OCP3
type OAuthExtraction struct {
	IdentityProviders []oauth.IdentityProvider
}

// OAuthTransform is an OAuth specific transform
type OAuthTransform struct {
	Config *Config
}

// Transform converts data collected from an OCP3 cluster to OCP4 CR's
func (e OAuthExtraction) Transform() (Output, error) {
	logrus.Info("OAuthTransform::Transform")

	var ocp4Cluster Cluster

	oauth, secrets, err := oauth.Translate(e.IdentityProviders)
	if err != nil {
		logrus.WithError(err).Fatalf("Unable to generate OAuth CRD")
	}

	ocp4Cluster.Master.OAuth = *oauth
	ocp4Cluster.Master.Secrets = secrets

	var manifests []Manifest
	if ocp4Cluster.Master.OAuth.Kind != "" {
		manifest := Manifest{Name: "100_CPMA-cluster-config-oauth.yaml", CRD: ocp4Cluster.Master.OAuth.GenYAML()}
		manifests = append(manifests, manifest)

		for _, secret := range ocp4Cluster.Master.Secrets {
			filename := "100_CPMA-cluster-config-secret-" + secret.Metadata.Name + ".yaml"
			m := Manifest{Name: filename, CRD: secret.GenYAML()}
			manifests = append(manifests, m)
		}
	}

	return ManifestOutput{Manifests: manifests}, nil
}

// Extract collects OAuth configuration from an OCP3 cluster
func (e OAuthTransform) Extract() Extraction {
	logrus.Info("OAuthTransform::Extract")
	content := e.Config.Fetch(e.Config.MasterConfigFile)
	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	var masterConfig configv1.MasterConfig
	var extraction OAuthExtraction
	var htContent, crtContent, keyContent []byte

	_, _, err := serializer.Decode(content, nil, &masterConfig)
	if err != nil {
		HandleError(err)
	}

	if masterConfig.OAuthConfig != nil {
		for _, identityProvider := range masterConfig.OAuthConfig.IdentityProviders {
			providerJSON, _ := identityProvider.Provider.MarshalJSON()
			provider := oauth.Provider{}
			json.Unmarshal(providerJSON, &provider)
			if provider.Kind == "HTPasswdPasswordIdentityProvider" {
				htContent = e.Config.Fetch(provider.File)
			}
			if provider.Kind == "KeystonePasswordIdentityProvider" {
				crtContent = e.Config.Fetch(provider.CertFile)
				keyContent = e.Config.Fetch(provider.KeyFile)
			}

			extraction.IdentityProviders = append(extraction.IdentityProviders,
				oauth.IdentityProvider{
					provider.Kind,
					provider.APIVersion,
					identityProvider.MappingMethod,
					identityProvider.Name,
					identityProvider.Provider,
					provider.File,
					htContent,
					crtContent,
					keyContent,
					identityProvider.UseAsChallenger,
					identityProvider.UseAsLogin,
				})
		}
	}

	return extraction
}

// Validate confirms we have recieved good OAuth configuration data during Extract
func (e OAuthExtraction) Validate() error {
	logrus.Warn("Oauth Transform Validation Not Implmeneted")
	return nil // Simulate fine
}

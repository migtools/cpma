package transform

import (
	"encoding/json"

	"github.com/fusor/cpma/pkg/transform/oauth"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

type OAuthExtraction struct {
	IdentityProviders []oauth.IdentityProvider
}

type OAuthTransform struct {
	Config *Config
}

func (e OAuthExtraction) Transform() (TransformOutput, error) {
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

	return ManifestTransformOutput{Manifests: manifests}, nil
}

func (c OAuthTransform) Extract() Extraction {
	logrus.Info("OAuthTransform::Extract")
	content := c.Config.Fetch(c.Config.MasterConfigFile)
	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	var masterConfig configv1.MasterConfig
	var extraction OAuthExtraction
	var htContent []byte

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
				htContent = c.Config.Fetch(provider.File)
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
					identityProvider.UseAsChallenger,
					identityProvider.UseAsLogin,
				})
		}
	}

	return extraction
}

func (c OAuthExtraction) Validate() error {
	return nil // Simulate fine
}

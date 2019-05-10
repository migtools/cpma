package ocp

import (
	"encoding/json"

	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/fusor/cpma/pkg/ocp4"
	"github.com/fusor/cpma/pkg/ocp4/oauth"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

type OAuthTransform struct {
	Config *Config
}

func (c OAuthTransform) Run(content []byte) (TransformOutput, error) {
	logrus.Info("OAuthTransform::Run")

	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	var masterConfig configv1.MasterConfig
	var identityProviders []ocp3.IdentityProvider
	var ocp4Cluster ocp4.Cluster
	var htContent []byte

	_, _, err := serializer.Decode(content, nil, &masterConfig)
	if err != nil {
		HandleError(err)
	}

	for _, identityProvider := range masterConfig.OAuthConfig.IdentityProviders {
		providerJSON, _ := identityProvider.Provider.MarshalJSON()
		provider := Provider{}
		json.Unmarshal(providerJSON, &provider)
		if provider.Kind == "HTPasswdPasswordIdentityProvider" {
			htContent = c.Config.Fetch(provider.File)
		}

		identityProviders = append(identityProviders,
			ocp3.IdentityProvider{
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

	if masterConfig.OAuthConfig != nil {
		oauth, secrets, err := oauth.Translate(identityProviders)
		if err != nil {
			logrus.WithError(err).Fatalf("Unable to generate OAuth CRD from %+v", masterConfig.OAuthConfig)
		}

		ocp4Cluster.Master.OAuth = *oauth
		ocp4Cluster.Master.Secrets = secrets
	}

	var manifests []ocp4.Manifest
	if ocp4Cluster.Master.OAuth.Kind != "" {
		manifest := ocp4.Manifest{Name: "100_CPMA-cluster-config-oauth.yaml", CRD: ocp4Cluster.Master.OAuth.GenYAML()}
		manifests = append(manifests, manifest)

		for _, secret := range ocp4Cluster.Master.Secrets {
			filename := "100_CPMA-cluster-config-secret-" + secret.Metadata.Name + ".yaml"
			m := ocp4.Manifest{Name: filename, CRD: secret.GenYAML()}
			manifests = append(manifests, m)
		}
	}

	return ManifestTransformOutput{Config: *c.Config, Manifests: manifests}, nil
}

func (c OAuthTransform) Extract() []byte {
	logrus.Info("OAuthTransform::Extract")
	return c.Config.Fetch(c.Config.MasterConfigFile)
}

func (c OAuthTransform) Validate() error {
	return nil // Simulate fine
}

package ocp

import (
	"encoding/json"
	"fmt"
	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/fusor/cpma/pkg/ocp4"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

type OAuthTransform struct {
	ConfigFile *ocp3.ConfigFile
	Migration  *Migration
}

func (m OAuthTransform) Run() (TransformOutput, error) {
	fmt.Println("OAuthTransform::Run")
	//fmt.Println(m.ConfigFile.Content)
	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	_, _, err := serializer.Decode(m.ConfigFile.Content, nil, &m.Migration.OCP3Cluster.MasterConfig)
	if err != nil {
		HandleError(err)
	}

	for _, identityProvider := range m.Migration.OCP3Cluster.MasterConfig.OAuthConfig.IdentityProviders {
		providerJSON, _ := identityProvider.Provider.MarshalJSON()
		provider := Provider{}
		json.Unmarshal(providerJSON, &provider)
		var HTFile ocp3.ConfigFile
		if provider.Kind == "HTPasswdPasswordIdentityProvider" {
			HTFile = (ocp3.ConfigFile{"htpasswd", provider.File, nil})
			m.Migration.Fetch(&HTFile)
		}

		m.Migration.OCP3Cluster.IdentityProviders = append(m.Migration.OCP3Cluster.IdentityProviders,
			ocp3.IdentityProvider{
				provider.Kind,
				provider.APIVersion,
				identityProvider.MappingMethod,
				identityProvider.Name,
				identityProvider.Provider,
				HTFile.Path,
				HTFile.Content,
				identityProvider.UseAsChallenger,
				identityProvider.UseAsLogin,
			})
	}

	m.Migration.OCP4Cluster.Master.Translate(m.Migration.OCP3Cluster)

	var manifests []ocp4.Manifest
	if m.Migration.OCP4Cluster.Master.OAuth.Kind != "" {
		manifest := ocp4.Manifest{Name: "100_CPMA-cluster-config-oauth.yaml", CRD: m.Migration.OCP4Cluster.Master.OAuth.GenYAML()}
		manifests = append(manifests, manifest)

		for _, secret := range m.Migration.OCP4Cluster.Master.Secrets {
			filename := "100_CPMA-cluster-config-secret-" + secret.Metadata.Name + ".yaml"
			m := ocp4.Manifest{Name: filename, CRD: secret.GenYAML()}
			manifests = append(manifests, m)
		}
	}

	return ManifestTransformOutput{
		Migration: *m.Migration,
		Manifests: manifests,
	}, nil
}

func (m OAuthTransform) Extract() {
	fmt.Println("OAuthTransform::Extract")
	m.Migration.Fetch(m.ConfigFile)
	//fmt.Println(m.ConfigFile.Content)
}

func (m OAuthTransform) Validate() error {
	return nil // Simulate fine
}

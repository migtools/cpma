package ocp

import (
	"encoding/json"
	"fmt"
	"github.com/fusor/cpma/pkg/ocp3"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func (m MasterConfigTransform) Run() (TransformOutput, error) {
	fmt.Println("MasterConfigTransform::Run")
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

	return FileTransformOutput{
		FileData: "[MasterConfigTransform output file contents]\n",
	}, nil
}

func (m MasterConfigTransform) Extract() {
	fmt.Println("MasterConfigTransform::Extract")
	m.Migration.Fetch(m.ConfigFile)
	//fmt.Println(m.ConfigFile.Content)
}

func (m MasterConfigTransform) Validate() error {
	return nil // Simulate fine
}

func HandleError(err error) error {
	return fmt.Errorf("An error has occurred: %s", err)
}

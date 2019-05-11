package transform

import (
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Containers struct {
	Registries map[string]Registries
}

type Registries struct {
	List []string `toml:"registries"`
}

type ImageCR struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       struct {
		RegistrySources RegistrySources `yaml:"registrySources"`
	} `yaml:"spec"`
}

type Metadata struct {
	Name        string
	Annotations map[string]string `yaml:"annotations"`
}

type RegistrySources struct {
	BlockedRegistries  []string `yaml:"blockedRegistries,omitempty"`
	InsecureRegistries []string `yaml:"insecureRegistries,omitempty"`
}

type RegistriesTransform struct {
	Config *Config
}

func (c RegistriesTransform) Run(content []byte) (TransformOutput, error) {
	logrus.Info("RegistriesTransform::Run")

	var containers Containers
	var manifests Manifests

	if _, err := toml.Decode(string(content), &containers); err != nil {
		// handle error
	}

	const (
		apiVersion = "config.openshift.io/v1"
		kind       = "Image"
		name       = "cluster"
		annokey    = "release.openshift.io/create-only"
		annoval    = "true"
	)

	var imageCR ImageCR
	imageCR.APIVersion = apiVersion
	imageCR.Kind = kind
	imageCR.Metadata.Name = name
	imageCR.Metadata.Annotations = make(map[string]string)
	imageCR.Metadata.Annotations[annokey] = annoval
	imageCR.Spec.RegistrySources.BlockedRegistries = containers.Registries["block"].List
	imageCR.Spec.RegistrySources.InsecureRegistries = containers.Registries["insecure"].List

	imageCRYAML, err := yaml.Marshal(&imageCR)
	if err != nil {
		HandleError(err)
	}

	manifest := Manifest{Name: "100_CPMA-cluster-config-registries.yaml", CRD: imageCRYAML}
	manifests = append(manifests, manifest)

	return ManifestTransformOutput{
		Config:    *c.Config,
		Manifests: manifests,
	}, nil
}

func (c RegistriesTransform) Extract() []byte {
	logrus.Info("RegistriesTransform::Extract")
	return c.Config.Fetch(c.Config.RegistriesConfigFile)
}

func (c RegistriesTransform) Validate() error {
	return nil // Simulate fine
}

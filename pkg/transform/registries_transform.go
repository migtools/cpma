package transform

import (
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type RegistriesExtraction struct {
	Registries map[string]Registries
}

type Registries struct {
	List []string `toml:"registries"`
}

type ImageCR struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   Metadata  `yaml:"metadata"`
	Spec       ImageSpec `yaml:"spec"`
}

type Metadata struct {
	Name        string
	Annotations map[string]string `yaml:"annotations"`
}

type ImageSpec struct {
	RegistrySources RegistrySources `yaml:"registrySources"`
}

type RegistrySources struct {
	BlockedRegistries  []string `yaml:"blockedRegistries,omitempty"`
	InsecureRegistries []string `yaml:"insecureRegistries,omitempty"`
}

type RegistriesTransform struct {
	Config *Config
}

func (e RegistriesExtraction) Transform() (TransformOutput, error) {
	logrus.Info("RegistriesTransform::Extraction")
	var manifests []Manifest

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
	imageCR.Spec.RegistrySources.BlockedRegistries = e.Registries["block"].List
	imageCR.Spec.RegistrySources.InsecureRegistries = e.Registries["insecure"].List

	imageCRYAML, err := yaml.Marshal(&imageCR)
	if err != nil {
		HandleError(err)
	}

	manifest := Manifest{Name: "100_CPMA-cluster-config-registries.yaml", CRD: imageCRYAML}
	manifests = append(manifests, manifest)

	return ManifestTransformOutput{
		Manifests: manifests,
	}, nil
}

func (c RegistriesTransform) Extract() Extraction {
	logrus.Info("RegistriesTransform::Extract")
	content := c.Config.Fetch(c.Config.RegistriesConfigFile)
	var extraction RegistriesExtraction
	if _, err := toml.Decode(string(content), &extraction); err != nil {
		HandleError(err)
	}
	return extraction
}

func (c RegistriesExtraction) Validate() error {
	logrus.Warn("Registries Transform Validation Not Implmeneted")
	return nil // Simulate fine
}

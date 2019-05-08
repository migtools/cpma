package image

import (
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

const (
	apiVersion = "config.openshift.io/v1"
	kind       = "Image"
	name       = "cluster"
	annokey    = "release.openshift.io/create-only"
	annoval    = "true"
)

// Translate converts OCPv3 registries.conf to OCPv4 Image  Custom Resources
func Translate(containers Containers) *ImageCR {

	var imageCR ImageCR
	imageCR.APIVersion = apiVersion
	imageCR.Kind = kind
	imageCR.Metadata.Name = name
	imageCR.Metadata.Annotations = make(map[string]string)
	imageCR.Metadata.Annotations[annokey] = annoval
	imageCR.Spec.RegistrySources.BlockedRegistries = containers.Registries["block"].List
	imageCR.Spec.RegistrySources.InsecureRegistries = containers.Registries["insecure"].List
	return &imageCR
}

// GenYAML returns a YAML of the OAuthCRD
func (imageCR *ImageCR) GenYAML() []byte {
	yamlBytes, err := yaml.Marshal(&imageCR)
	if err != nil {
		logrus.Fatal(err)
	}

	return yamlBytes
}

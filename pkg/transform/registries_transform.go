package transform

import (
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/fusor/cpma/pkg/env"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// RegistriesExtraction holds registry information extracted from an OCP3 cluster
type RegistriesExtraction struct {
	Registries map[string]Registries
}

// Registries holds a list of Registries
type Registries struct {
	List []string `toml:"registries"`
}

// ImageCR is an Image Cluster Resource
type ImageCR struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   Metadata  `yaml:"metadata"`
	Spec       ImageSpec `yaml:"spec"`
}

// Metadata is the Metadata for an Image Cluster Resource
type Metadata struct {
	Name        string
	Annotations map[string]string `yaml:"annotations"`
}

// ImageSpec is a Spec for an ImageCR
type ImageSpec struct {
	RegistrySources RegistrySources `yaml:"registrySources"`
}

// RegistrySources holds lists of blocked and insecure registries from an OCP3 cluster
type RegistrySources struct {
	BlockedRegistries  []string `yaml:"blockedRegistries,omitempty"`
	InsecureRegistries []string `yaml:"insecureRegistries,omitempty"`
}

// RegistriesTransform is a registry specific transform
type RegistriesTransform struct {
	Config *Config
}

// Transform contains registry configuration collected from an OCP3 cluster
func (e RegistriesExtraction) Transform() (Output, error) {
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
		return nil, err
	}

	manifest := Manifest{Name: "100_CPMA-cluster-config-registries.yaml", CRD: imageCRYAML}
	manifests = append(manifests, manifest)

	return ManifestOutput{
		Manifests: manifests,
	}, nil
}

// Extract collects registry information from an OCP3 cluster
func (e RegistriesTransform) Extract() (Extraction, error) {
	logrus.Info("RegistriesTransform::Extract")
	content, err := e.Config.Fetch(env.Config().GetString("RegistriesConfigFile"))
	if err != nil {
		return nil, err
	}

	var extraction RegistriesExtraction
	if _, err := toml.Decode(string(content), &extraction); err != nil {
		return nil, err
	}
	return extraction, nil
}

// Validate registry data collected from an OCP3 cluster
func (e RegistriesExtraction) Validate() error {
	if len(e.Registries["block"].List) == 0 && len(e.Registries["insecure"].List) == 0 {
		return errors.New("no configured registries detected, not generating a cr")
	}
	return nil
}

// Type retrurn transform type
func (e RegistriesTransform) Type() string {
	return "Registries"
}

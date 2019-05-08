package ocp

import (
	"github.com/fusor/cpma/pkg/ocp4"
)

type Provider struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	File       string `json:"file"`
}

type Config struct {
	MasterConfigFile     string
	NodeConfigFile       string
	RegistriesConfigFile string
	OutputDir            string
	Hostname             string
}

type ManifestTransformOutput struct {
	Config    Config
	Manifests []ocp4.Manifest
}

type TransformRunner struct {
	Config string
}

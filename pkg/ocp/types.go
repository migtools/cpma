package ocp

import (
	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/fusor/cpma/pkg/ocp4"
)

type Migration struct {
	OCP3Cluster ocp3.Cluster
	OCP4Cluster ocp4.Cluster
	OutputDir   string
}

type Provider struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	File       string `json:"file"`
}

type Config struct {
	MasterConfigFile string
	NodeConfigFile   string
	RunnerConfig     string
}

type FileTransformOutput struct {
	FileData string
}

type MasterConfigTransform struct {
	ConfigFile *ocp3.ConfigFile
	Migration  *Migration
}

type TransformRunner struct {
	Config string
}

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

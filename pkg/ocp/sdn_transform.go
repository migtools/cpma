package ocp

import (
	"fmt"
	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/fusor/cpma/pkg/ocp4"
)

type SDNTransform struct {
	ConfigFile *ocp3.ConfigFile
	Migration  *Migration
}

func (m SDNTransform) Run() (TransformOutput, error) {
	fmt.Println("SDNTransform::Run")
	var manifests []ocp4.Manifest
	//create manifests from files
	return ManifestTransformOutput{
		Migration: *m.Migration,
		Manifests: manifests,
	}, nil
}

func (m SDNTransform) Extract() {
	fmt.Println("SDNTransform::Extract")
	// Retrieve file(s)... master-config, whatever
}

func (m SDNTransform) Validate() error {
	return nil // Simulate fine
}

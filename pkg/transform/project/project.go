package project

import (
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
)

const (
	apiVersion = "operator.openshift.io/v1"
	kind       = "Project"
	name       = "cluster"
)

// Translate ProjectPolicyConfig definitions
func Translate(projectConfig legacyconfigv1.ProjectConfig) (*configv1.Project, error) {
	var projectCR configv1.Project

	projectCR.APIVersion = apiVersion
	projectCR.Kind = kind
	projectCR.Name = name

	if projectConfig.ProjectRequestMessage != "" {
		projectCR.Spec.ProjectRequestMessage = projectConfig.ProjectRequestMessage
	}

	if projectConfig.ProjectRequestTemplate != "" {
		projectCR.Spec.ProjectRequestTemplate.Name = projectConfig.ProjectRequestTemplate
	}

	return &projectCR, nil
}

// Validate registry data collected from an OCP3 cluster
func Validate(e legacyconfigv1.ProjectConfig) error {
	return nil
}

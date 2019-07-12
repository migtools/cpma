package project

import (
	"strings"

	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/sirupsen/logrus"
)

const (
	apiVersion = "config.openshift.io/v1"
	kind       = "Project"
	name       = "cluster"
)

// Translate ProjectPolicyConfig definitions
func Translate(projectConfig legacyconfigv1.ProjectConfig) (*configv1.Project, error) {
	var projectCR configv1.Project

	projectCR.APIVersion = apiVersion
	projectCR.Kind = kind
	projectCR.Name = name

	if projectConfig.DefaultNodeSelector != "" {
		logrus.Info("ProjectConfig.DefaultNodeSelector is handled by scheduler")
	}

	if projectConfig.ProjectRequestMessage != "" {
		projectCR.Spec.ProjectRequestMessage = projectConfig.ProjectRequestMessage
	}

	if projectConfig.ProjectRequestTemplate != "" {
		i := strings.Index(projectConfig.ProjectRequestTemplate, "/")
		prefix := projectConfig.ProjectRequestTemplate[0 : i+1]
		projectRequestTemplate := strings.TrimLeft(projectConfig.ProjectRequestTemplate, prefix)
		projectCR.Spec.ProjectRequestTemplate.Name = projectRequestTemplate
	}

	return &projectCR, nil
}

// Validate registry data collected from an OCP3 cluster
func Validate(e legacyconfigv1.ProjectConfig) error {
	return nil
}

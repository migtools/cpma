package scheduler

import (
	"errors"

	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/sirupsen/logrus"
)

const (
	apiVersion = "operator.openshift.io/v1"
	kind       = "Scheduler"
	name       = "cluster"
)

// Translate ProjectPolicyConfig definitions
func Translate(masterConfig legacyconfigv1.MasterConfig) (*configv1.Scheduler, error) {
	var schedulerCR configv1.Scheduler

	schedulerCR.APIVersion = apiVersion
	schedulerCR.Kind = kind
	schedulerCR.Name = name

	if masterConfig.ProjectConfig.DefaultNodeSelector != "" {
		logrus.Info("ProjectConfig.DefaultNodeSelector is handled by scheduler")
		schedulerCR.Spec.DefaultNodeSelector = masterConfig.ProjectConfig.DefaultNodeSelector
	}

	return &schedulerCR, nil
}

// Validate registry data collected from an OCP3 cluster
func Validate(e legacyconfigv1.MasterConfig) error {
	if len(e.ProjectConfig.DefaultNodeSelector) == 0 {
		return errors.New("DefaultNodeSelector can't be empty")
	}

	return nil
}

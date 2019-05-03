package ocp3

import (
	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

func init() {
	configv1.InstallLegacy(scheme.Scheme)
}

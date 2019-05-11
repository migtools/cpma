package ocp3

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

type IdentityProvider struct {
	Kind            string
	APIVersion      string
	MappingMethod   string
	Name            string
	Provider        runtime.RawExtension
	HTFileName      string
	HTFileData      []byte
	UseAsChallenger bool
	UseAsLogin      bool
}

func init() {
	configv1.InstallLegacy(scheme.Scheme)
}

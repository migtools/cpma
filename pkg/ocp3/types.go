package ocp3

import (
	"k8s.io/apimachinery/pkg/runtime"
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

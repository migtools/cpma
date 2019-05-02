package ocp3

import (
	configv1 "github.com/openshift/api/legacyconfig/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Cluster struct {
	Hostname          string
	MasterConfig      configv1.MasterConfig
	NodeConfig        configv1.NodeConfig
	ETCDConfig        ETCDConfig
	RuntimeConfig     RuntimeConfig
	IdentityProviders []IdentityProvider
}

type ETCDConfig struct {
	TLSCipherSuites string
}

type RuntimeConfig struct {
	Type               string
	BlockedRegistries  []string
	SearchRegistries   []string
	InsecureRegistries []string
}

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

type HTPasswd struct {
	Data string //whatever actually goes here
}

type ConfigFile struct {
	Type    string
	Path    string
	Content []byte
}

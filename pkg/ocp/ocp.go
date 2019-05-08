package ocp

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/fusor/cpma/env"
	"github.com/fusor/cpma/internal/io"
	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/fusor/cpma/pkg/ocp4"
	"github.com/fusor/cpma/pkg/ocp4/oauth"
	"github.com/fusor/cpma/pkg/ocp4/sdn"
	"github.com/fusor/cpma/pkg/ocp4/secrets"
	configv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

type ConfigFile struct {
	Hostname string
	Path     string
	Content  []byte
}

type OAuthConfig struct {
	ConfigFile
	OCP3    configv1.MasterConfig
	OAuth   oauth.OAuthCRD
	Secrets []secrets.Secret
}

type SDNConfig struct {
	ConfigFile
	OCP3 configv1.MasterNetworkConfig
	SDN  sdn.NetworkCR
}

type ConfigNode struct {
	ConfigFile
	OCP3 ocp3.Node
	OCP4 ocp4.Node
}

type Translator interface {
	Add(string)
	Extract()
	Transform()
	GenYAML() ocp4.Manifests
}

// GetFile allows to mock file retrieval
var GetFile = io.GetFile

var source = env.Config().GetString("Source")

func (sdnConfig *SDNConfig) Add(hostname string) {
	path := env.Config().GetString("MasterConfigFile")

	if path == "" {
		path = "/etc/origin/master/master-config.yaml"
	}
	sdnConfig.ConfigFile.Hostname = hostname
	sdnConfig.ConfigFile.Path = path
}

func (oauthConfig *OAuthConfig) Add(hostname string) {
	path := env.Config().GetString("MasterConfigFile")

	if path == "" {
		path = "/etc/origin/master/master-config.yaml"
	}
	oauthConfig.ConfigFile.Hostname = hostname
	oauthConfig.ConfigFile.Path = path
}

func (config *ConfigNode) Add(hostname string) {
	nodef := env.Config().GetString("NodeConfigFile")

	if nodef == "" {
		nodef = "/etc/origin/node/node-config.yaml"
	}

	config.ConfigFile.Hostname = hostname
	config.ConfigFile.Path = nodef
}

func (sdnConfig *SDNConfig) Decode() {
	var masterConfig configv1.MasterConfig
	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	_, _, err := serializer.Decode(sdnConfig.Content, nil, &masterConfig)
	if err != nil {
		logrus.Fatal(err)
	}

	sdnConfig.OCP3 = masterConfig.NetworkConfig
}

// Decode unmarshals OCP3 MasterConfig and sets OAuth
func (oauthConfig *OAuthConfig) Decode() {
	var masterConfig configv1.MasterConfig
	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	_, _, err := serializer.Decode(oauthConfig.Content, nil, &masterConfig)
	if err != nil {
		logrus.Fatal(err)
	}

	oauthConfig.OCP3.OAuthConfig = masterConfig.OAuthConfig
}

func (config *ConfigNode) Decode() {
	config.OCP3.Decode(config.ConfigFile.Content)
}

// DumpManifests creates OCDs files
func DumpManifests(manifests ocp4.Manifests) {
	for _, manifest := range manifests {
		maniftestfile := filepath.Join(env.Config().GetString("OutputDir"), "manifests", manifest.Name)
		os.MkdirAll(path.Dir(maniftestfile), 0755)
		err := ioutil.WriteFile(maniftestfile, manifest.CRD, 0644)
		logrus.Printf("CRD:Added: %s", maniftestfile)
		if err != nil {
			logrus.Panic(err)
		}
	}
}

func Fetch(configFile *ConfigFile) {
	localF := filepath.Join(env.Config().GetString("OutputDir"), configFile.Hostname, configFile.Path)
	configFile.Content = GetFile(configFile.Hostname, configFile.Path, localF)
	logrus.Printf("File:Loaded: %s", localF)
}

// GenYAML returns the list of OAuth CRDs
func (oauthConfig *OAuthConfig) GenYAML() ocp4.Manifests {
	var manifests ocp4.Manifests

	// Generate yaml for oauth config
	crd := oauthConfig.OAuth.GenYAML()
	manifests = ocp4.OAuthManifest(oauthConfig.OAuth.Kind, crd, manifests)

	for _, secretManifest := range oauthConfig.Secrets {
		crd := secretManifest.GenYAML()
		manifests = ocp4.SecretsManifest(secretManifest, crd, manifests)
	}

	return manifests
}

// GenYAML returns the list of translated CRDs
func (sdnConfig *SDNConfig) GenYAML() ocp4.Manifests {
	var manifests ocp4.Manifests
	crd := sdnConfig.SDN.GenYAML()
	manifests = ocp4.SDNManifest(crd, manifests)
	return manifests
}

// GenYAML returns the list of translated CRDs
func (config *ConfigNode) GenYAML() ocp4.Manifests {
	var manifests ocp4.Manifests

	nodeManifests := config.OCP4.GenYAML()

	for _, manifest := range nodeManifests {
		manifests = append(manifests, manifest)
	}
	return manifests
}

// Extract fetch then decode OCP3 OAuth component
func (oauthConfig *OAuthConfig) Extract() {
	Fetch(&oauthConfig.ConfigFile)
	oauthConfig.Decode()
}

// Extract fetch then decode OCP3 component
func (sdnConfig *SDNConfig) Extract() {
	Fetch(&sdnConfig.ConfigFile)
	sdnConfig.Decode()
}

// Extract fetch then decode OCP3 component
func (config *ConfigNode) Extract() {
	Fetch(&config.ConfigFile)
	config.Decode()
}

// Transform OCP3 to OCP4
func (oauthConfig *OAuthConfig) Transform() {
	if oauthConfig.OCP3.OAuthConfig != nil {
		logrus.Debugln("Transforming oauth config")
		oauth, secretList, err := oauth.Transform(oauthConfig.OCP3.OAuthConfig)

		if err != nil {
			logrus.WithError(err).Fatalf("Unable to generate OAuth CRD from %+v", oauthConfig.OCP3.OAuthConfig)
		}
		oauthConfig.OAuth = *oauth
		oauthConfig.Secrets = secretList
	}
}

func (sdnConfig *SDNConfig) Transform() {
	if &sdnConfig.OCP3 != nil {
		logrus.Debugln("Translating SDN config")
		networkCR := sdn.Transform(sdnConfig.OCP3)
		sdnConfig.SDN = *networkCR
	}
}

// Transform OCP3 to OCP4
func (config *ConfigNode) Transform() {
	config.OCP4.Transform(config.OCP3.Config)
}

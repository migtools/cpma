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
)

type ConfigFile struct {
	Hostname string
	Path     string
	Content  []byte
}

type OAuthTranslator struct {
	ConfigFile
	OCP3    configv1.MasterConfig
	OAuth   oauth.OAuthCRD
	Secrets []secrets.Secret
}

type SDNTranslator struct {
	ConfigFile
	OCP3 configv1.MasterNetworkConfig
	SDN  sdn.NetworkCR
}

type Translator interface {
	Add(string)
	Extract()
	Transform()
	Load()
}

// GetFile allows to mock file retrieval
var GetFile = io.GetFile

var source = env.Config().GetString("Source")

func (sdnConfig *SDNTranslator) Add(hostname string) {
	path := env.Config().GetString("MasterConfigFile")

	if path == "" {
		path = "/etc/origin/master/master-config.yaml"
	}
	sdnConfig.ConfigFile.Hostname = hostname
	sdnConfig.ConfigFile.Path = path
}

func (oauthConfig *OAuthTranslator) Add(hostname string) {
	path := env.Config().GetString("MasterConfigFile")

	if path == "" {
		path = "/etc/origin/master/master-config.yaml"
	}
	oauthConfig.ConfigFile.Hostname = hostname
	oauthConfig.ConfigFile.Path = path
}

// Decode unmarshals OCP3 MasterConfig and sets OAuth
func (oauthConfig *OAuthTranslator) Decode() {
	masterConfig := ocp3.MasterDecode(oauthConfig.Content)
	oauthConfig.OCP3.OAuthConfig = masterConfig.OAuthConfig
}

func (sdnConfig *SDNTranslator) Decode() {
	masterConfig := ocp3.MasterDecode(sdnConfig.Content)
	sdnConfig.OCP3 = masterConfig.NetworkConfig
}

// DumpManifests creates Manifests file from OCDs
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

func fetch(configFile *ConfigFile) {
	localF := filepath.Join(env.Config().GetString("OutputDir"), configFile.Hostname, configFile.Path)
	configFile.Content = GetFile(configFile.Hostname, configFile.Path, localF)
	logrus.Printf("File:Loaded: %s", localF)
}

// Extract fetch then decode OCP3 OAuth component
func (oauthConfig *OAuthTranslator) Extract() {
	fetch(&oauthConfig.ConfigFile)
	oauthConfig.Decode()
}

// Extract fetch then decode OCP3 component
func (sdnConfig *SDNTranslator) Extract() {
	fetch(&sdnConfig.ConfigFile)
	sdnConfig.Decode()
}

func (oauthConfig *OAuthTranslator) Load() {
	var manifests ocp4.Manifests

	// Generate yaml for oauth config
	crd := oauthConfig.OAuth.GenYAML()
	manifests = ocp4.OAuthManifest(oauthConfig.OAuth.Kind, crd, manifests)

	for _, secretManifest := range oauthConfig.Secrets {
		crd := secretManifest.GenYAML()
		manifests = ocp4.SecretsManifest(secretManifest, crd, manifests)
	}

	DumpManifests(manifests)
}

func (sdnConfig *SDNTranslator) Load() {
	var manifests ocp4.Manifests
	crd := sdnConfig.SDN.GenYAML()
	manifests = ocp4.SDNManifest(crd, manifests)
	DumpManifests(manifests)
}

// Transform OAuthTranslator from OCP3 to OCP4
func (oauthConfig *OAuthTranslator) Transform() {
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

// Transform SDNTranslator from OCP3 to OCP4
func (sdnConfig *SDNTranslator) Transform() {
	if &sdnConfig.OCP3 != nil {
		logrus.Debugln("Translating SDN config")
		networkCR := sdn.Transform(sdnConfig.OCP3)
		sdnConfig.SDN = *networkCR
	}
}

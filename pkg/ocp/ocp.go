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
	Extract()
	Transform()
	Load()
}

// GetFile allows to mock file retrieval
var GetFile = io.GetFile

var source = env.Config().GetString("Source")

func NewOAuthTranslator(hostname string) *OAuthTranslator {
	oauthTranslator := new(OAuthTranslator)
	path := env.Config().GetString("MasterConfigFile")

	if path == "" {
		path = "/etc/origin/master/master-config.yaml"
	}
	oauthTranslator.ConfigFile.Hostname = hostname
	oauthTranslator.ConfigFile.Path = path
	return oauthTranslator
}

func NewSDNTranslator(hostname string) *SDNTranslator {
	sdnTranslator := new(SDNTranslator)
	path := env.Config().GetString("MasterConfigFile")

	if path == "" {
		path = "/etc/origin/master/master-config.yaml"
	}
	sdnTranslator.ConfigFile.Hostname = hostname
	sdnTranslator.ConfigFile.Path = path
	return sdnTranslator
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
func (oauthTranslator *OAuthTranslator) Extract() {
	fetch(&oauthTranslator.ConfigFile)
	masterConfig := ocp3.MasterDecode(oauthTranslator.Content)
	oauthTranslator.OCP3.OAuthConfig = masterConfig.OAuthConfig
}

// Extract fetch then decode OCP3 component
func (sdnTranslator *SDNTranslator) Extract() {
	fetch(&sdnTranslator.ConfigFile)
	masterConfig := ocp3.MasterDecode(sdnTranslator.Content)
	sdnTranslator.OCP3 = masterConfig.NetworkConfig
}

func (oauthTranslator *OAuthTranslator) Load() {
	var manifests ocp4.Manifests

	// Generate yaml for oauth config
	crd := oauthTranslator.OAuth.GenYAML()
	manifests = ocp4.OAuthManifest(oauthTranslator.OAuth.Kind, crd, manifests)

	for _, secretManifest := range oauthTranslator.Secrets {
		crd := secretManifest.GenYAML()
		manifests = ocp4.SecretsManifest(secretManifest, crd, manifests)
	}

	DumpManifests(manifests)
}

func (sdnTranslator *SDNTranslator) Load() {
	var manifests ocp4.Manifests
	crd := sdnTranslator.SDN.GenYAML()
	manifests = ocp4.SDNManifest(crd, manifests)
	DumpManifests(manifests)
}

// Transform OAuthTranslator from OCP3 to OCP4
func (oauthTranslator *OAuthTranslator) Transform() {
	if oauthTranslator.OCP3.OAuthConfig != nil {
		logrus.Debugln("Transforming oauth config")
		oauth, secretList, err := oauth.Transform(oauthTranslator.OCP3.OAuthConfig)

		if err != nil {
			logrus.WithError(err).Fatalf("Unable to generate OAuth CRD from %+v", oauthTranslator.OCP3.OAuthConfig)
		}
		oauthTranslator.OAuth = *oauth
		oauthTranslator.Secrets = secretList
	}
}

// Transform SDNTranslator from OCP3 to OCP4
func (sdnTranslator *SDNTranslator) Transform() {
	if &sdnTranslator.OCP3 != nil {
		logrus.Debugln("Translating SDN config")
		networkCR := sdn.Transform(sdnTranslator.OCP3)
		sdnTranslator.SDN = *networkCR
	}
}

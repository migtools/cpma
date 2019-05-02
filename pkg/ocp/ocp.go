package ocp

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/fusor/cpma/pkg/ocp4"
	"github.com/fusor/cpma/pkg/sftpclient"
	"github.com/sirupsen/logrus"
)

const MasterConfigFile = "/etc/origin/master/master-config.yaml"
const NodeConfigFile = "/etc/origin/node/node-config.yaml"
const ETCDConfigFile = "/etc/etcd/etcd.conf"

func (migration *Migration) Decode(configFile ocp3.ConfigFile) {
	migration.OCP3Cluster.Decode(configFile)
}

// GenYAML returns the list of translated CRDs
func (migration *Migration) GenYAML() ocp4.Manifests {
	var manifests ocp4.Manifests

	masterManifests, err := migration.OCP4Cluster.Master.GenYAML()
	if err != nil {
		return nil
	}
	for _, manifest := range masterManifests {
		manifests = append(manifests, manifest)
	}
	return manifests
}

// DumpManifests creates OCDs files
func (migration *Migration) DumpManifests(manifests []ocp4.Manifest) {
	for _, manifest := range manifests {
		maniftestfile := filepath.Join(migration.OutputDir, "manifests", manifest.Name)
		os.MkdirAll(path.Dir(maniftestfile), 0755)
		err := ioutil.WriteFile(maniftestfile, manifest.CRD, 0644)
		logrus.Printf("CR manifest created: %s", maniftestfile)
		if err != nil {
			logrus.Panic(err)
		}
	}
}

// Fetch retrieves file from Host
func (migration *Migration) Fetch(configFile *ocp3.ConfigFile) {
	dst := filepath.Join(migration.OutputDir, migration.OCP3Cluster.Hostname, configFile.Path)
	sftpclient.Fetch(migration.OCP3Cluster.Hostname, configFile.Path, dst)

	f, err := ioutil.ReadFile(filepath.Join(migration.OutputDir, migration.OCP3Cluster.Hostname, configFile.Path))
	if err != nil {
		logrus.Warning(err)
	}
	configFile.Content = f
}

func (migration *Migration) LoadOCP3Configs() {
	// Get master config so we can start figuring out what additional files we need
	masterConfig := ocp3.ConfigFile{"master", "/etc/origin/master/master-config.yaml", nil}
	migration.Fetch(&masterConfig)
	migration.Decode(masterConfig)

	// Start compiling a list of additional files to retrieve
	configFiles := []ocp3.ConfigFile{}
	configFiles = append(configFiles, ocp3.ConfigFile{"node", "/etc/origin/node/node-config.yaml", nil})
	configFiles = append(configFiles, ocp3.ConfigFile{"etcd", "/etc/etcd/etcd.conf", nil})
	//configFiles = append(configFiles, ocp3.ConfigFile{"crio", "/etc/crio/crio.conf", nil})

	for _, identityProvider := range migration.OCP3Cluster.MasterConfig.OAuthConfig.IdentityProviders {
		providerJSON, _ := identityProvider.Provider.MarshalJSON()
		provider := Provider{}
		json.Unmarshal(providerJSON, &provider)
		var HTFile ocp3.ConfigFile
		if provider.Kind == "HTPasswdPasswordIdentityProvider" {
			HTFile = (ocp3.ConfigFile{"htpasswd", provider.File, nil})
			migration.Fetch(&HTFile)
		}

		migration.OCP3Cluster.IdentityProviders = append(migration.OCP3Cluster.IdentityProviders,
			ocp3.IdentityProvider{
				provider.Kind,
				provider.APIVersion,
				identityProvider.MappingMethod,
				identityProvider.Name,
				identityProvider.Provider,
				HTFile.Path,
				HTFile.Content,
				identityProvider.UseAsChallenger,
				identityProvider.UseAsLogin,
			})
	}

	for _, configFile := range configFiles {
		migration.Fetch(&configFile)
		migration.Decode(configFile)
	}

}

// Translate OCP3 to OCP4
func (migration *Migration) Translate() {
	migration.OCP4Cluster.Master.Translate(migration.OCP3Cluster)
}

package env

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/fusor/cpma/internal/config"
	"github.com/fusor/cpma/internal/sftpclient"
	v1 "github.com/openshift/origin/pkg/cmd/server/apis/config/v1"
	log "github.com/sirupsen/logrus"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
)

type ClusterV3 struct {
	Nodes []NodeConfig
}

type ClusterCMD interface {
	addNode(name, file, path string) int
	Show() string
}

type Cmd interface {
	FetchSrc() int
	Show() string
}

type Info struct {
	SrCluster  ClusterV3       `mapstructure:"srcluster"`
	DsCluster  ClusterV4       `mapstructure:"dscluster"`
	SFTP       sftpclient.Info `mapstructure:"Source"`
	OutputPath string          `mapstructure:"outputPath"`
	LocalOnly  bool            `mapstructure:"localOnly"`
}

type NodeConfig struct {
	Name             string
	FileName         string
	Path             string
	Payload          []byte
	MstConfig        *V1MasterConfig
	NdeConfig        *V1NodeConfig
	OAuthIDProviders []IdentityProviders
}

type ClusterV4 struct {
	CRManifests []*CRManifest
}

type CRManifest struct {
	Kind         string
	ManifestFile string
}

type V1MasterConfig v1.MasterConfig

type V1NodeConfig v1.NodeConfig

type IdentityProviders struct {
	Name          string
	Challenge     bool
	Login         bool
	MappingMethod string
	Provider      ProviderInfo
}

type ProviderInfo struct {
	APIVersion    string   `json:"apiVersion"`
	Ca            string   `json:ca`
	ClientID      string   `json:"clientID"`
	ClientSecret  string   `json:"clientSecret"`
	File          string   `json:file`
	Hostname      string   `json:hostname`
	Kind          string   `json:"kind"`
	Organizations []string `json:organizations`
	Teams         []string `json:teams`
}

func addNode(name, filename, path string) *NodeConfig {
	x := NodeConfig{}
	x.Name = name
	x.FileName = filename
	x.Path = path
	return &x
}

func (config *Info) CROAuth() {
	for i := range config.SrCluster.Nodes {
		if config.SrCluster.Nodes[i].MstConfig != nil {
			for j := range config.SrCluster.Nodes[i].OAuthIDProviders {
				val := config.SrCluster.Nodes[i].OAuthIDProviders[j]
				switch {
				case strings.Contains(val.Name, "htpasswd_auth"):
					config.CROAuthHTPasswd(val)
				case strings.Contains(val.Name, "github"):
					config.CROAuthGithub(val)
				case strings.Contains(val.Name, "google"):
					config.CROAuthGoogle(val)
				}
			}
		}
	}
}

func (config *Info) CROAuthGoogle(identityProvider IdentityProviders) {
}

func (config *Info) CROAuthGithub(identityProvider IdentityProviders) {
	const templ1 = `apiVersion: config.openshift.io/v1
kind: OAuth
metadata:
	name: cluster
spec:
	identityProviders:
	- name: {{.Name}}
		challenge: false
		login: {{.Login}}
		mappingMethod: {{.MappingMethod}}
		type: GitHub
		github:
			{{- if .Provider.Ca}}
			ca:
				name: {{.Provider.Ca}}
			{{- end}}
			clientID: {{.Provider.ClientID}}
			clientSecret:
				name: github-secret
			{{- if .Provider.Hostname}} hostname: {{.Provider.Hostname}}{{end}}
			{{- if .Provider.Organizations }}
			organizations:
				{{- range .Provider.Organizations}}
				- {{.}}
				{{- end}}
			{{- end}}
			{{- if .Provider.Teams }}
			teams:
				{{- range .Provider.Teams}}
				- {{.}}
				{{- end}}
			{{- end}}`

	var manifest1 = template.Must(template.New("OAuthCRGithub1").Parse(templ1))
	result := new(CRManifest)
	result.Kind = "OAuth"

	newpath := filepath.Join(config.OutputPath, "destination/oauth")
	os.MkdirAll(newpath, os.ModePerm)
	result.ManifestFile = filepath.Join(newpath, "github.yaml")

	manifestfile, err := os.Create(result.ManifestFile)
	if err != nil {
		log.Fatal(err)
	}
	defer manifestfile.Close()

	if err := manifest1.Execute(manifestfile, identityProvider); err != nil {
		log.Fatal(err)
	}

	config.DsCluster.CRManifests = append(config.DsCluster.CRManifests, result)
}

func (config *Info) CROAuthHTPasswd(identityProvider IdentityProviders) {
	const templ1 = `apiVersion: config.openshift.io/v1
kind: OAuth
metadata:
  name: cluster
spec:
  identityProviders:
  - name: my_htpasswd_provider
    challenge: {{.Challenge}}
    login: {{.Login}}
    mappingMethod: {{.MappingMethod}}
    type: HTPasswd
    htpasswd:
      fileData:
        name: htpass-secret`

	var manifest1 = template.Must(template.New("OAuthCR").Parse(templ1))

	result := new(CRManifest)
	result.Kind = "OAuth"

	newpath := filepath.Join(config.OutputPath, "destination/oauth")
	os.MkdirAll(newpath, os.ModePerm)
	result.ManifestFile = filepath.Join(newpath, "htpassword.yaml")

	manifestfile, err := os.Create(result.ManifestFile)
	if err != nil {
		log.Fatal(err)
	}
	defer manifestfile.Close()

	if err := manifest1.Execute(manifestfile, identityProvider); err != nil {
		log.Fatal(err)
	}

	config.DsCluster.CRManifests = append(config.DsCluster.CRManifests, result)
}

func (config *Info) FetchSrc() int {
	sftpclient := config.SFTP.NewClient()
	defer sftpclient.Close()

	for i := range config.SrCluster.Nodes {
		srcFilePath := filepath.Join(config.SrCluster.Nodes[i].Path, config.SrCluster.Nodes[i].FileName)
		dstFilePath := filepath.Join(config.OutputPath, "source", srcFilePath)

		sftpclient.GetFile(srcFilePath, dstFilePath)

		payload, err := ioutil.ReadFile(dstFilePath)
		if err != nil {
			log.Fatal(err)
		}
		config.SrCluster.Nodes[i].Payload = payload
	}

	return 1
}

func (cluster *ClusterV3) load(list [][]string) {
	for _, nc := range list {
		cluster.Nodes = append(cluster.Nodes, *addNode(nc[0], nc[1], nc[2]))
	}
}

func (config *Info) LoadSrc() int {
	for i := range config.SrCluster.Nodes {
		srcFilePath := filepath.Join(config.SrCluster.Nodes[i].Path, config.SrCluster.Nodes[i].FileName)
		dstFilePath := filepath.Join(config.OutputPath, "source", srcFilePath)

		payload, err := ioutil.ReadFile(dstFilePath)
		if err != nil {
			log.Fatal(err)
		}
		config.SrCluster.Nodes[i].Payload = payload
	}

	return 1
}

func (config *Info) Parse() {
	for i := range config.SrCluster.Nodes {
		if config.SrCluster.Nodes[i].Name == "master" {
			config.SrCluster.Nodes[i].MstConfig = config.SrCluster.Nodes[i].ParseMaster(config)

			// Workaround missing plugin extension context
			for j := range config.SrCluster.Nodes[i].MstConfig.OAuthConfig.IdentityProviders {
				if fmt.Sprint(reflect.TypeOf(config.SrCluster.Nodes[i].MstConfig.OAuthConfig.IdentityProviders[j].Provider)) == "runtime.RawExtension" {
					jsonProvider := string(config.SrCluster.Nodes[i].MstConfig.OAuthConfig.IdentityProviders[j].Provider.Raw)

					var identityProvider = IdentityProviders{}
					identityProvider.Name = config.SrCluster.Nodes[i].MstConfig.OAuthConfig.IdentityProviders[j].Name
					identityProvider.Challenge = config.SrCluster.Nodes[i].MstConfig.OAuthConfig.IdentityProviders[j].UseAsChallenger
					identityProvider.Login = config.SrCluster.Nodes[i].MstConfig.OAuthConfig.IdentityProviders[j].UseAsLogin
					identityProvider.MappingMethod = config.SrCluster.Nodes[i].MstConfig.OAuthConfig.IdentityProviders[j].MappingMethod

					if err := json.Unmarshal([]byte(jsonProvider), &identityProvider.Provider); err != nil {
						fmt.Printf("Unmarshal to ProviderInfo failed: %v", err)
					}

					config.SrCluster.Nodes[i].OAuthIDProviders = append(config.SrCluster.Nodes[i].OAuthIDProviders, identityProvider)
				}
			}
		}
	}
}

func (node NodeConfig) ParseMaster(config *Info) *V1MasterConfig {
	var mstconf V1MasterConfig

	jsonData, err := kyaml.ToJSON(node.Payload)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(jsonData, &mstconf); err != nil {
		log.Printf("unmarshal error %v", err)
	}
	return &mstconf
}

func (cluster *ClusterV3) ParseNode(config Info) int {
	for i := range cluster.Nodes {
		fmt.Printf("ParseNode: %v", cluster.Nodes[i].FileName)
		return 1
	}
	return 0
}

func (cluster ClusterV4) Show() string {
	som := make([]string, 100)

	for i := range cluster.CRManifests {
		som = append(som, fmt.Sprintf("qDst Cluster: %s\n", *cluster.CRManifests[i]))
	}

	return strings.Join(som, "")
}

func (cluster ClusterV3) Show() string {
	var payload = "not loaded"
	som := make([]string, 100)

	for i := range cluster.Nodes {
		if len(cluster.Nodes[i].Payload) > 0 {
			payload = "loaded"
		}
		som = append(som, fmt.Sprintf("Src Cluster: {Name:%s File: %s Payload: %s}\n",
			cluster.Nodes[i].Name, cluster.Nodes[i].Path+cluster.Nodes[i].FileName, payload))
	}

	return strings.Join(som, "")
}

func (info *Info) Show() string {
	return fmt.Sprintf("CPMA info:\n") +
		info.SrCluster.Show() +
		info.DsCluster.Show() +
		fmt.Sprintf("%#v\n", info.SFTP) +
		fmt.Sprintf("Output Path:%s\n", info.OutputPath)
}

// New returns a instance of the application settings.
func New() *Info {
	var info Info

	if err := config.Config().Unmarshal(&info); err != nil {
		log.Fatalf("unable to parse configuration: %v", err)
	}

	// TODO: make switch for config data or ask user for data
	list := make([][]string, 2)
	list[0] = []string{"master", "master-config.yaml", "/etc/origin/master"}
	list[1] = []string{"node", "node-config.yaml", "/etc/origin/node"}

	info.SrCluster = ClusterV3{}
	info.DsCluster = ClusterV4{}
	info.SrCluster.load(list)
	return &info
}

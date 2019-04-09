package env

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/fusor/cpma/internal/sftpclient"
	v1 "github.com/openshift/origin/pkg/cmd/server/apis/config/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"

	"github.com/spf13/viper"
)

type Cluster struct {
	Nodes []NodeConfig
}

type ClusterOld struct {
	Nodes map[string]*NodeConfig
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
	SrCluster  Cluster         `mapstructure:"srcluster"`
	DsCluster  Cluster         `mapstructure:"dscluster"`
	SFTP       sftpclient.Info `mapstructure:"Source"`
	OutputPath string          `mapstructure:"outputPath"`
}

type NodeConfig struct {
	Name         string
	FileName     string
	Path         string
	Payload      []byte
	MstConfig    *V1MasterConfig
	NdeConfig    *V1NodeConfig
	PlugProvider ProviderInfo
}

type V1MasterConfig v1.MasterConfig

type V1NodeConfig v1.NodeConfig

type ProviderInfo struct {
	APIVersion   string `json:"apiVersion"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	Kind         string `json:"kind"`
}

func addNode(name, filename, path string) *NodeConfig {
	x := NodeConfig{}
	x.Name = name
	x.FileName = filename
	x.Path = path
	return &x
}

func (config *Info) FetchSrc() int {
	sftpclient := config.SFTP.NewClient()
	defer sftpclient.Close()

	for i := range config.SrCluster.Nodes {
		srcFilePath := filepath.Join(config.SrCluster.Nodes[i].Path, config.SrCluster.Nodes[i].FileName)
		dstFilePath := filepath.Join(config.OutputPath, srcFilePath)

		sftpclient.GetFile(srcFilePath, dstFilePath)

		payload, err := ioutil.ReadFile(dstFilePath)
		if err != nil {
			log.Fatal(err)
		}
		config.SrCluster.Nodes[i].Payload = payload
	}

	return 1
}

func (cluster *Cluster) load(list [][]string) {
	for _, nc := range list {
		cluster.Nodes = append(cluster.Nodes, *addNode(nc[0], nc[1], nc[2]))
	}
}

func (config *Info) Parse() {
	for i := range config.SrCluster.Nodes {
		if config.SrCluster.Nodes[i].Name == "master" {
			config.SrCluster.Nodes[i].MstConfig = config.SrCluster.Nodes[i].ParseMaster(config)

			// WorkAround: until context for plugin extension
			for j := range config.SrCluster.Nodes[i].MstConfig.OAuthConfig.IdentityProviders {
				if fmt.Sprint(reflect.TypeOf(config.SrCluster.Nodes[i].MstConfig.OAuthConfig.IdentityProviders[j].Provider)) == "runtime.RawExtension" {
					foo := string(config.SrCluster.Nodes[i].MstConfig.OAuthConfig.IdentityProviders[j].Provider.Raw)

					// TODO: Replace PlugProvider with []IdentityProviders
					if err := json.Unmarshal([]byte(foo), &config.SrCluster.Nodes[i].PlugProvider); err != nil {
						fmt.Printf("unmarshal to ProviderInfo failed: %v", err)
					}
				}
			}
		} else if config.SrCluster.Nodes[i].Name == "node" {
			config.DsCluster.ParseNode(*config)
		}
	}
}

func (node NodeConfig) ParseMaster(config *Info) *V1MasterConfig {
	sftpclient := config.SFTP.NewClient()
	defer sftpclient.Close()
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

func (cluster *Cluster) ParseNode(config Info) int {
	for i := range cluster.Nodes {
		fmt.Printf("ParseNode: %v", cluster.Nodes[i].FileName)
		return 1
	}
	return 0
}

func (cluster Cluster) Show() string {
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
		fmt.Sprintf("Dst Cluster:{%s}\n", info.DsCluster) +
		fmt.Sprintf("%#v\n", info.SFTP) +
		fmt.Sprintf("Output Path:%s\n", info.OutputPath)
}

// New returns a instance of the application settings.
func New() *Info {
	var info Info

	if err := viper.Unmarshal(&info); err != nil {
		log.Fatalf("unable to parse configuration: %v", err)
	}

	// TODO: make switch for config data or ask user for data
	list := make([][]string, 2)
	list[0] = []string{"master", "master-config.yaml", "/etc/origin/master"}
	list[1] = []string{"node", "node-config.yaml", "/etc/origin/node"}

	info.SrCluster = Cluster{}
	info.SrCluster.load(list)
	return &info
}

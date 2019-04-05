package env

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/fusor/cpma/internal/sftpclient"
	v1 "github.com/openshift/origin/pkg/cmd/server/apis/config/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"

	"github.com/spf13/viper"
)

type Cluster struct {
	Nodes map[string]NodeConfig
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
	FileName  string
	Path      string
	Payload   []byte
	MstConfig V1MasterConfig
	NdeConfig V1NodeConfig
}

type V1MasterConfig v1.MasterConfig

type V1NodeConfig v1.NodeConfig

func addNode(filename, path string) *NodeConfig {
	x := NodeConfig{}
	x.FileName = filename
	x.Path = path
	return &x
}

func (config *Info) FetchSrc() int {
	sftpclient := config.SFTP.NewClient()
	defer sftpclient.Close()

	for key, nodeconfig := range config.SrCluster.Nodes {
		srcFilePath := nodeconfig.Path + "/" + nodeconfig.FileName
		dstFilePath := filepath.Join(config.OutputPath, srcFilePath)
		sftpclient.GetFile(srcFilePath, dstFilePath)

		payload, err := ioutil.ReadFile(dstFilePath)
		if err != nil {
			log.Fatal(err)
		}

		src := config.SrCluster.Nodes[key]
		src.Payload = payload
		config.SrCluster.Nodes[key] = src
		return 1
	}
	return 0
}

func (cluster *Cluster) load(list [][]string) {
	for _, nc := range list {
		cluster.Nodes[nc[0]] = *addNode(nc[1], nc[2])
	}
}

func (config *Info) Parse() {
	for key, _ := range config.SrCluster.Nodes {
		if key == "master" {
			src := config.SrCluster.Nodes[key]
			src.ParseMaster(*config)
		} else if key == "node" {
			config.DsCluster.ParseNode(*config)
		}
	}
}

func (node *NodeConfig) ParseMaster(config Info) {
	sftpclient := config.SFTP.NewClient()
	defer sftpclient.Close()

	jsonData, err := kyaml.ToJSON(node.Payload)
	if err != nil {
		log.Fatal(err)
	}

	error := json.Unmarshal(jsonData, &node.MstConfig)
	if error != nil {
		fmt.Printf("err was %v", err)
	}

	fmt.Printf("%+v\n", node.MstConfig.ServingInfo.BindAddress)
	fmt.Printf("%+v\n", node.MstConfig.OAuthConfig)
	return node
}

func (cluster *Cluster) ParseNode(config Info) int {
	for _, nodeconfig := range cluster.Nodes {
		fmt.Println(fmt.Sprintf("%v", nodeconfig))
		return 1
	}
	return 0
}

func (srcluster Cluster) Show() string {
	var payload = "not loaded"
	som := make([]string, 100)

	for name, nodeconfig := range srcluster.Nodes {
		if len(nodeconfig.Payload) > 0 {
			payload = "loaded"
		}
		som = append(som, fmt.Sprintf("Src Cluster:{Name:%s File: %s Payload: %s}\n", name, nodeconfig.Path+nodeconfig.FileName, payload))
	}
	return strings.Join(som, "")
}

func (info *Info) Show() string {
	return fmt.Sprintf("CPAM info:\n") +
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
	info.SrCluster.Nodes = make(map[string]NodeConfig)
	info.SrCluster.load(list)
	return &info
}

func (c *V1MasterConfig) UnmarshMaster(data []byte) error {
	return json.Unmarshal(data, c)
}

func (c *V1NodeConfig) UnmarshNode(data []byte) error {
	return json.Unmarshal(data, c)
}

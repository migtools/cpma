package env

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/fusor/cpma/internal/sftpclient"
	"github.com/spf13/viper"
)

type ClusterCMD interface {
	addNode(name, file, path string) int
	Show() string
}
type Cluster struct {
	Nodes map[string]NodeConfig
}

type NodeConfig struct {
	FileName string
	Path     string
	Payload  string
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
		src.Payload = string(payload)
		config.SrCluster.Nodes[key] = src
	}
	return 0
}

func addNode(filename, path string) *NodeConfig {
	x := NodeConfig{}
	x.FileName = filename
	x.Path = path
	return &x
}

func (cluster *Cluster) load(list [][]string) {
	for _, nc := range list {
		cluster.Nodes[nc[0]] = *addNode(nc[1], nc[2])
	}
}

func (srcluster Cluster) Show() string {
	var payload = ""
	som := make([]string, 100)

	for name, nodeconfig := range srcluster.Nodes {
		if nodeconfig.Payload != "" {
			payload = "loaded"
		}
		som = append(som, fmt.Sprintf("Src Cluster:{Name:%s File: %s Payload: %s}\n", name, nodeconfig.Path+nodeconfig.FileName, payload))
	}
	return strings.Join(som, "")
}

func (info *Info) Show() string {
	return fmt.Sprintf("\nCPAM info:\n") +
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
	info.FetchSrc()
	return &info
}

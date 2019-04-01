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

type Clusters map[string]NodeConfig

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
	SrCluster  Clusters        `mapstructure:"srcluster"`
	DsCluster  Clusters        `mapstructure:"dscluster"`
	SFTP       sftpclient.Info `mapstructure:"Source"`
	OutputPath string          `mapstructure:"outputPath"`
}

func (config *Info) FetchSrc() int {
	sftpclient := config.SFTP.NewClient()
	defer sftpclient.Close()

	for key, nodeconfig := range config.SrCluster {
		srcFilePath := nodeconfig.Path + "/" + nodeconfig.FileName
		dstFilePath := filepath.Join(config.OutputPath, "data"+srcFilePath)
		sftpclient.GetFile(srcFilePath, dstFilePath)

		payload, err := ioutil.ReadFile(dstFilePath)
		if err != nil {
			log.Fatal(err)
		}

		src := config.SrCluster[key]
		src.Payload = string(payload)
		config.SrCluster[key] = src
	}
	return 0
}

func (cluster Clusters) addNode(name, file, path string) int {
	x := NodeConfig{}
	x.setConfig(file, path)
	cluster[name] = x
	return 0
}

func (cluster *Clusters) load(list [][]string) {
	for _, nc := range list {
		cluster.addNode(nc[0], nc[1], nc[2])
	}
}

func (config *NodeConfig) setConfig(filename, path string) {
	config.FileName = filename
	config.Path = path
}

func (srcluster *Clusters) Show() string {
	var payload = ""
	som := make([]string, 100)

	for name, nodeconfig := range *srcluster {
		if nodeconfig.Payload != "" {
			payload = "loaded"
		}
		som = append(som, fmt.Sprintf("info.SrcCluster:(Name:%s File: %s Payload: %s)\n", name, nodeconfig.Path+nodeconfig.FileName, payload))
	}
	return strings.Join(som, "")
}

func (info *Info) Show() string {
	return fmt.Sprintf("\nCPAM info:\n") +
		info.SrCluster.Show() +
		fmt.Sprintf("%#v\n", info.DsCluster) +
		fmt.Sprintf("%#v\n", info.SFTP) +
		fmt.Sprintf("%#v\n", info.OutputPath) +
		fmt.Sprintf("\n")
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

	info.SrCluster = make(Clusters)
	info.SrCluster.load(list)

	return &info
}

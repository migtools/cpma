package env

import (
	"log"
	"path/filepath"

	"github.com/fusor/cpma/internal/sftpclient"
	"github.com/spf13/viper"
)

type ClusterCMD interface {
	addNode(name, file, path string) int
}

type Clusters map[string]NodeConfig

type NodeConfig struct {
	FileName string
	Path     string
}

type Cmd interface {
	setSFTP() int
	Fetch() int        //SFTP       sftpclient.Info
	getConfig() string // -> OutputPath
}

type Info struct {
	Cluster    Clusters        `mapstructure:"cluster"`
	SFTP       sftpclient.Info `mapstructure:"Source"`
	OutputPath string          `mapstructure:"outputPath"`
}

func (config *Info) Fetch() int {
	sftpclient := config.SFTP.NewClient()
	defer sftpclient.Close()

	for _, cluster := range config.Cluster {
		srcFilePath := cluster.Path + "/" + cluster.FileName
		dstFilePath := "data" + srcFilePath
		sftpclient.GetFile(srcFilePath, filepath.Join(config.OutputPath, dstFilePath))
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

	info.Cluster = make(Clusters)
	info.Cluster.load(list)

	return &info
}

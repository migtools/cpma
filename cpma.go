package main

import (
	"fmt"

	"github.com/fusor/cpma/cmd"
	"github.com/fusor/cpma/env"
	_ "github.com/fusor/cpma/internal/log"
	log "github.com/sirupsen/logrus"
)

func main() {
	cmd.Execute()

	config := env.New()

	config.FetchSrc()
	config.Parse()

	for i := range config.SrCluster.Nodes {
		fmt.Println(fmt.Sprintf("%s", config.SrCluster.Nodes[i].FileName))
		if config.SrCluster.Nodes[i].MstConfig != nil {
			log.Println(config.SrCluster.Nodes[i].MstConfig.ServingInfo.BindAddress)
			log.Println(config.SrCluster.Nodes[i].MstConfig.OAuthConfig.MasterPublicURL)
			log.Println(config.SrCluster.Nodes[i].PlugProvider.ClientSecret)
		}
	}
	log.Print(config.Show())
}

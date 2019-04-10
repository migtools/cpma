package main

import (
	"fmt"

	"github.com/fusor/cpma/cmd"
	"github.com/fusor/cpma/env"
	_ "github.com/fusor/cpma/pkg/log"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Println("CPMA started")

	cmd.Execute()

	config := env.New()

	config.FetchSrc()
	config.Parse()

	for i := range config.SrCluster.Nodes {
		fmt.Println(fmt.Sprintf("%s", config.SrCluster.Nodes[i].FileName))
		if config.SrCluster.Nodes[i].MstConfig != nil {
			fmt.Println("===>")
			fmt.Println(fmt.Sprintf("%s", config.SrCluster.Nodes[i].MstConfig.ServingInfo.BindAddress))
			fmt.Println(fmt.Sprintf("%s", config.SrCluster.Nodes[i].MstConfig.OAuthConfig.MasterPublicURL))
			fmt.Println(fmt.Sprintf("%s", config.SrCluster.Nodes[i].PlugProvider.ClientSecret))
		}
	}
	log.Print(config.Show())
	log.Println("CPMA finished")
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fusor/cpma/cmd"
	"github.com/fusor/cpma/env"
)

func main() {
	log.Print("CPMA started")
	f, err := os.OpenFile("cpma.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("CPMA Log started")
	cmd.Execute()

	config := env.New()
	config.FetchSrc()
	config.Parse()

	fmt.Printf("%+v\n", config.SrCluster.Nodes["master"])

	fmt.Printf("%+v\n", config.SrCluster.Nodes["master"].MstConfig.ServingInfo.BindAddress)
	fmt.Printf("%+v\n", config.SrCluster.Nodes["master"].MstConfig.OAuthConfig)

	log.Print(config.Show())
	log.Print("CPMA finished")
}

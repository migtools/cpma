package main

import (
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

	//x := config.SrCluster.Nodes["master"].MstConfig.OAuthConfig.AssetPublicURL
	//z := fmt.Sprintf("%s", x)
	//fmt.Println(z)

	log.Print(config.Show())
}

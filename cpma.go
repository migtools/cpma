package main

import (
	"log"

	"github.com/fusor/cpma/env"
	"github.com/fusor/cpma/sftpclient"
)

var (
	configFile = "config.yaml"
)

func main() {
	cfg, err := env.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("unable to read configuration file: %v", err)
	}

	srcFilePath := "./essay/config.yaml"
	dstFilePath := "./data/config.yaml"
	sftpclient.GetFile(
		cfg.Source.HostName,
		cfg.Source.UserName,
		cfg.Source.SSHKey,
		srcFilePath, dstFilePath,
	)
}

package main

import (
	"github.com/fusor/cpma/cmd"
	"github.com/fusor/cpma/env"
	_ "github.com/fusor/cpma/internal/log"
	log "github.com/sirupsen/logrus"
)

func main() {
	cmd.Execute()

	config := env.New()

	if config.LocalOnly {
		config.LoadSrc()
	} else {
		config.FetchSrc()
	}

	config.Parse()

	log.Print(config.Show())
}

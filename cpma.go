package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/fusor/cpma/cmd"
	"github.com/fusor/cpma/env"
	_ "github.com/fusor/cpma/internal/log"
	ocp3 "github.com/fusor/cpma/ocp3config"
)

func main() {
	cmd.Execute()
	e := env.New()

	// TODO: Passing *e.Info here is not exactly nice. Fix?
	ocp3config := ocp3.New()
	ocp3config.Fetch(e)

	m := ocp3config.ParseMaster()

	log.Println(m)
}

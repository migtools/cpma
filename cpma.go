package main

import (
	"github.com/fusor/cpma/cmd"
	"github.com/fusor/cpma/ocp3"
	"github.com/fusor/cpma/ocp4"
)

func main() {
	cmd.Execute()

	ocp3config := ocp3.New()
	ocp3config.Fetch()

	mc := ocp3config.ParseMaster()
	clusterV4 := ocp4.Cluster{}
	clusterV4.Translate(mc)
	clusterV4.GenYAML()
}

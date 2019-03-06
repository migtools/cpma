package main

import (
	"github.com/fusor/cpma/cmd"
	"github.com/fusor/cpma/sftpclient"
	"github.com/spf13/viper"
)

func main() {
	cmd.Execute()

	srcFilePath := "./essay/config.yaml"
	dstFilePath := "./data/config.yaml"

	sftpclient.GetFile(
		viper.GetString("source.hostname"),
		viper.GetString("source.username"),
		viper.GetString("source.sshkey"),
		srcFilePath, dstFilePath,
	)
}

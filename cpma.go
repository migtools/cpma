package main

import (
	"github.com/fusor/cpma/sftpclient"
)

func main() {
	//config := setup.LoadConfiguration(configFileName)

	cluster := "mars8.ddns.net"
	user := "qwerty66"
	keyfile := "/home/gildub/.ssh/test_user"
	srcFilePath := "./essay/config.yaml"
	dstFilePath := "./data/config.yaml"

	sftpclient.GetFile(cluster, user, keyfile, srcFilePath, dstFilePath)
}

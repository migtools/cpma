package main

import (
	"path/filepath"

	"github.com/fusor/cpma/cmd"
	"github.com/fusor/cpma/env"
)

func main() {
	cmd.Execute()
	config := env.New()

	sftpclient := config.SFTP.NewClient()
	defer sftpclient.Close()

	srcFilePath := "/etc/origin/master/master-config.yaml"
	dstFilePath := "./master-config.yaml"
	sftpclient.GetFile(srcFilePath, filepath.Join(config.OutputPath, dstFilePath))

	srcFilePath = "/etc/origin/node/node-config.yaml"
	dstFilePath = "./node-config.yaml"
	sftpclient.GetFile(srcFilePath, filepath.Join(config.OutputPath, dstFilePath))

}

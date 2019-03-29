package main

import (
	"fmt"
	"path/filepath"

	"github.com/fusor/cpma/cmd"
	"github.com/fusor/cpma/env"
)

func main() {
	cmd.Execute()
	config := env.New()

	sftpclient := config.SFTP.NewClient()
	defer sftpclient.Close()

	fmt.Println(config)
	fmt.Println()

	for _, cluster := range config.Cluster {
		srcFilePath := cluster.Path + "/" + cluster.FileName
		dstFilePath := "data" + srcFilePath
		sftpclient.GetFile(srcFilePath, filepath.Join(config.OutputPath, dstFilePath))
	}
}

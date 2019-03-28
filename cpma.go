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

	configfiles := map[string]map[string]string{
		"master": map[string]string{
			"name": "master-config.yaml",
			"path": "/etc/origin/master",
		},
		"node": map[string]string{
			"name": "node-config.yaml",
			"path": "/etc/origin/node",
		},
	}

	for _, data := range configfiles {
		srcFilePath := data["path"] + "/" + data["name"]
		dstFilePath := "data" + srcFilePath
		sftpclient.GetFile(srcFilePath, filepath.Join(config.OutputPath, dstFilePath))
	}
}

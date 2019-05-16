package io

import (
	"io/ioutil"
	"os"

	"github.com/fusor/cpma/pkg/io/sftpclient"
)

// GetFile allows to mock file retrieval
var GetFile = fetchFile

// Fetch first tries to retrieve file from local disk (outputDir/<Hostname>/).
// If it fails then connects to Hostname to retrieve file and stores it locally
// To force a network connection remove outputDir/... prior to exec.
func fetchFile(host, src, target string) ([]byte, error) {
	if fileExists(target) {
		return ioutil.ReadFile(target)
	} else {
		sftpclient.Fetch(host, src, target)
		netFile, err := ioutil.ReadFile(target)
		if err != nil {
			return nil, err
		}
		return netFile, nil
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

package io

import (
	"io/ioutil"

	"github.com/fusor/cpma/internal/io/sftpclient"
	"github.com/sirupsen/logrus"
)

var GetFile = fetchFile

// Fetch first tries to retrieve file from local disk (outputDir/<Hostname>/).
// If it fails then connects to Hostname to retrieve file and stores it locally
// To force a network connection remove outputDir/... prior to exec.
func fetchFile(host, src, target string) []byte {
	f, err := ioutil.ReadFile(target)
	if err != nil {
		sftpclient.Fetch(host, src, target)
		netFile, err := ioutil.ReadFile(target)
		if err != nil {
			logrus.Fatal(err)
		}
		return netFile
	} else {
		return f
	}
}

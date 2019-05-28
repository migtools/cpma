package io

import (
	"io/ioutil"

	"github.com/fusor/cpma/pkg/io/sftpclient"
	"github.com/sirupsen/logrus"
)

// GetFile first tries to retrieve file from local disk (outputDir/<Hostname>/).
// If it fails then connects to Hostname to retrieve file and stores it locally
// To force a network connection remove outputDir/... prior to exec.
var GetFile = func(host, src, target string) ([]byte, error) {
	f, err := ioutil.ReadFile(target)
	if err != nil {
		sftpclient.Fetch(host, src, target)
		netFile, err := ioutil.ReadFile(target)
		if err != nil {
			return nil, err
		}
		return netFile, nil
	}
	return f, nil
}

// FetchEnv Fetch env vars from the OCP3 cluster
func FetchEnv(host, envVar string) (string, error) {
	output, err := sftpclient.GetEnvVar(host, envVar)
	if err != nil {
		return "", err
	}
	logrus.Debugf("Env:loaded: %s", envVar)

	return output, nil
}

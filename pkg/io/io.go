package io

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io/remotehost"
	"github.com/sirupsen/logrus"

	configv1 "github.com/openshift/api/legacyconfig/v1"
)

// FetchFile first tries to retrieve file from local disk (outputDir/<Hostname>/).
// If it fails then connects to Hostname to retrieve file and stores it locally
// To force a network connection remove outputDir/... prior to exec.
var FetchFile = func(src string) ([]byte, error) {
	dst := filepath.Join(env.Config().GetString("OutputDir"), env.Config().GetString("Source"), src)
	f, err := ioutil.ReadFile(dst)
	if err != nil {
		host := env.Config().GetString("Source")
		remotehost.Fetch(host, src, dst)
		netFile, err := ioutil.ReadFile(src)
		if err != nil {
			return nil, err
		}
		return netFile, nil
	}
	return f, nil
}

// FetchEnv Fetch env vars from the OCP3 cluster
func FetchEnv(host, envVar string) (string, error) {
	output, err := remotehost.GetEnvVar(host, envVar)
	if err != nil {
		return "", err
	}
	logrus.Debugf("Env:loaded: %s", envVar)

	return output, nil
}

// FetchStringSource fetches a string from an OCP3 cluster
func FetchStringSource(stringSource configv1.StringSource) (string, error) {
	if stringSource.Value != "" {
		return stringSource.Value, nil
	}

	if stringSource.File != "" {
		fileContent, err := FetchFile(stringSource.File)
		if err != nil {
			return "", nil
		}

		fileString := strings.TrimSuffix(string(fileContent), "\n")
		return fileString, nil
	}

	if stringSource.Env != "" {
		env, err := FetchEnv(env.Config().GetString("Source"), stringSource.Env)
		if err != nil {
			return "", nil
		}

		return env, nil
	}

	return "", nil
}

// ReadFile reads a file and returns its contents
func ReadFile(file string) ([]byte, error) {
	src := filepath.Join(env.Config().GetString("OutputDir"), file)
	return ioutil.ReadFile(src)
}

// WriteFile writes data to a file
func WriteFile(content []byte, file string) error {
	dst := filepath.Join(env.Config().GetString("OutputDir"), file)
	os.MkdirAll(path.Dir(dst), 0750)
	return ioutil.WriteFile(dst, content, 0640)
}

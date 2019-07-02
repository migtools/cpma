package io

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io/remotehost"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// FetchFile first tries to retrieve file from local disk (outputDir/<Hostname>/).
// If it fails then connects to Hostname to retrieve file and stores it locally
// To force a network connection remove outputDir/... prior to exec.
var FetchFile = func(src string) ([]byte, error) {
	dst := filepath.Join(env.Config().GetString("Hostname"), src)
	f, err := ReadFile(dst)
	if err != nil {
		host := env.Config().GetString("Hostname")

		cmd := fmt.Sprintf("sudo cat %s", src)
		output, err := remotehost.RunCMD(host, cmd)
		if err != nil {
			return nil, err
		}

		if output == "" {
			msg := fmt.Sprintf("Empty or missing file: %s", dst)
			return nil, errors.New(msg)
		}

		err = WriteFile([]byte(output), dst)
		if err != nil {
			logrus.Errorf("Unable to save: %s", dst)
			return nil, err
		}

		netFile, err := ReadFile(dst)
		if err != nil {
			return nil, err
		}
		return netFile, nil

	}
	return f, nil
}

// FetchEnv Fetch env vars from the OCP3 cluster
func FetchEnv(host, envVar string) (string, error) {
	cmd := fmt.Sprintf("print $%s", envVar)
	output, err := remotehost.RunCMD(host, cmd)
	if err != nil {
		return "", errors.Wrap(err, "Can't fetch env variable")
	}
	logrus.Debugf("Env:loaded: %s", envVar)

	return output, nil
}

// FetchStringSource fetches a string from an OCP3 cluster
func FetchStringSource(stringSource legacyconfigv1.StringSource) (string, error) {
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
		env, err := FetchEnv(env.Config().GetString("Hostname"), stringSource.Env)
		if err != nil {
			return "", nil
		}

		return env, nil
	}

	return "", nil
}

// ReadFile reads a file in OutputDir and returns its contents
func ReadFile(file string) ([]byte, error) {
	src := filepath.Join(env.Config().GetString("OutputDir"), file)
	return ioutil.ReadFile(src)
}

// WriteFile writes data to a file in OutputDir
func WriteFile(content []byte, file string) error {
	dst := filepath.Join(env.Config().GetString("OutputDir"), file)
	os.MkdirAll(path.Dir(dst), 0750)
	return ioutil.WriteFile(dst, content, 0640)
}

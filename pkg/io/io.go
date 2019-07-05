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

// FetchFile first tries to retrieve file from local disk (workDir/<Hostname>/).
// If it fails then connects to Hostname to retrieve file and stores it locally
// To force a network connection remove workDir/... prior to exec.
var FetchFile = func(src string) ([]byte, error) {
	var f []byte
	var err error

	if env.Config().GetBool("FetchFromRemote") {
		f, err = fetchFromRemote(src)
	} else {
		f, err = fetchFromLocal(src)
	}

	return f, err
}

func fetchFromRemote(src string) ([]byte, error) {
	host := env.Config().GetString("Hostname")
	dst := filepath.Join(host, src)

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

func fetchFromLocal(src string) ([]byte, error) {
	localSrc := filepath.Join(env.Config().GetString("WorkDir"), env.Config().GetString("Source"), src)
	f, err := ioutil.ReadFile(localSrc)
	if err != nil {
		return nil, err
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

// ReadFile reads a file in WorkDir and returns its contents
func ReadFile(file string) ([]byte, error) {
	src := filepath.Join(env.Config().GetString("WorkDir"), file)
	return ioutil.ReadFile(src)
}

// WriteFile writes data to a file in WorkDir
func WriteFile(content []byte, file string) error {
	dst := filepath.Join(env.Config().GetString("WorkDir"), file)
	os.MkdirAll(path.Dir(dst), 0750)
	return ioutil.WriteFile(dst, content, 0640)
}

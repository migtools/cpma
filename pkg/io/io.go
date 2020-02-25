package io

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/konveyor/cpma/pkg/env"
	"github.com/konveyor/cpma/pkg/io/remotehost"
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
		f, err = FetchFromLocal(src)
	}

	return f, err
}

func fetchFromRemote(src string) ([]byte, error) {
	host := env.Config().GetString("Hostname")
	dst := filepath.Join(host, src)

	logrus.Debugf("Fetching Remote File %s:%s", host, src)
	cmd := fmt.Sprintf("sudo sh -c 'if [[ ! -f %s ]]; then echo not-found; fi'", src)
	output0, err := remotehost.RunCMD(host, cmd)
	if err != nil {
		return nil, errors.Wrapf(err, "Error accessing file %s", dst)
	}

	if output0 == "not-found" {
		msg := fmt.Sprintf("File %s not found", dst)
		return nil, errors.New(msg)
	}

	cmd = fmt.Sprintf("sudo cat %s", src)
	output, err := remotehost.RunCMD(host, cmd)
	if err != nil {
		return nil, errors.Wrapf(err, "Error accessing file %s", dst)
	}

	if output == "" {
		msg := fmt.Sprintf("Empty file: %s", dst)
		return nil, errors.New(msg)
	}

	if err := WriteFile([]byte(output), dst); err != nil {
		logrus.Errorf("Unable to save file: %s", dst)
		return nil, err
	}

	netFile, err := ReadFile(dst)
	if err != nil {
		return nil, err
	}
	return netFile, nil
}

// FetchFromLocal retrieve file from local WorkDir
func FetchFromLocal(src string) ([]byte, error) {
	localSrc := filepath.Join(env.Config().GetString("WorkDir"), env.Config().GetString("Hostname"), src)
	logrus.Debugf("Fetching Local File %s", localSrc)
	if fileExists(localSrc) {
		f, err := ioutil.ReadFile(localSrc)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	msg := fmt.Sprintf("No such file %s", localSrc)
	return nil, errors.New(msg)
}

// FetchEnv Fetch env vars from either the source cluster or localhost
func FetchEnv(host, envVar string) (string, error) {
	var output string

	if env.Config().GetBool("FetchFromRemote") {
		var err error
		cmd := fmt.Sprintf("printf $%s", envVar)
		output, err = remotehost.RunCMD(host, cmd)
		if err != nil {
			return "", errors.Wrap(err, "Can't fetch env variable")
		}
	} else {
		output = os.Getenv(envVar)
	}

	logrus.Debugf("Env:loaded: %s", envVar)

	return output, nil
}

// FetchStringSource fetches a string from an source cluster
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

// ReadFile reads a file from WorkDir and returns its contents
func ReadFile(file string) ([]byte, error) {
	src := filepath.Join(env.Config().GetString("WorkDir"), file)
	return ioutil.ReadFile(src)
}

// WriteFile writes data to a file into WorkDir
func WriteFile(content []byte, file string) error {
	dst := filepath.Join(env.Config().GetString("WorkDir"), file)
	os.MkdirAll(path.Dir(dst), 0750)
	return ioutil.WriteFile(dst, content, 0640)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

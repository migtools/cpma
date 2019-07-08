package e2e

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/fusor/cpma/pkg/transform"

	"github.com/fusor/cpma/pkg/env"
	"github.com/pkg/errors"
)

// OpenClusterSession will ensure that cluster session is open
func OpenClusterSession(tmpDir string) error {
	cmd := exec.Command("which", "oc")
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "Cant locate oc binary ")
	}

	clusterAddr := os.Getenv("CPMA_HOSTNAME")
	login := os.Getenv("CPMA_LOGIN")
	passwd := os.Getenv("CPMA_PASSWD")
	kubeconfig, exists := os.LookupEnv("KUBECONFIG")
	if !exists {
		kubeconfig = path.Join(tmpDir, "kubeconfig")
		os.Setenv("KUBECONFIG", kubeconfig)
	}

	binary := "oc"
	commandArgs := []string{
		"login", clusterAddr,
		"-u", login,
		"-p", passwd,
		"--insecure-skip-tls-verify",
		"--config", kubeconfig}
	cmd = exec.Command(binary, commandArgs...)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "Cant open cluster session")
	}
	return nil
}

// RunCpma build and execute the tool
// on provided set of environment variables
func RunCpma() error {
	cmd := exec.Command("make", "build")
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "Couldn't build a binary")
	}

	err = env.InitConfig()
	if err != nil {
		return errors.Wrap(err, "Can't initialize config")
	}
	binary := path.Join("bin", "cpma")
	cmd = exec.Command(binary) //, commandArgs...)
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "Can't execute the binary")
	}
	return nil
}

// ReadReport reads and unmarshal the report into report struceture from transform
func ReadReport(pathToReport string) (report *transform.ReportOutput, err error) {
	srcReport, err := ioutil.ReadFile(pathToReport)
	if err != nil {
		return nil, errors.Wrap(err, "Error while reading report")
	}
	err = json.Unmarshal(srcReport, &report)
	if err != nil {
		return nil, errors.Wrap(err, "Can't unmarshal report to report structure.")
	}
	return
}

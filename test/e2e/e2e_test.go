package e2e

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"reflect"
	"testing"

	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/transform"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestReport(t *testing.T) {
	var (
		e2eTestDataDir string
		e2eTestOut     string
		e2eTestSrc     string

		err error
	)

	e2eTestDataDir = path.Join("test", "e2e", "testdata")
	e2eTestOut = path.Join(e2eTestDataDir, "out")
	e2eTestSrc = path.Join(e2eTestDataDir, "src")

	err = openClusterSession(e2eTestOut)
	assert.NoError(t, err, "Could not open cluster session")

	os.Chdir("../..")
	os.Setenv("CPMA_WORKDIR", e2eTestOut)
	os.Setenv("CPMA_CREATECONFIG", "no")
	os.Setenv("CPMA_CONFIGSOURCE", "remote")
	os.Setenv("CPMA_INSECUREHOSTKEY", "true")

	err = runCpma()
	assert.NoError(t, err, "Couldn't execute CPMA")

	sourceReport := path.Join(e2eTestSrc, "report.json")
	targetReport := path.Join(e2eTestOut, "report.json")

	srcReport, err := readReport(sourceReport)
	assert.NoError(t, err, "Couldn't process source report")
	outReport, err := readReport(targetReport)
	assert.NoError(t, err, "Couldn't process target report")

	assert.True(t, reflect.DeepEqual(&srcReport, &outReport), "Reports are not equal")

	os.RemoveAll(e2eTestOut)
}

func TestManifestsReports(t *testing.T) {
	var (
		e2eTestDataDir string
		e2eTestOut     string

		err error
	)

	e2eTestDataDir = path.Join("test", "e2e", "testdata")
	e2eTestOut = path.Join(e2eTestDataDir, "out")

	err = openClusterSession(e2eTestOut)
	assert.NoError(t, err, "Could not open cluster session")

	os.Chdir("../..")
	os.Setenv("CPMA_WORKDIR", e2eTestOut)
	os.Setenv("CPMA_CREATECONFIG", "no")
	os.Setenv("CPMA_CONFIGSOURCE", "remote")
	os.Setenv("CPMA_INSECUREHOSTKEY", "true")

	err = runCpma()
	assert.NoError(t, err, "Couldn't execute CPMA")

	testCases := []struct {
		name      string
		reports   string
		manifests string
	}{
		{
			name:      "Only reports mode",
			manifests: "false",
			reports:   "true",
		},
		{
			name:      "Only manifests mode",
			manifests: "true",
			reports:   "false",
		},
		{
			name:      "Both allowed",
			manifests: "true",
			reports:   "true",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("CPMA_MANIFESTS", tc.manifests)
			os.Setenv("CPMA_REPORTS", tc.reports)
			err = runCpma()
			assert.NoError(t, err, "Couldn't execute CPMA")

			if env.Config().GetString("Manifests") == "true" && env.Config().GetString("Reports") == "false" {
				_, err := os.Stat(path.Join(e2eTestOut, "report.json"))
				os.IsNotExist(err)
				assert.Equal(t, true, os.IsNotExist(err))
			}

			if env.Config().GetString("Reports") == "true" && env.Config().GetString("Manifests") == "false" {
				_, err := os.Stat(path.Join(e2eTestOut, "manifests"))
				os.IsNotExist(err)
				assert.Equal(t, true, os.IsNotExist(err))
			}

			if env.Config().GetString("Manifests") == "true" && env.Config().GetString("Manifests") == "true" {
				_, err := os.Stat(path.Join(e2eTestOut, "report.json"))
				assert.Equal(t, nil, err)
				_, err = os.Stat(path.Join(e2eTestOut, "manifests"))
				assert.Equal(t, nil, err)
			}
		})
	}
	os.Unsetenv("CPMA_MANIFESTS")
	os.Unsetenv("CPMA_REPORTS")
}

// openClusterSession will ensure that cluster session is open
func openClusterSession(tmpDir string) error {
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

// runCpma build and execute the tool
// on provided set of environment variables
func runCpma() error {
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

// readReport reads and unmarshal the report into report struceture from transform
func readReport(pathToReport string) (report *transform.ReportOutput, err error) {
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
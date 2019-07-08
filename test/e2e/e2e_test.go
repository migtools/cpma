package e2e

import (
	"os"
	"path"
	"reflect"
	"testing"

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

	err = OpenClusterSession(e2eTestOut)
	assert.NoError(t, err, "Could not open cluster session")

	os.Chdir("../..")
	os.Setenv("CPMA_OUTPUTDIR", e2eTestOut)
	os.Setenv("CPMA_CREATECONFIG", "no")
	os.Setenv("CPMA_CONFIGSOURCE", "remote")
	os.Setenv("CPMA_INSECUREHOSTKEY", "true")

	err = RunCpma()
	assert.NoError(t, err, "Couldn't execute CPMA")

	sourceReport := path.Join(e2eTestSrc, "report.json")
	targetReport := path.Join(e2eTestOut, "report.json")

	srcReport, err := ReadReport(sourceReport)
	assert.NoError(t, err, "Couldn't process source report")
	outReport, err := ReadReport(targetReport)
	assert.NoError(t, err, "Couldn't process target report")

	assert.True(t, reflect.DeepEqual(&srcReport, &outReport), "Reports are not equal")

	os.RemoveAll(e2eTestOut)
}

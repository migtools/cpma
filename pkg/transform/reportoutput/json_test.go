package reportoutput

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/konveyor/cpma/pkg/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONOutput(t *testing.T) {
	reportJSON, err := ioutil.ReadFile("testdata/reportexample.json")
	require.NoError(t, err)

	report := &ReportOutput{}

	err = json.Unmarshal(reportJSON, report)
	require.NoError(t, err)

	jsonFileName = "reportactual.json"
	env.Config().Set("WorkDir", "testdata")
	jsonOutput(*report)

	expectedJSON, err := ioutil.ReadFile("testdata/reportexample.json")
	require.NoError(t, err)

	actualJSON, err := ioutil.ReadFile("testdata/reportactual.json")
	require.NoError(t, err)

	assert.Equal(t, expectedJSON, actualJSON)

	err = os.Remove("testdata/reportactual.json")
	require.NoError(t, err)
}

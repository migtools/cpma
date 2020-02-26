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

func TestHTMLOutput(t *testing.T) {
	reportJSON, err := ioutil.ReadFile("testdata/reportexample.json")
	require.NoError(t, err)

	report := &ReportOutput{}

	err = json.Unmarshal(reportJSON, report)
	require.NoError(t, err)

	htmlFileName = "reportactual.html"
	env.Config().Set("WorkDir", "testdata")
	htmlOutput(*report)

	expectedHTML, err := ioutil.ReadFile("testdata/reportexpected.html")
	require.NoError(t, err)

	actualHTML, err := ioutil.ReadFile("testdata/reportactual.html")
	require.NoError(t, err)

	assert.Equal(t, expectedHTML, actualHTML)

	err = os.Remove("testdata/reportactual.html")
	require.NoError(t, err)
}

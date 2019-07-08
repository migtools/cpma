package transform_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/api"
	"github.com/fusor/cpma/pkg/transform"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClusterReportExtractionTransform(t *testing.T) {
	apiResources := api.Resources{
		PersistentVolumeList: cpmatest.CreateTestPVList(),
		NodeList:             cpmatest.CreateTestNodeList(),
		StorageClassList:     cpmatest.CreateStorageClassList(),
		NamespaceList:        cpmatest.CreateTestNameSpaceList(),
	}

	clusterReportExtraction := transform.ClusterReportExtraction{apiResources}

	actualClusterReport, err := clusterReportExtraction.Transform()
	require.NoError(t, err)

	report := transform.ReportOutput{
		ClusterReport: actualClusterReport[0].(transform.ReportOutput).ClusterReport,
	}

	actualClusterReportJSON, err := json.MarshalIndent(report, "", " ")
	require.NoError(t, err)

	expectedClusterReportJSON, err := ioutil.ReadFile("testdata/expected-report-cluster.json")
	require.NoError(t, err)

	assert.Equal(t, expectedClusterReportJSON, actualClusterReportJSON)
}

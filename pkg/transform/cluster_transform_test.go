package transform_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/api"
	"github.com/fusor/cpma/pkg/transform"
	cpmatest "github.com/fusor/cpma/pkg/transform/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClusterExtractionTransform(t *testing.T) {
	apiResources := api.Resources{
		QuotaList:            cpmatest.CreateTestQuotaList(),
		PersistentVolumeList: cpmatest.CreateTestPVList(),
		NodeList:             cpmatest.CreateTestNodeList(),
		StorageClassList:     cpmatest.CreateStorageClassList(),
		NamespaceList:        cpmatest.CreateTestNameSpaceList(),
		RBACResources: api.RBACResources{
			UsersList:                      cpmatest.CreateUserList(),
			GroupList:                      cpmatest.CreateGroupList(),
			ClusterRolesList:               cpmatest.CreateClusterRoleList(),
			ClusterRolesBindingsList:       cpmatest.CreateClusterRoleBindingsList(),
			SecurityContextConstraintsList: cpmatest.CreateSCCList(),
		},
	}

	clusterExtraction := transform.ClusterExtraction{apiResources}

	actualClusterOutput, err := clusterExtraction.Transform()
	require.NoError(t, err)

	report := transform.ReportOutput{
		ClusterReport: actualClusterOutput[0].(transform.ReportOutput).ClusterReport,
	}

	actualClusterReportJSON, err := json.MarshalIndent(report, "", " ")
	require.NoError(t, err)

	expectedClusterReportJSON, err := ioutil.ReadFile("testdata/expected-report-cluster.json")
	require.NoError(t, err)

	assert.Equal(t, expectedClusterReportJSON, actualClusterReportJSON)
}

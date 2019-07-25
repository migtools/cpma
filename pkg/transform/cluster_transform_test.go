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
		QuotaList:            cpmatest.CreateTestClusterQuotaList(),
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

	manifests := actualClusterOutput[0].(transform.ManifestOutput).Manifests
	expectedClusterQuotaCRD, err := ioutil.ReadFile("testdata/expected-CR-cluster-quota.yaml")
	require.NoError(t, err)
	assert.Equal(t, "100_CPMA-cluster-quota-resource-test-quota1.yaml", manifests[0].Name)
	assert.Equal(t, expectedClusterQuotaCRD, manifests[0].CRD)
	expectedResourceQuotaCRD, err := ioutil.ReadFile("testdata/expected-CR-resource-quota.yaml")
	require.NoError(t, err)
	assert.Equal(t, "100_CPMA-namespacetest1-resource-quota-resourcequota1.yaml", manifests[1].Name)
	assert.Equal(t, expectedResourceQuotaCRD, manifests[1].CRD)

	report := transform.ReportOutput{
		ClusterReport: actualClusterOutput[1].(transform.ReportOutput).ClusterReport,
	}
	actualClusterReportJSON, err := json.MarshalIndent(report, "", " ")
	require.NoError(t, err)
	expectedClusterReportJSON, err := ioutil.ReadFile("testdata/expected-report-cluster.json")
	require.NoError(t, err)
	assert.Equal(t, expectedClusterReportJSON, actualClusterReportJSON)
}

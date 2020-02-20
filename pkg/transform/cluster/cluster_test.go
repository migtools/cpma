package cluster_test

import (
	"testing"

	"github.com/fusor/cpma/pkg/api"
	"github.com/fusor/cpma/pkg/transform"
	"github.com/fusor/cpma/pkg/transform/cluster"
	cpmatest "github.com/fusor/cpma/pkg/transform/internal/test"
	o7tapiauth "github.com/openshift/api/authorization/v1"
	o7tapiquota "github.com/openshift/api/quota/v1"
	o7tapiroute "github.com/openshift/api/route/v1"
	"github.com/stretchr/testify/assert"

	k8sapicore "k8s.io/api/core/v1"
	k8sapistorage "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestReportQuotas(t *testing.T) {
	testkey := resource.Quantity{
		Format: resource.DecimalSI,
	}
	testkey.Set(int64(99))

	expectedQuotaReports := make([]cluster.QuotaReport, 0)
	expectedQuotaReports = append(expectedQuotaReports, cluster.QuotaReport{
		Name: "test-quota1",
		Quota: k8sapicore.ResourceQuotaSpec{
			Hard: k8sapicore.ResourceList{
				"testkey": testkey},
		},
		Selector: o7tapiquota.ClusterResourceQuotaSelector{},
	})

	testCases := []struct {
		name                 string
		inputQuotaList       *o7tapiquota.ClusterResourceQuotaList
		expectedQuotaReports []cluster.QuotaReport
	}{
		{
			name:                 "generate cluster quota report",
			inputQuotaList:       cpmatest.CreateTestClusterQuotaList(),
			expectedQuotaReports: expectedQuotaReports,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clusterReport := &cluster.Report{}
			clusterReport.ReportQuotas(api.Resources{
				QuotaList: tc.inputQuotaList,
			})

			assert.Equal(t, tc.expectedQuotaReports, clusterReport.Quotas)
		})
	}
}

func TestReportNodes(t *testing.T) {
	clusterReport := &cluster.Report{}

	// Create node labels
	masterNodeLabels := make(map[string]string)
	masterNodeLabels["node-role.kubernetes.io/master"] = "true"

	masterNodeCapacity := make(k8sapicore.ResourceList)
	// Add CPU node usage
	cpuQuantity := resource.Quantity{
		Format: resource.DecimalSI,
	}
	cpuQuantity.Set(int64(2))
	masterNodeCapacity["cpu"] = cpuQuantity

	// Add node memory usage
	memoryQuantity := resource.Quantity{
		Format: resource.BinarySI,
	}
	memoryQuantity.Set(int64(2048))
	masterNodeCapacity["memory"] = memoryQuantity

	// Add pods
	podsQuantity := resource.Quantity{
		Format: resource.DecimalSI,
	}
	podsQuantity.Set(int64(10))
	masterNodeCapacity["pods"] = podsQuantity

	// Add resources that are available for scheduling
	allocatableResources := make(k8sapicore.ResourceList)

	allocatableMemoryQuantity := resource.Quantity{
		Format: resource.BinarySI,
	}
	allocatableMemoryQuantity.Set(int64(2048))
	masterNodeCapacity["memory"] = allocatableMemoryQuantity

	// Add pod count
	podList := &k8sapicore.PodList{}
	podList.Items = make([]k8sapicore.Pod, 0)
	podList.Items = append(podList.Items, k8sapicore.Pod{
		Spec: k8sapicore.PodSpec{
			NodeName: "test-master",
		},
	})
	podList.Items = append(podList.Items, k8sapicore.Pod{
		Spec: k8sapicore.PodSpec{
			NodeName: "test-master",
		},
	})
	podList.Items = append(podList.Items, k8sapicore.Pod{
		Spec: k8sapicore.PodSpec{
			NodeName: "test-master",
		},
	})
	podList.Items = append(podList.Items, k8sapicore.Pod{
		Spec: k8sapicore.PodSpec{
			NodeName: "not-this-node",
		},
	})

	namespaceList := make([]api.NamespaceResources, 0)
	namespaceList = append(namespaceList, api.NamespaceResources{
		PodList: podList,
	})

	// Init fake nodes
	nodes := make([]k8sapicore.Node, 0)
	nodes = append(nodes, k8sapicore.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "test-master",
			Labels: masterNodeLabels,
		},
		Status: k8sapicore.NodeStatus{
			Capacity:    masterNodeCapacity,
			Allocatable: allocatableResources,
		},
	})

	testCases := []struct {
		name                string
		actualClusterReport *cluster.Report
	}{
		{
			name:                "generate node report",
			actualClusterReport: clusterReport,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.actualClusterReport.ReportNodes(api.Resources{
				NodeList:      cpmatest.CreateTestNodeList(),
				NamespaceList: namespaceList,
			})

			assert.Equal(t, "test-master", tc.actualClusterReport.Nodes[0].Name)
			assert.Equal(t, true, tc.actualClusterReport.Nodes[0].MasterNode)
			assert.Equal(t, int64(2), tc.actualClusterReport.Nodes[0].Resources.CPU.Value())
			assert.Equal(t, int64(2048), tc.actualClusterReport.Nodes[0].Resources.MemoryConsumed.Value())
			assert.Equal(t, int64(2048), tc.actualClusterReport.Nodes[0].Resources.MemoryCapacity.Value())
			assert.Equal(t, int64(3), tc.actualClusterReport.Nodes[0].Resources.RunningPods.Value())
			assert.Equal(t, int64(10), tc.actualClusterReport.Nodes[0].Resources.PodCapacity.Value())
		})
	}
}

func TestReportPods(t *testing.T) {
	expectedPodRepors := make([]cluster.PodReport, 0)
	expectedPodRepors = append(expectedPodRepors, cluster.PodReport{Name: "test-pod1"})
	expectedPodRepors = append(expectedPodRepors, cluster.PodReport{Name: "test-pod2"})

	testCases := []struct {
		name              string
		inputPodList      *k8sapicore.PodList
		expectedPodRepors []cluster.PodReport
	}{
		{
			name:              "generate pod report",
			inputPodList:      cpmatest.CreateTestPodList(),
			expectedPodRepors: expectedPodRepors,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reportedNamespace := &cluster.NamespaceReport{}
			cluster.ReportPods(reportedNamespace, tc.inputPodList)
			assert.Equal(t, tc.expectedPodRepors, reportedNamespace.Pods)
		})
	}
}

func TestReportNamespaceResources(t *testing.T) {
	expectedCPU := &resource.Quantity{
		Format: resource.DecimalSI,
	}
	expectedCPU.Set(int64(2))
	expectedMemory := &resource.Quantity{
		Format: resource.BinarySI,
	}
	expectedMemory.Set(int64(2))

	expectedResources := &cluster.ContainerResourcesReport{
		ContainerCount: 2,
		CPUTotal:       expectedCPU,
		MemoryTotal:    expectedMemory,
	}

	testCases := []struct {
		name              string
		inputPodList      *k8sapicore.PodList
		expectedResources cluster.ContainerResourcesReport
	}{
		{
			name:              "generate resource report",
			inputPodList:      cpmatest.CreateTestPodResourceList(),
			expectedResources: *expectedResources,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reportedNamespace := &cluster.NamespaceReport{}
			cluster.ReportResources(reportedNamespace, tc.inputPodList)
			assert.Equal(t, tc.expectedResources, reportedNamespace.Resources)
		})
	}
}

func TestReportRoutes(t *testing.T) {
	alternateBackends := make([]o7tapiroute.RouteTargetReference, 0)
	alternateBackends = append(alternateBackends, o7tapiroute.RouteTargetReference{
		Kind: "testkind",
		Name: "testname",
	})

	to := o7tapiroute.RouteTargetReference{
		Kind: "testkindTo",
		Name: "testTo",
	}

	tls := &o7tapiroute.TLSConfig{
		Termination: o7tapiroute.TLSTerminationEdge,
	}

	expectedRouteReport := make([]cluster.RouteReport, 0)
	expectedRouteReport = append(expectedRouteReport, cluster.RouteReport{
		Name:              "route1",
		Host:              "testhost",
		Path:              "testpath",
		AlternateBackends: alternateBackends,
		TLS:               tls,
		To:                to,
		WildcardPolicy:    "None",
	})

	testCases := []struct {
		name                string
		inputRouteList      *o7tapiroute.RouteList
		expectedRouteReport []cluster.RouteReport
	}{
		{
			name:                "generate route report",
			inputRouteList:      cpmatest.CreateTestRouteList(),
			expectedRouteReport: expectedRouteReport,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reportedNamespace := &cluster.NamespaceReport{}
			cluster.ReportRoutes(reportedNamespace, tc.inputRouteList)
			assert.Equal(t, tc.expectedRouteReport, reportedNamespace.Routes)
		})
	}
}

func TestPVReport(t *testing.T) {
	resources := make(k8sapicore.ResourceList)
	cpu := resource.Quantity{
		Format: resource.DecimalSI,
	}
	cpu.Set(int64(1))
	resources["cpu"] = cpu

	memory := resource.Quantity{
		Format: resource.BinarySI,
	}
	memory.Set(int64(1))
	resources["memory"] = memory

	driver := k8sapicore.PersistentVolumeSource{
		NFS: &k8sapicore.NFSVolumeSource{
			Server: "example.com",
		},
	}

	expectedPVReport := make([]cluster.PVReport, 0)
	expectedPVReport = append(expectedPVReport, cluster.PVReport{
		Name:          "testpv",
		Driver:        driver,
		StorageClass:  "testclass",
		Capacity:      resources,
		Phase:         k8sapicore.VolumePending,
		ReclaimPolicy: k8sapicore.PersistentVolumeReclaimPolicy("testpolicy"),
	})

	testCases := []struct {
		name             string
		inputPVList      *k8sapicore.PersistentVolumeList
		expectedPVReport []cluster.PVReport
	}{
		{
			name:             "generate route report",
			inputPVList:      cpmatest.CreateTestPVList(),
			expectedPVReport: expectedPVReport,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clusterReport := &cluster.Report{}
			clusterReport.ReportPVs(api.Resources{
				PersistentVolumeList: tc.inputPVList,
			})

			assert.Equal(t, tc.expectedPVReport, clusterReport.PVs)
		})
	}
}

func TestStorageClassReport(t *testing.T) {
	expectedStorageClassReport := make([]cluster.StorageClassReport, 0)
	expectedStorageClassReport = append(expectedStorageClassReport, cluster.StorageClassReport{
		Name:        "testclass",
		Provisioner: "testprovisioner",
	})

	testCases := []struct {
		name                       string
		inputStorageClassList      *k8sapistorage.StorageClassList
		expectedStorageClassReport []cluster.StorageClassReport
	}{
		{
			name:                       "generate storage class report",
			inputStorageClassList:      cpmatest.CreateStorageClassList(),
			expectedStorageClassReport: expectedStorageClassReport,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clusterReport := &cluster.Report{}
			clusterReport.ReportStorageClasses(api.Resources{
				StorageClassList: tc.inputStorageClassList,
			})

			assert.Equal(t, tc.expectedStorageClassReport, clusterReport.StorageClasses)
		})
	}
}

func TestRBACReport(t *testing.T) {
	userList := cpmatest.CreateUserList()
	groupList := cpmatest.CreateGroupList()
	clusterRoleList := cpmatest.CreateClusterRoleList()
	clusterRoleBindingsList := cpmatest.CreateClusterRoleBindingsList()
	sccList := cpmatest.CreateSCCList()

	roleList := &o7tapiauth.RoleList{}
	roleList.Items = make([]o7tapiauth.Role, 0)

	roleList.Items = append(roleList.Items, o7tapiauth.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testrole1",
		},
	})

	namespaceList := make([]api.NamespaceResources, 0)
	namespaceList = append(namespaceList, api.NamespaceResources{
		NamespaceName: "testnamespace1",
		RolesList:     roleList,
	})

	expectedUsers := make([]cluster.OpenshiftUser, 0)
	expectedUsers = append(expectedUsers, cluster.OpenshiftUser{
		Name:       "testuser1",
		FullName:   "full name1",
		Identities: []string{"test-identity1", "test-identity2"},
		Groups:     []string{"group1", "group2"},
	})
	expectedUsers = append(expectedUsers, cluster.OpenshiftUser{
		Name:       "testuser2",
		FullName:   "full name2",
		Identities: []string{"test-identity1", "test-identity2"},
		Groups:     []string{"group1", "group2"},
	})

	expectedGroups := make([]cluster.OpenshiftGroup, 0)
	expectedGroups = append(expectedGroups, cluster.OpenshiftGroup{
		Name:  "testgroup1",
		Users: []string{"testuser1"},
	})
	expectedGroups = append(expectedGroups, cluster.OpenshiftGroup{
		Name:  "testgroup2",
		Users: []string{"testuser2"},
	})

	expectedRoles := make([]cluster.OpenshiftRole, 0)
	expectedRoles = append(expectedRoles, cluster.OpenshiftRole{
		Name: "testrole1",
	})

	expectedNamespaceRoles := make([]cluster.OpenshiftNamespaceRole, 0)
	expectedNamespaceRoles = append(expectedNamespaceRoles, cluster.OpenshiftNamespaceRole{
		Namespace: "testnamespace1",
		Roles:     expectedRoles,
	})

	expectedClusterRoles := make([]cluster.OpenshiftClusterRole, 0)
	expectedClusterRoles = append(expectedClusterRoles, cluster.OpenshiftClusterRole{
		Name: "testrole1",
	})

	expectedClusterRoleBindings := make([]cluster.OpenshiftClusterRoleBinding, 0)
	expectedClusterRoleBindings = append(expectedClusterRoleBindings, cluster.OpenshiftClusterRoleBinding{
		Name:       "testbinding1",
		UserNames:  []string{"testuser1"},
		GroupNames: []string{"testgroup1"},
	})

	expectedSCC := make([]cluster.OpenshiftSecurityContextConstraints, 0)
	expectedSCC = append(expectedSCC, cluster.OpenshiftSecurityContextConstraints{
		Name:       "testscc1",
		Users:      []string{"testuser1", "testrole:serviceaccount:testnamespace1:testsa"},
		Groups:     []string{"testgroup1"},
		Namespaces: []string{"testnamespace1"},
	})

	expectedRBACReport := &cluster.RBACReport{
		Users:                      expectedUsers,
		Groups:                     expectedGroups,
		Roles:                      expectedNamespaceRoles,
		ClusterRoles:               expectedClusterRoles,
		ClusterRoleBindings:        expectedClusterRoleBindings,
		SecurityContextConstraints: expectedSCC,
	}

	testCases := []struct {
		name               string
		inputRBAC          api.Resources
		expectedRBACReport cluster.RBACReport
	}{
		{
			name: "generate RBAC report",
			inputRBAC: api.Resources{
				NamespaceList: namespaceList,
				RBACResources: api.RBACResources{
					UsersList:                      userList,
					GroupList:                      groupList,
					ClusterRolesList:               clusterRoleList,
					ClusterRolesBindingsList:       clusterRoleBindingsList,
					SecurityContextConstraintsList: sccList,
				},
			},
			expectedRBACReport: *expectedRBACReport,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clusterReport := &cluster.Report{}
			clusterReport.Namespaces = make([]cluster.NamespaceReport, 0)
			clusterReport.Namespaces = append(clusterReport.Namespaces, cluster.NamespaceReport{
				Name: "testnamespace1",
			})
			clusterReport.ReportRBAC(tc.inputRBAC)

			assert.Equal(t, tc.expectedRBACReport, clusterReport.RBACReport)
		})
	}

}

func TestReportMisssingGVs(t *testing.T) {
	expectedMissingGVs := make([]cluster.NewGVsReport, 0)
	expectedMissingGVs = append(expectedMissingGVs, cluster.NewGVsReport{
		GroupVersion: "testgroupversion/v1",
	})

	testCases := []struct {
		name                  string
		inputGroupVersions    *metav1.APIGroupList
		inputDstGroupVersions *metav1.APIGroupList
		expectedNewGVs        []cluster.NewGVsReport
	}{
		{
			name:                  "generate missing groupversions report",
			inputGroupVersions:    cpmatest.CreateTestClusterGroupVersions("testgroupversion", "v1beta1"),
			inputDstGroupVersions: cpmatest.CreateTestClusterGroupVersions("testgroupversion", "v1"),
			expectedNewGVs:        expectedMissingGVs,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clusterNewGVs := &cluster.Report{}
			clusterNewGVs.ReportNewGVs(transform.NewGroupVersions(tc.inputGroupVersions, tc.inputDstGroupVersions))
			assert.Equal(t, tc.expectedNewGVs, clusterNewGVs.NewGVs)
		})
	}
}

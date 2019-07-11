package cluster_test

import (
	"testing"

	"github.com/fusor/cpma/pkg/api"
	"github.com/fusor/cpma/pkg/transform/cluster"
	cpmatest "github.com/fusor/cpma/pkg/transform/internal/test"
	O7tapiroute "github.com/openshift/api/route/v1"
	"github.com/stretchr/testify/assert"
	k8sapicore "k8s.io/api/core/v1"
	k8sapistorage "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	k8smachinery "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
		ObjectMeta: k8smachinery.ObjectMeta{
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
	alternateBackends := make([]O7tapiroute.RouteTargetReference, 0)
	alternateBackends = append(alternateBackends, O7tapiroute.RouteTargetReference{
		Kind: "testkind",
		Name: "testname",
	})

	to := O7tapiroute.RouteTargetReference{
		Kind: "testkindTo",
		Name: "testTo",
	}

	tls := &O7tapiroute.TLSConfig{
		Termination: O7tapiroute.TLSTerminationEdge,
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
		inputRouteList      *O7tapiroute.RouteList
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
		Name:         "testpv",
		Driver:       driver,
		StorageClass: "testclass",
		Capacity:     resources,
		Phase:        k8sapicore.VolumePending,
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

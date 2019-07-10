package cluster

import (
	"github.com/fusor/cpma/pkg/api"
	O7tapiroute "github.com/openshift/api/route/v1"
	"github.com/sirupsen/logrus"
	k8sapicore "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	k8smachinery "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Report represents json report of k8s resources
type Report struct {
	Nodes          []NodeReport         `json:"nodes"`
	Namespaces     []NamespaceReport    `json:"namespaces,omitempty"`
	PVs            []PVReport           `json:"pvs,omitempty"`
	StorageClasses []StorageClassReport `json:"storageClasses,omitempty"`
}

// NodeReport represents json report of k8s nodes
type NodeReport struct {
	Name       string        `json:"name"`
	MasterNode bool          `json:"masterNode"`
	Resources  NodeResources `json:"resources"`
}

// NodeResources represents a json report of Node resources
type NodeResources struct {
	CPU            *resource.Quantity `json:"cpu"`
	MemoryConsumed *resource.Quantity `json:"memoryConsumed"`
	MemoryCapacity *resource.Quantity `json:"memoryCapacity"`
	RunningPods    *resource.Quantity `json:"runningPods"`
	PodCapacity    *resource.Quantity `json:"podCapacity"`
}

// NamespaceReport represents json report of k8s namespaces
type NamespaceReport struct {
	Name         string                   `json:"name"`
	LatestChange k8smachinery.Time        `json:"latestChange,omitempty"`
	Resources    ContainerResourcesReport `json:"resources,omitempty"`
	Pods         []PodReport              `json:"pods,omitempty"`
	Routes       []RouteReport            `json:"routes,omitempty"`
}

// PodReport represents json report of k8s pods
type PodReport struct {
	Name string `json:"name"`
}

// RouteReport represents json report of k8s pods
type RouteReport struct {
	Name              string                             `json:"name"`
	Host              string                             `json:"host"`
	Path              string                             `json:"path,omitempty"`
	AlternateBackends []O7tapiroute.RouteTargetReference `json:"alternateBackends,omitempty"`
	TLS               *O7tapiroute.TLSConfig             `json:"tls,omitempty"`
	To                O7tapiroute.RouteTargetReference   `json:"to,omitempty"`
	WildcardPolicy    O7tapiroute.WildcardPolicyType     `json:"wildcardPolicy"`
}

// ContainerResourcesReport represents json report for aggregated container resources
type ContainerResourcesReport struct {
	ContainerCount int                `json:"containerCount"`
	CPUTotal       *resource.Quantity `json:"cpuTotal"`
	MemoryTotal    *resource.Quantity `json:"memoryTotal"`
}

// PVReport represents json report of k8s PVs
type PVReport struct {
	Name         string                            `json:"name"`
	Driver       k8sapicore.PersistentVolumeSource `json:"driver"`
	StorageClass string                            `json:"storageClass,omitempty"`
	Capacity     k8sapicore.ResourceList           `json:"capacity,omitempty"`
	Phase        k8sapicore.PersistentVolumePhase  `json:"phase,omitempty"`
}

// StorageClassReport represents json report of k8s storage classes
type StorageClassReport struct {
	Name        string `json:"name"`
	Provisioner string `json:"provisioner"`
}

// GenClusterReport inserts report values into structures for json output
func GenClusterReport(apiResources api.Resources) (clusterReport Report) {
	clusterReport.ReportNodes(apiResources)
	clusterReport.ReportNamespaces(apiResources)
	clusterReport.ReportPVs(apiResources)
	clusterReport.ReportStorageClasses(apiResources)
	return
}

// ReportNodes fills in information about nodes
func (clusterReport *Report) ReportNodes(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportNodes")

	for _, node := range apiResources.NodeList.Items {
		nodeReport := &NodeReport{
			Name: node.ObjectMeta.Name,
		}

		isMaster, ok := node.ObjectMeta.Labels["node-role.kubernetes.io/master"]
		nodeReport.MasterNode = ok && isMaster == "true"

		ReportNodeResources(nodeReport, node.Status, apiResources)
		clusterReport.Nodes = append(clusterReport.Nodes, *nodeReport)
	}
}

// ReportNodeResources parse and insert info about consumed resources
func ReportNodeResources(repotedNode *NodeReport, nodeStatus k8sapicore.NodeStatus, apiResources api.Resources) {
	repotedNode.Resources.CPU = nodeStatus.Capacity.Cpu()

	repotedNode.Resources.MemoryCapacity = nodeStatus.Capacity.Memory()

	memConsumed := new(resource.Quantity)
	memCapacity, _ := nodeStatus.Capacity.Memory().AsInt64()
	memAllocatable, _ := nodeStatus.Allocatable.Memory().AsInt64()
	memConsumed.Set(memCapacity - memAllocatable)
	memConsumed.Format = resource.BinarySI
	repotedNode.Resources.MemoryConsumed = memConsumed

	var runningPodsCount int64
	for _, resources := range apiResources.NamespaceList {
		for _, pod := range resources.PodList.Items {
			if pod.Spec.NodeName == repotedNode.Name {
				runningPodsCount++
			}
		}
	}
	podsRunning := new(resource.Quantity)
	podsRunning.Set(runningPodsCount)
	podsRunning.Format = resource.DecimalSI
	repotedNode.Resources.RunningPods = podsRunning

	repotedNode.Resources.PodCapacity = nodeStatus.Capacity.Pods()
}

// ReportNamespaces fills in information about Namespaces
func (clusterReport *Report) ReportNamespaces(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportNamespaces")

	for _, resources := range apiResources.NamespaceList {
		reportedNamespace := NamespaceReport{Name: resources.NamespaceName}

		ReportPods(&reportedNamespace, resources.PodList)
		ReportResources(&reportedNamespace, resources.PodList)
		ReportRoutes(&reportedNamespace, resources.RouteList)
		clusterReport.Namespaces = append(clusterReport.Namespaces, reportedNamespace)
	}
}

// ReportPods creates info about cluster pods
func ReportPods(reportedNamespace *NamespaceReport, podList *k8sapicore.PodList) {
	for _, pod := range podList.Items {
		reportedPod := &PodReport{Name: pod.Name}
		reportedNamespace.Pods = append(reportedNamespace.Pods, *reportedPod)

		// Update namespace touch timestamp
		if pod.ObjectMeta.CreationTimestamp.Time.Unix() > reportedNamespace.LatestChange.Time.Unix() {
			reportedNamespace.LatestChange = pod.ObjectMeta.CreationTimestamp
		}
	}
}

// ReportResources create report about namespace resources
func ReportResources(reportedNamespace *NamespaceReport, podList *k8sapicore.PodList) {
	resources := ContainerResourcesReport{
		CPUTotal:    &resource.Quantity{Format: resource.DecimalSI},
		MemoryTotal: &resource.Quantity{Format: resource.BinarySI},
	}
	reportedNamespace.Resources = resources

	for _, pod := range podList.Items {
		ReportContainerResources(reportedNamespace, &pod)
	}
}

// ReportContainerResources create report about container resources
func ReportContainerResources(reportedNamespace *NamespaceReport, pod *k8sapicore.Pod) {
	cpuTotal := reportedNamespace.Resources.CPUTotal.Value()
	memoryTotal := reportedNamespace.Resources.MemoryTotal.Value()

	for _, container := range pod.Spec.Containers {
		cpuTotal += container.Resources.Requests.Cpu().Value()
		memoryTotal += container.Resources.Requests.Memory().Value()
	}
	reportedNamespace.Resources.CPUTotal.Set(cpuTotal)
	reportedNamespace.Resources.MemoryTotal.Set(memoryTotal)
	reportedNamespace.Resources.ContainerCount += len(pod.Spec.Containers)
}

// ReportRoutes create report about routes
func ReportRoutes(reportedNamespace *NamespaceReport, routeList *O7tapiroute.RouteList) {
	for _, route := range routeList.Items {
		reportedRoute := &RouteReport{
			Name:              route.Name,
			AlternateBackends: route.Spec.AlternateBackends,
			Host:              route.Spec.Host,
			Path:              route.Spec.Path,
			To:                route.Spec.To,
			TLS:               route.Spec.TLS,
			WildcardPolicy:    route.Spec.WildcardPolicy,
		}

		reportedNamespace.Routes = append(reportedNamespace.Routes, *reportedRoute)
	}
}

// ReportPVs create report oabout pvs
func (clusterReport *Report) ReportPVs(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportPVs")
	pvList := apiResources.PersistentVolumeList

	// Go through all PV and save required information to report
	for _, pv := range pvList.Items {
		reportedPV := &PVReport{
			Name:         pv.Name,
			Driver:       pv.Spec.PersistentVolumeSource,
			StorageClass: pv.Spec.StorageClassName,
			Capacity:     pv.Spec.Capacity,
			Phase:        pv.Status.Phase,
		}

		clusterReport.PVs = append(clusterReport.PVs, *reportedPV)
	}
}

// ReportStorageClasses create report about storage classes
func (clusterReport *Report) ReportStorageClasses(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportStorageClasses")
	// Go through all storage classes and save required information to report
	storageClassList := apiResources.StorageClassList
	for _, storageClass := range storageClassList.Items {
		reportedStorageClass := &StorageClassReport{
			Name:        storageClass.Name,
			Provisioner: storageClass.Provisioner,
		}

		clusterReport.StorageClasses = append(clusterReport.StorageClasses, *reportedStorageClass)
	}
}

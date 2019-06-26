package transform

import (
	"strconv"

	"github.com/fusor/cpma/pkg/api"
	O7tapiroute "github.com/openshift/api/route/v1"
	"github.com/sirupsen/logrus"

	k8sapicore "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// ClusterReportName is the cluster report name
const ClusterReportName = "ClusterReport"

// ClusterReportExtraction holds data extracted from k8s API resources
type ClusterReportExtraction struct {
	api.Resources
}

// ClusterTransform reprents transform for k8s API resources
type ClusterTransform struct {
}

// Transform converts data collected from an OCP3 API into a useful output
func (e ClusterReportExtraction) Transform() ([]Output, error) {
	logrus.Info("ClusterTransform::Transform")

	clusterReport, err := genClusterReport(api.Resources{
		PersistentVolumeList: e.PersistentVolumeList,
		StorageClassList:     e.StorageClassList,
		NamespaceMap:         e.NamespaceMap,
		NodeList:             e.NodeList,
	})
	if err != nil {
		return nil, err
	}

	output := ReportOutput{
		ClusterReport: clusterReport,
	}

	outputs := []Output{output}
	return outputs, nil
}

// genClusterReport inserts report values into structures for json output
func genClusterReport(apiResources api.Resources) (clusterReport ClusterReport, err error) {
	clusterReport.reportNodes(apiResources)
	clusterReport.reportNamespaces(apiResources)
	clusterReport.reportPVs(apiResources)
	clusterReport.reportStorageClasses(apiResources)
	return
}

// reportNodes fills in information about nodes
func (clusterReport *ClusterReport) reportNodes(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportNodes")

	for _, node := range apiResources.NodeList.Items {
		nodeReport := &NodeReport{
			Name: node.ObjectMeta.Name,
		}

		isMaster, ok := node.ObjectMeta.Labels["node-role.kubernetes.io/master"]
		if ok && isMaster == "true" {
			nodeReport.MasterNode = true
		} else {
			nodeReport.MasterNode = false
		}

		reportResources(nodeReport, node.Status, apiResources)
		clusterReport.Nodes = append(clusterReport.Nodes, *nodeReport)
	}
}

// reportResources parse and insert info about consumed resources
func reportResources(repotedNode *NodeReport, nodeStatus k8sapicore.NodeStatus, apiResources api.Resources) {
	repotedNode.Resources.CPU = nodeStatus.Capacity.Cpu().String()

	repotedNode.Resources.MemoryCapacity = nodeStatus.Capacity.Memory().String()

	memConsumed := new(resource.Quantity)
	memCapacity, _ := nodeStatus.Capacity.Memory().AsInt64()
	memAllocatable, _ := nodeStatus.Allocatable.Memory().AsInt64()
	memConsumed.Set(memCapacity - memAllocatable)
	memConsumed.Format = resource.BinarySI
	repotedNode.Resources.MemoryConsumed = memConsumed.String()

	runningPods := 0
	for _, resources := range apiResources.NamespaceMap {
		for _, pod := range resources.PodList.Items {
			if pod.Spec.NodeName == repotedNode.Name {
				runningPods++
			}
		}
	}
	repotedNode.Resources.RunningPods = strconv.Itoa(runningPods) + "/" + nodeStatus.Capacity.Pods().String()
}

// reportNamespaces fills in information about Namespaces
func (clusterReport *ClusterReport) reportNamespaces(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportNamespaces")

	for namespaceName, resources := range apiResources.NamespaceMap {
		reportedNamespace := NamespaceReport{
			Name: namespaceName,
		}

		reportPods(&reportedNamespace, resources.PodList)
		reportRoutes(&reportedNamespace, resources.RouteList)

		clusterReport.Namespaces = append(clusterReport.Namespaces, reportedNamespace)
	}
}

// reportPods creates info about cluster pods
func reportPods(reportedNamespace *NamespaceReport, podList *k8sapicore.PodList) {
	for _, pod := range podList.Items {
		reportedPod := &PodReport{
			Name: pod.Name,
		}

		reportedNamespace.Pods = append(reportedNamespace.Pods, *reportedPod)
	}
}

func reportRoutes(reportedNamespace *NamespaceReport, routeList *O7tapiroute.RouteList) {
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

func (clusterReport *ClusterReport) reportPVs(apiResources api.Resources) {
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

func (clusterReport *ClusterReport) reportStorageClasses(apiResources api.Resources) {
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

// Validate no need to validate it, data is exctracted from API
func (e ClusterReportExtraction) Validate() error {
	return nil
}

// Extract collects data for cluster report
func (e ClusterTransform) Extract() (Extraction, error) {
	extraction := &ClusterReportExtraction{}

	nodeList, err := api.ListNodes()
	if err != nil {
		return nil, err
	}
	extraction.NodeList = nodeList

	namespacesList, err := api.ListNamespaces()
	if err != nil {
		return nil, err
	}

	// Map all namespaces to their resources
	extraction.NamespaceMap = make(map[string]*api.NamespaceResources)
	for _, namespace := range namespacesList.Items {
		namespaceResources := &api.NamespaceResources{}

		podsList, err := api.ListPods(namespace.Name)
		if err != nil {
			return nil, err
		}
		namespaceResources.PodList = podsList

		routesList, err := api.ListRoutes(namespace.Name)
		if err != nil {
			return nil, err
		}
		namespaceResources.RouteList = routesList

		extraction.NamespaceMap[namespace.Name] = namespaceResources
	}

	pvList, err := api.ListPVs()
	if err != nil {
		return nil, err
	}
	extraction.PersistentVolumeList = pvList

	storageClassList, err := api.ListStorageClasses()
	if err != nil {
		return nil, err
	}
	extraction.StorageClassList = storageClassList

	return *extraction, nil
}

// Name returns a human readable name for the transform
func (e ClusterTransform) Name() string {
	return ClusterReportName
}

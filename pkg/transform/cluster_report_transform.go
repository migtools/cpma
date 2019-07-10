package transform

import (
	"github.com/fusor/cpma/pkg/api"
	"github.com/fusor/cpma/pkg/transform/cluster"
	"github.com/sirupsen/logrus"
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

	clusterReport := cluster.GenClusterReport(api.Resources{
		PersistentVolumeList: e.PersistentVolumeList,
		StorageClassList:     e.StorageClassList,
		NamespaceList:        e.NamespaceList,
		NodeList:             e.NodeList,
	})

	output := ReportOutput{
		ClusterReport: clusterReport,
	}

	outputs := []Output{output}
	return outputs, nil
}

// Validate no need to validate it, data is exctracted from API
func (e ClusterReportExtraction) Validate() (err error) { return }

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
	namespaceListSize := len(namespacesList.Items)
	extraction.NamespaceList = make([]api.NamespaceResources, namespaceListSize, namespaceListSize)
	for i, namespace := range namespacesList.Items {
		namespaceResources := api.NamespaceResources{NamespaceName: namespace.Name}

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

		extraction.NamespaceList[i] = namespaceResources
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

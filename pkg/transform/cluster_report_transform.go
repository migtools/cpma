package transform

import (
	"github.com/fusor/cpma/pkg/api"
	"github.com/fusor/cpma/pkg/transform/clusterreport"
)

// ClusterReportName is the cluster report name
const ClusterReportName = "ClusterReport"

// ClusterReportExtraction holds data extracted from k8s API resources
type ClusterReportExtraction struct {
	api.Resources
}

// ClusterReportTransform reprents transform for k8s API resources
type ClusterReportTransform struct {
}

// Transform converts data collected from an OCP3 API into a useful output
func (e ClusterReportExtraction) Transform() ([]Output, error) {
	clusterReport, err := clusterreport.Report(api.Resources{
		PersistentVolumeList: e.PersistentVolumeList,
		StorageClassList:     e.StorageClassList,
		NamespaceMap:         e.NamespaceMap,
	})

	if err != nil {
		return nil, err
	}

	output := ClusterOutput{clusterReport}

	outputs := []Output{output}
	return outputs, nil
}

// Validate no need to validate it, data is exctracted from API
func (e ClusterReportExtraction) Validate() error {
	return nil
}

// Extract collects data for cluster report
func (e ClusterReportTransform) Extract() (Extraction, error) {
	extraction := &ClusterReportExtraction{}

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
func (e ClusterReportTransform) Name() string {
	return ClusterReportName
}

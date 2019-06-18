package transform

import (
	"github.com/fusor/cpma/pkg/api"
	"github.com/sirupsen/logrus"
	k8sapicore "k8s.io/api/core/v1"
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
	logrus.Info("ClusterReportTransform::Transform")

	clusterReport, err := genClusterReport(api.Resources{
		PersistentVolumeList: e.PersistentVolumeList,
		StorageClassList:     e.StorageClassList,
		NamespaceMap:         e.NamespaceMap,
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

func genClusterReport(apiResources api.Resources) (ClusterReport, error) {
	clusterReport := ClusterReport{}

	clusterReport.reportNamespaces(apiResources)

	clusterReport.reportPVs(apiResources)

	clusterReport.reportStorageClasses(apiResources)

	return clusterReport, nil
}

func (clusterReport *ClusterReport) reportNamespaces(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportNamespaces")

	// Go through all required namespace resources and report them
	for namespaceName, resources := range apiResources.NamespaceMap {
		reportedNamespace := NamespaceReport{
			Name: namespaceName,
		}

		reportPods(&reportedNamespace, resources.PodList)

		clusterReport.Namespaces = append(clusterReport.Namespaces, reportedNamespace)
	}
}

func reportPods(reportedNamespace *NamespaceReport, podList *k8sapicore.PodList) {
	for _, pod := range podList.Items {
		reportedPod := &PodReport{
			Name: pod.Name,
		}

		reportedNamespace.Pods = append(reportedNamespace.Pods, *reportedPod)
	}
}

func (clusterReport *ClusterReport) reportPVs(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportPVs")
	pvList := apiResources.PersistentVolumeList

	// Go through all PV and save required information to report
	for _, pv := range pvList.Items {
		reportedPV := &PVReport{
			Name:         pv.Name,
			StorageClass: pv.Spec.StorageClassName,
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

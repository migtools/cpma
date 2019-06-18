package clusterreport

import (
	"github.com/fusor/cpma/pkg/api"
	"github.com/sirupsen/logrus"
	k8sapicore "k8s.io/api/core/v1"
)

// ClusterReport represents json report of k8s resources
type ClusterReport struct {
	Namespaces     []Namespace    `json:"namespaces,omitempty"`
	PVs            []PV           `json:"pvs,omitempty"`
	StorageClasses []StorageClass `json:"storageClasses,omitempty"`
}

// Namespace represents json report of k8s namespaces
type Namespace struct {
	Name string `json:"name"`
	Pods []Pod  `json:"pods,omitempty"`
}

// Pod represents json report of k8s pods
type Pod struct {
	Name string `json:"name"`
}

// PV represents json report of k8s PVs
type PV struct {
	Name         string `json:"name"`
	StorageClass string `json:"storageClass,omitempty"`
}

// StorageClass represents json report of k8s storage classes
type StorageClass struct {
	Name        string `json:"name"`
	Provisioner string `json:"provisioner"`
}

// Report collecting data about OCP3 resources
func Report(apiResources api.Resources) (*ClusterReport, error) {
	clusterReport := &ClusterReport{}

	clusterReport.reportNamespaces(apiResources)

	clusterReport.reportPVs(apiResources)

	clusterReport.reportStorageClasses(apiResources)

	return clusterReport, nil
}

func (cluserReport *ClusterReport) reportNamespaces(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportNamespaces")

	// Go through all required namespace resources and report them
	for namespaceName, resources := range apiResources.NamespaceMap {
		reportedNamespace := Namespace{
			Name: namespaceName,
		}

		reportPods(&reportedNamespace, resources.PodList)

		cluserReport.Namespaces = append(cluserReport.Namespaces, reportedNamespace)
	}
}

func reportPods(reportedNamespace *Namespace, podList *k8sapicore.PodList) {
	for _, pod := range podList.Items {
		reportedPod := &Pod{
			Name: pod.Name,
		}

		reportedNamespace.Pods = append(reportedNamespace.Pods, *reportedPod)
	}
}

func (cluserReport *ClusterReport) reportPVs(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportPVs")
	pvList := apiResources.PersistentVolumeList
	// Go through all PV and save required information to report
	for _, pv := range pvList.Items {
		reportedPV := &PV{
			Name:         pv.Name,
			StorageClass: pv.Spec.StorageClassName,
		}

		cluserReport.PVs = append(cluserReport.PVs, *reportedPV)
	}
}

func (cluserReport *ClusterReport) reportStorageClasses(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportStorageClasses")
	// Go through all storage classes and save required information to report
	storageClassList := apiResources.StorageClassList
	for _, storageClass := range storageClassList.Items {
		reportedStorageClass := &StorageClass{
			Name:        storageClass.Name,
			Provisioner: storageClass.Provisioner,
		}

		cluserReport.StorageClasses = append(cluserReport.StorageClasses, *reportedStorageClass)
	}
}

package clusterreport

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/fusor/cpma/pkg/api"
	"github.com/fusor/cpma/pkg/env"
	"github.com/sirupsen/logrus"
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

// Start collecting data about OCP3 resources
func Start() error {
	clusterReport := &ClusterReport{}

	if err := clusterReport.reportNamespaces(); err != nil {
		return err
	}

	if err := clusterReport.reportPVs(); err != nil {
		return err
	}

	if err := clusterReport.dumpToJSON(); err != nil {
		return err
	}

	return nil
}

func (cluserReport *ClusterReport) reportNamespaces() error {
	logrus.Debug("ClusterReport::ReportNamespaces")
	namespacesList, err := api.ListNamespaces()
	if err != nil {
		return err
	}

	// get namespaces names as a slice
	namespacesNames := make([]string, 0, len(namespacesList.Items))
	for _, namespace := range namespacesList.Items {
		namespacesNames = append(namespacesNames, namespace.Name)
	}

	// Go through all required namespace resources and report them
	for _, namespaceName := range namespacesNames {
		reportedNamespace := Namespace{
			Name: namespaceName,
		}
		reportPods(namespaceName, &reportedNamespace)

		cluserReport.Namespaces = append(cluserReport.Namespaces, reportedNamespace)
	}

	return nil
}

func reportPods(namespaceName string, reportedNamespace *Namespace) error {
	podsList, err := api.ListPods(namespaceName)
	if err != nil {
		return err
	}

	for _, pod := range podsList.Items {
		reportedPod := &Pod{
			Name: pod.Name,
		}

		reportedNamespace.Pods = append(reportedNamespace.Pods, *reportedPod)
	}

	return nil
}

func (cluserReport *ClusterReport) reportPVs() error {
	logrus.Debug("ClusterReport::ReportPVs")
	pvList, err := api.ListPVs()
	if err != nil {
		return err
	}

	// Go through all PV and save required information to report
	for _, pv := range pvList.Items {
		reportedPV := &PV{
			Name:         pv.Name,
			StorageClass: pv.Spec.StorageClassName,
		}

		cluserReport.PVs = append(cluserReport.PVs, *reportedPV)
	}

	return nil
}

func (cluserReport *ClusterReport) reportStorageClasses() error {
	logrus.Debug("ClusterReport::ReportStorageClasses")
	storageClassList, err := api.ListStorageClasses()
	if err != nil {
		return err
	}

	// Go through all storage classes and save required information to report
	for _, storageClass := range storageClassList.Items {
		reportedStorageClass := &StorageClass{
			Name:        storageClass.Name,
			Provisioner: storageClass.Provisioner,
		}

		cluserReport.StorageClasses = append(cluserReport.StorageClasses, *reportedStorageClass)
	}

	return nil
}

func (cluserReport *ClusterReport) dumpToJSON() error {
	jsonFile := filepath.Join(env.Config().GetString("OutputDir"), "cluster-report.json")

	file, err := json.MarshalIndent(&cluserReport, "", " ")
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(jsonFile, file, 0644); err != nil {
		return err
	}

	logrus.Debugf("Cluster report added to %s", jsonFile)
	return nil
}

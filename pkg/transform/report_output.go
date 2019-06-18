package transform

import (
	"encoding/json"

	"github.com/fusor/cpma/pkg/io"
	"github.com/sirupsen/logrus"
)

// ReportOutput holds a collection of reports to be written to fil

type ReportOutput struct {
	ClusterReport    ClusterReport     `json:"cluster"`
	ComponentReports []ComponentReport `json:"components"`
}

type ComponentReport struct {
	Component string   `json:"component"`
	Reports   []Report `json:"reports"`
}

// ClusterReport represents json report of k8s resources
type ClusterReport struct {
	Namespaces     []NamespaceReport    `json:"namespaces,omitempty"`
	PVs            []PVReport           `json:"pvs,omitempty"`
	StorageClasses []StorageClassReport `json:"storageClasses,omitempty"`
}

// NamespaceReport represents json report of k8s namespaces
type NamespaceReport struct {
	Name string      `json:"name"`
	Pods []PodReport `json:"pods,omitempty"`
}

// PodReport represents json report of k8s pods
type PodReport struct {
	Name string `json:"name"`
}

// PVReport represents json report of k8s PVs
type PVReport struct {
	Name         string `json:"name"`
	StorageClass string `json:"storageClass,omitempty"`
}

// StorageClassReport represents json report of k8s storage classes
type StorageClassReport struct {
	Name        string `json:"name"`
	Provisioner string `json:"provisioner"`
}

// ReportOutputFlush flush reports to disk
var ReportOutputFlush = func(r ReportOutput) error {
	logrus.Info("Flushing reports to disk")
	DumpReports(r)
	return nil
}

// Flush reports to files
func (r ReportOutput) Flush() error {
	return ReportOutputFlush(r)
}

// DumpReports creates OCDs files
func DumpReports(r ReportOutput) {
	var existingReports ReportOutput

	jsonFile := "report.json"

	jsonData, err := io.ReadFile(jsonFile)
	if err != nil {
		logrus.Errorf("unable to read to report file: %s", jsonFile)
	}

	err = json.Unmarshal(jsonData, &existingReports)
	if err != nil {
		logrus.Errorf("unable to unmarshal existing report json")
	}

	for _, namespace := range r.ClusterReport.Namespaces {
		existingReports.ClusterReport.Namespaces = append(existingReports.ClusterReport.Namespaces, namespace)
	}

	for _, pv := range r.ClusterReport.PVs {
		existingReports.ClusterReport.PVs = append(existingReports.ClusterReport.PVs, pv)
	}

	for _, sc := range r.ClusterReport.StorageClasses {
		existingReports.ClusterReport.StorageClasses = append(existingReports.ClusterReport.StorageClasses, sc)
	}

	for _, componentReport := range r.ComponentReports {
		existingReports.ComponentReports = append(existingReports.ComponentReports, componentReport)
	}

	jsonReports, err := json.MarshalIndent(existingReports, "", " ")
	if err != nil {
		logrus.Errorf("unable to marshal reports")
	}

	err = io.WriteFile(jsonReports, jsonFile)
	if err != nil {
		logrus.Errorf("unable to write to report file: %s", jsonFile)
	}
}

package transform

import (
	"encoding/json"

	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/cluster"
	"github.com/sirupsen/logrus"
)

// ReportOutput holds a collection of reports to be written to file
type ReportOutput struct {
	ClusterReport    cluster.Report    `json:"cluster,omitempty"`
	ComponentReports []ComponentReport `json:"components,omitempty"`
}

// ComponentReport holds a collection of ocp3 config reports
type ComponentReport struct {
	Component string   `json:"component"`
	Reports   []Report `json:"reports"`
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
	if env.Config().GetString("Mode") == env.OnlyManifestMode {
		return
	}

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

	for _, node := range r.ClusterReport.Nodes {
		existingReports.ClusterReport.Nodes = append(existingReports.ClusterReport.Nodes, node)
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

	for _, users := range r.ClusterReport.RBACReport.Users {
		existingReports.ClusterReport.RBACReport.Users = append(existingReports.ClusterReport.RBACReport.Users, users)
	}

	for _, groups := range r.ClusterReport.RBACReport.Groups {
		existingReports.ClusterReport.RBACReport.Groups = append(existingReports.ClusterReport.RBACReport.Groups, groups)
	}

	for _, roles := range r.ClusterReport.RBACReport.Roles {
		existingReports.ClusterReport.RBACReport.Roles = append(existingReports.ClusterReport.RBACReport.Roles, roles)
	}

	for _, clusterRoles := range r.ClusterReport.RBACReport.ClusterRoles {
		existingReports.ClusterReport.RBACReport.ClusterRoles = append(existingReports.ClusterReport.RBACReport.ClusterRoles, clusterRoles)
	}

	for _, clusterRoles := range r.ClusterReport.RBACReport.ClusterRoleBinding {
		existingReports.ClusterReport.RBACReport.ClusterRoleBinding = append(existingReports.ClusterReport.RBACReport.ClusterRoleBinding, clusterRoles)
	}

	for _, scc := range r.ClusterReport.RBACReport.SecurityContextConstraints {
		existingReports.ClusterReport.RBACReport.SecurityContextConstraints = append(existingReports.ClusterReport.RBACReport.SecurityContextConstraints, scc)
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

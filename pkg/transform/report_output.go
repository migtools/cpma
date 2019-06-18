package transform

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/clusterreport"
	"github.com/sirupsen/logrus"
)

// ReportOutput holds a collection of reports to be written to fil
type ReportOutput struct {
	Component string   `json:"component"`
	Reports   []Report `json:"reports"`
}

// ClusterOutput represents report of k8s resources
type ClusterOutput struct {
	ClusterReport *clusterreport.ClusterReport `json:"cluster"`
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
	var existingReports []ReportOutput

	jsonFile := "report.json"

	jsonData, err := io.ReadFile(jsonFile)
	if err != nil {
		logrus.Errorf("unable to read to report file: %s", jsonFile)
	}

	err = json.Unmarshal(jsonData, &existingReports)
	if err != nil {
		logrus.Errorf("unable to unmarshal existing report json")
	}

	existingReports = append(existingReports, r)

	jsonReports, err := json.MarshalIndent(existingReports, "", " ")
	if err != nil {
		logrus.Errorf("unable to marshal reports")
	}

	err = io.WriteFile(jsonReports, jsonFile)
	if err != nil {
		logrus.Errorf("unable to write to report file: %s", jsonFile)
	}
}

func (clusterOutput ClusterOutput) dumpToJSON() error {
	clusterReport := clusterOutput.ClusterReport

	jsonFile := filepath.Join(env.Config().GetString("OutputDir"), "cluster-report.json")

	file, err := json.MarshalIndent(&clusterReport, "", " ")
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(jsonFile, file, 0644); err != nil {
		return err
	}

	logrus.Debugf("Cluster report added to %s", jsonFile)
	return nil
}

// Flush reports to files
func (clusterOutput ClusterOutput) Flush() error {
	err := clusterOutput.dumpToJSON()
	if err != nil {
		return err
	}

	return nil
}

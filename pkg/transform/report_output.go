package transform

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/fusor/cpma/pkg/env"
	"github.com/sirupsen/logrus"
)

// ReportOutput holds a collection of reports to be written to fil
type ReportOutput struct {
	Component string   `json:"component"`
	Reports   []Report `json:"reports"`
}

// ReportOutputFlush flush reports to disk
var reportOutputFlush = func(r ReportOutput) error {
	logrus.Info("Flushing reports to disk")
	DumpReports(r)
	return nil
}

// Flush reports to files
func (r ReportOutput) Flush() error {
	return reportOutputFlush(r)
}

// DumpReports creates OCDs files
func DumpReports(r ReportOutput) {
	var existingReports []ReportOutput
	jsonfile := filepath.Join(env.Config().GetString("OutputDir"), "report.json")

	jsonData, err := ioutil.ReadFile(jsonfile)
	if err != nil {
		logrus.Errorf("unable to read to report file: %s", jsonfile)
	}

	json.Unmarshal(jsonData, &existingReports)
	existingReports = append(existingReports, r)

	jsonReports, err := json.Marshal(existingReports)
	if err != nil {
		logrus.Errorf("unable to marshal reports")
	}

	err = ioutil.WriteFile(jsonfile, jsonReports, 0644)
	if err != nil {
		logrus.Errorf("unable to write to report file: %s", jsonfile)
	}
}

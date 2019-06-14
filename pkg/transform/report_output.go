package transform

import (
	"encoding/json"

	"github.com/fusor/cpma/pkg/io"
	"github.com/sirupsen/logrus"
)

// ReportOutput holds a collection of reports to be written to fil
type ReportOutput struct {
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

package transform

import (
	"encoding/json"
	"os"
	"path"
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

	datafile := filepath.Join(env.Config().GetString("OutputDir"), "report.txt")
	os.MkdirAll(path.Dir(datafile), 0755)

	jsonReports, err := json.Marshal(r)
	if err != nil {
		logrus.Errorf("unable to marshal reports")
	}

	f, err := os.OpenFile(datafile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Errorf("unable to open report file: %s", datafile)
	}

	if _, err := f.Write(jsonReports); err != nil {
		logrus.Errorf("unable to open report file: %s", datafile)
	}
}

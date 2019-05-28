package transform

import (
	"encoding/json"
	"fmt"
	"os"
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

	jsonfile := filepath.Join(env.Config().GetString("OutputDir"), "report.json")
	htmlfile := filepath.Join(env.Config().GetString("OutputDir"), "report.html")

	jsonReports, err := json.Marshal(r)
	if err != nil {
		logrus.Errorf("unable to marshal reports")
	}

	jf, err := os.OpenFile(jsonfile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Errorf("unable to open report file: %s", jsonfile)
	}
	defer jf.Close()

	hf, err := os.OpenFile(htmlfile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Errorf("unable to open report file: %s", htmlfile)
	}
	defer hf.Close()

	_, err = jf.Write(jsonReports)
	if err != nil {
		logrus.Errorf("unable to write to report file: %s", jsonfile)
	}

	_, err = jf.Write([]byte(","))
	if err != nil {
		logrus.Errorf("unable to write to report file: %s", jsonfile)
	}

	for _, report := range r.Reports {
		var bgcolor string
		switch report.Confidence {
		case "red":
			bgcolor = "#FF0000"
		default:
			bgcolor = "#00FF00"
		}
		hf.Write([]byte(fmt.Sprintf("<tr bgcolor=%s><td>%s</td><td>%s</td><td>%s</td><td>%v</td></tr>\n", bgcolor, r.Component, report.Name, report.Kind, report.Supported)))
	}
}

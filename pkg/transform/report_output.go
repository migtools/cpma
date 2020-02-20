package transform

import (
	"github.com/konveyor/cpma/pkg/transform/reportoutput"
	"github.com/sirupsen/logrus"
)

// Report represents structure for final output
type Report struct {
	Report reportoutput.ReportOutput
}

// Flush reports to files
func (r Report) Flush() error {
	return ReportOutputFlush(r)
}

// ReportOutputFlush flush reports to disk
var ReportOutputFlush = func(r Report) error {
	logrus.Info("Flushing reports to disk")
	reportoutput.DumpReports(r.Report)
	return nil
}

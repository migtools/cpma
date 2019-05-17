package transform

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/fusor/cpma/pkg/env"
	"github.com/sirupsen/logrus"
)

// ReportOutput holds a collection of manifests to be written to fil
type ReportOutput struct {
	Manifests []Manifest
}

var reportOutputFlush = func(manifests []Manifest) error {
	logrus.Info("Writing report to disk")
	DumpData(manifests)
	return nil
}

// Flush manifests to files
func (m ReportOutput) Flush() error {
	return reportOutputFlush(m.Manifests)
}

func DumpData(dataList []Manifest) {
	for _, data := range dataList {
		datafile := filepath.Join(env.Config().GetString("OutputDir"), "report.txt")
		os.MkdirAll(path.Dir(datafile), 0755)
		err := ioutil.WriteFile(datafile, data.CRD, 0644)
		logrus.Printf("File:Added: %s", datafile)
		if err != nil {
			logrus.Panic(err)
		}
	}
}

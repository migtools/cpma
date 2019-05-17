package etl

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/fusor/cpma/pkg/env"
	"github.com/sirupsen/logrus"
)

// DataOutput holds a collection of data to be written to fil
type DataOutput struct {
	DataList []Data
}

// DataOutputFlush flushes data to disk
var DataOutputFlush = func(data []Data) error {
	logrus.Info("Flushing data to disk")
	DumpData(data)
	return nil
}

// Flush data to files
func (d DataOutput) Flush() error {
	return DataOutputFlush(d.DataList)
}

// DumpData creates OCDs files
func DumpData(dataList []Data) {
	for _, data := range dataList {
		datafile := filepath.Join(env.Config().GetString("OutputDir"), data.Type, data.Name)
		os.MkdirAll(path.Dir(datafile), 0755)
		err := ioutil.WriteFile(datafile, data.File, 0644)
		logrus.Printf("File:Added: %s", datafile)
		if err != nil {
			logrus.Panic(err)
		}
	}
}

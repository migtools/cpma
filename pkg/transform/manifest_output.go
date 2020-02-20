package transform

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/fusor/cpma/pkg/env"
	"github.com/sirupsen/logrus"
)

// ManifestOutput holds a collection of manifests to be written to fil
type ManifestOutput struct {
	Manifests []Manifest
}

// ManifestOutputFlush flush manifests to disk
var ManifestOutputFlush = func(manifests []Manifest) error {
	logrus.Info("Flushing manifests to disk")
	for _, manifest := range manifests {
		maniftestfile := filepath.Join(env.Config().GetString("WorkDir"), "manifests", manifest.Name)
		os.MkdirAll(path.Dir(maniftestfile), 0755)
		if err := ioutil.WriteFile(maniftestfile, manifest.CRD, 0644); err != nil {
			logrus.Panic(err)
		}
		logrus.Printf("CRD:Added: %s", maniftestfile)
	}
	return nil
}

// Flush manifests to files
func (m ManifestOutput) Flush() error {
	return ManifestOutputFlush(m.Manifests)
}

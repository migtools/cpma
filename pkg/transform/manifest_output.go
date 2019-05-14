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

// Flush calls DumpManifests to write file data
func (m ManifestOutput) Flush() error {
	logrus.Info("Writing file data:")
	DumpManifests(m.Manifests)
	return nil
}

// DumpManifests creates OCDs files
func DumpManifests(manifests []Manifest) {
	for _, manifest := range manifests {
		maniftestfile := filepath.Join(env.Config().GetString("OutputDir"), "manifests", manifest.Name)
		os.MkdirAll(path.Dir(maniftestfile), 0755)
		err := ioutil.WriteFile(maniftestfile, manifest.CRD, 0644)
		logrus.Printf("CRD:Added: %s", maniftestfile)
		if err != nil {
			logrus.Panic(err)
		}
	}
}

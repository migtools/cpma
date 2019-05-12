package transform

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/fusor/cpma/env"
	"github.com/sirupsen/logrus"
)

type ManifestTransformOutput struct {
	Manifests []Manifest
}

func (m ManifestTransformOutput) Flush() error {
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

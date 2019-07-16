package transform

import (
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io/remotehost"
	"github.com/sirupsen/logrus"
)

// DockerComponentName is the Docker component string
const DockerComponentName = "Docker"

// DockerExtraction holds Docker data extracted from OCP3
type DockerExtraction struct {
}

// DockerTransform is an Docker specific transform
type DockerTransform struct {
}

// Transform converts data collected from an OCP3 into a useful output
func (e DockerExtraction) Transform() ([]Output, error) {
	logrus.Info("DockerTransform::Transform")
	reports, err := e.buildReportOutput()
	if err != nil {
		return nil, err
	}
	outputs := []Output{reports}
	return outputs, nil
}

func (e DockerExtraction) buildReportOutput() (Output, error) {
	if env.Config().GetString("Mode") == env.OnlyManifestMode {
		logrus.Debug("Skipping DockerTransform report, only manifests mode was set")
		return ReportOutput{}, nil
	}

	componentReport := ComponentReport{
		Component: DockerComponentName,
	}

	componentReport.Reports = append(componentReport.Reports,
		Report{
			Name:       "Docker",
			Kind:       "Container Runtime",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "The Docker runtime has been replaced with CRI-O",
		})

	reportOutput := ReportOutput{
		ComponentReports: []ComponentReport{componentReport},
	}

	return reportOutput, nil
}

// Extract collects Docker configuration from an OCP3 cluster
func (e DockerTransform) Extract() (Extraction, error) {
	logrus.Info("DockerTransform::Extract")
	// Testing remote connection
	if env.Config().GetBool("FetchFromRemote") {
		_, err := remotehost.NewSSHSession(env.Config().GetString("Hostname"))
		if err != nil {
			return nil, err
		}
	}
	var extraction DockerExtraction
	return extraction, nil
}

// Validate confirms we have recieved good Docker configuration data during Extract
func (e DockerExtraction) Validate() error {
	return nil
}

// Name returns a human readable name for the transform
func (e DockerTransform) Name() string {
	return DockerComponentName
}

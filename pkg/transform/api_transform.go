package transform

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fusor/cpma/pkg/decode"
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/sirupsen/logrus"
)

// APIComponentName is the API component string
const APIComponentName = "API"

// APIExtraction holds API data extracted from OCP3
type APIExtraction struct {
	HTTPServingInfo ServingInfo
}

// ServingInfo contains information to serve a service
type ServingInfo struct {
	BindAddress string
}

// APITransform is an API specific transform
type APITransform struct {
}

// Transform converts data collected from an OCP3 into a useful output
func (e APIExtraction) Transform() ([]Output, error) {
	logrus.Info("APITransform::Transform")
	reports, err := e.buildReportOutput()
	if err != nil {
		return nil, err
	}
	outputs := []Output{reports}
	return outputs, nil
}

func (e APIExtraction) buildReportOutput() (Output, error) {
	componentReport := ComponentReport{
		Component: CrioComponentName,
	}

	portArray := strings.Split(e.HTTPServingInfo.BindAddress, ":")
	port := portArray[len(portArray)-1]

	var confidence = NoConfidence
	if port == "6443" {
		confidence = HighConfidence
	}

	componentReport.Reports = append(componentReport.Reports,
		Report{
			Name:       "API",
			Kind:       "Port",
			Supported:  false,
			Confidence: confidence,
			Comment:    fmt.Sprintf("The API Port for Openshift 4 is 6443 and is non-configurable. Your OCP 3 cluster is currently configured to use port %v", port),
		})

	reportOutput := ReportOutput{
		ComponentReports: []ComponentReport{componentReport},
	}

	return reportOutput, nil
}

// Extract collects API configuration from an OCP3 cluster
func (e APITransform) Extract() (Extraction, error) {
	logrus.Info("APITransform::Extract")
	content, err := io.FetchFile(env.Config().GetString("MasterConfigFile"))
	if err != nil {
		return nil, err
	}

	masterConfig, err := decode.MasterConfig(content)
	if err != nil {
		return nil, err
	}

	var extraction APIExtraction

	if masterConfig.ServingInfo.BindAddress != "" {
		extraction.HTTPServingInfo.BindAddress = masterConfig.ServingInfo.BindAddress
	}

	return extraction, nil
}

// Validate confirms we have recieved good API configuration data during Extract
func (e APIExtraction) Validate() error {
	if e.HTTPServingInfo.BindAddress == "" {
		return errors.New("could not determine API Port")
	}

	return nil
}

// Name returns a human readable name for the transform
func (e APITransform) Name() string {
	return APIComponentName
}

package transform

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fusor/cpma/pkg/decode"
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/apicert"
	"github.com/fusor/cpma/pkg/transform/reportoutput"
	"github.com/sirupsen/logrus"

	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
)

// APIComponentName is the API component string
const APIComponentName = "API"

// APIExtraction holds API data extracted from OCP3
type APIExtraction struct {
	ServingInfo legacyconfigv1.ServingInfo
}

// APITransform is an API specific transform
type APITransform struct {
}

// Transform converts data collected from an OCP3 into a useful output
func (e APIExtraction) Transform() ([]Output, error) {
	outputs := []Output{}

	if env.Config().GetBool("Manifests") {
		logrus.Info("APITransform::Transform:Manifests")
		manifests, err := e.buildManifestOutput()
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, manifests)
	}

	if env.Config().GetBool("Reporting") {
		logrus.Info("APITransform::Transform:Reports")
		e.buildReportOutput()
	}

	return outputs, nil
}

func (e APIExtraction) buildManifestOutput() (Output, error) {
	var manifests []Manifest

	APISecretCR, err := apicert.Translate(e.ServingInfo)
	if err != nil {
		return nil, err
	}

	if APISecretCR == nil {
		return nil, nil
	}

	APISecretCRYAML, err := GenYAML(APISecretCR)
	if err != nil {
		return nil, err
	}

	manifest := Manifest{Name: "100_CPMA-cluster-config-APISecret.yaml", CRD: APISecretCRYAML}
	manifests = append(manifests, manifest)

	return ManifestOutput{
		Manifests: manifests,
	}, nil
}

func (e APIExtraction) buildReportOutput() {
	componentReport := reportoutput.ComponentReport{
		Component: APIComponentName,
	}

	portArray := strings.Split(e.ServingInfo.BindAddress, ":")
	port := portArray[len(portArray)-1]

	var confidence = NoConfidence
	if port == "6443" {
		confidence = HighConfidence
	}

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "API",
			Kind:       "Port",
			Supported:  false,
			Confidence: confidence,
			Comment:    fmt.Sprintf("The API Port for Openshift 4 is 6443 and is non-configurable. Your OCP 3 cluster is currently configured to use port %v", port),
		})

	FinalReportOutput.Report.ComponentReports = append(FinalReportOutput.Report.ComponentReports, componentReport)
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
		extraction.ServingInfo.BindAddress = masterConfig.ServingInfo.BindAddress
	}

	if masterConfig.ServingInfo.CertInfo.CertFile != "" && masterConfig.ServingInfo.CertInfo.KeyFile != "" {
		extraction.ServingInfo.CertInfo = masterConfig.ServingInfo.CertInfo
	}

	return extraction, nil
}

// Validate confirms we have recieved API configuration data during Extract
func (e APIExtraction) Validate() error {
	if e.ServingInfo.BindAddress == "" {
		return errors.New("could not determine API Port")
	}

	if e.ServingInfo.CertInfo.CertFile == "" && e.ServingInfo.CertInfo.KeyFile == "" {
		return errors.New("could not determine API Certificate")
	}

	return nil
}

// Name returns a human readable name for the transform
func (e APITransform) Name() string {
	return APIComponentName
}

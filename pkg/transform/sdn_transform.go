package transform

import (
	"fmt"

	"github.com/fusor/cpma/pkg/decode"
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/sdn"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/sirupsen/logrus"
)

// SDNComponentName is the SDN component string
const SDNComponentName = "SDN"

// SDNExtraction is an SDN specific extraction
type SDNExtraction struct {
	legacyconfigv1.MasterConfig
}

// SDNTransform is an SDN specific transform
type SDNTransform struct {
}

const clusterNetworkComment = `Networks must be configured during installation,
 hostSubnetLength was replaced with hostPrefix in OCP4, default value was set to 23`

// Transform converts data collected from an OCP3 into a useful output
func (e SDNExtraction) Transform() ([]Output, error) {
	logrus.Info("SDNTransform::Transform")
	manifests, err := e.buildManifestOutput()
	if err != nil {
		return nil, err
	}
	reports, err := e.buildReportOutput()
	if err != nil {
		return nil, err
	}
	outputs := []Output{manifests, reports}
	return outputs, nil
}

func (e SDNExtraction) buildManifestOutput() (Output, error) {
	var manifests []Manifest

	networkCR, err := sdn.Translate(e.MasterConfig)
	if err != nil {
		return nil, err
	}

	networkCRYAML, err := GenYAML(networkCR)
	if err != nil {
		return nil, err
	}

	manifest := Manifest{Name: "100_CPMA-cluster-config-sdn.yaml", CRD: networkCRYAML}
	manifests = append(manifests, manifest)

	return ManifestOutput{
		Manifests: manifests,
	}, nil
}

func (e SDNExtraction) buildReportOutput() (Output, error) {
	componentReport := ComponentReport{
		Component: SDNComponentName,
	}

	for _, n := range e.MasterConfig.NetworkConfig.ClusterNetworks {
		cidrComment := fmt.Sprintf("Networks must be configured during installation, it's possible to use %s", n.CIDR)
		componentReport.Reports = append(componentReport.Reports,
			Report{
				Name:       "CIDR",
				Kind:       "ClusterNetwork",
				Supported:  true,
				Confidence: ModerateConfidence,
				Comment:    cidrComment,
			})

		componentReport.Reports = append(componentReport.Reports,
			Report{
				Name:       "HostSubnetLength",
				Kind:       "ClusterNetwork",
				Supported:  false,
				Confidence: NoConfidence,
				Comment:    clusterNetworkComment,
			})
	}

	componentReport.Reports = append(componentReport.Reports,
		Report{
			Name:       e.MasterConfig.NetworkConfig.ServiceNetworkCIDR,
			Kind:       "ServiceNetwork",
			Supported:  true,
			Confidence: ModerateConfidence,
			Comment:    "Networks must be configured during installation",
		})

	componentReport.Reports = append(componentReport.Reports,
		Report{
			Name:       "",
			Kind:       "ExternalIPNetworkCIDRs",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Configuration of ExternalIPNetworkCIDRs is not supported in OCP4",
		})

	componentReport.Reports = append(componentReport.Reports,
		Report{
			Name:       "",
			Kind:       "IngressIPNetworkCIDR",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Translation of this configuration is not supported, refer to ingress operator configuration for more information",
		})

	reportOutput := ReportOutput{
		ComponentReports: []ComponentReport{componentReport},
	}

	return reportOutput, nil
}

// Extract collects SDN configuration information from an OCP3 cluster
func (e SDNTransform) Extract() (Extraction, error) {
	logrus.Info("SDNTransform::Extract")

	content, err := io.FetchFile(env.Config().GetString("MasterConfigFile"))
	if err != nil {
		return nil, err
	}

	masterConfig, err := decode.MasterConfig(content)
	if err != nil {
		return nil, err
	}

	var extraction SDNExtraction
	extraction.MasterConfig = *masterConfig

	return extraction, nil
}

// Validate the data extracted from the OCP3 cluster
func (e SDNExtraction) Validate() error {
	err := sdn.Validate(e.MasterConfig)
	if err != nil {
		return err
	}

	return nil
}

// Name returns a human readable name for the transform
func (e SDNTransform) Name() string {
	return SDNComponentName
}

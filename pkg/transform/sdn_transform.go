package transform

import (
	"strconv"

	"github.com/fusor/cpma/pkg/decode"
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/sdn"
	configv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/sirupsen/logrus"
)

// SDNExtraction is an SDN specific extraction
type SDNExtraction struct {
	configv1.MasterConfig
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

	networkCRYAML, err := sdn.GenYAML(networkCR)
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
	reportOutput := ReportOutput{
		Component: SDNComponentName,
	}

	for _, n := range e.MasterConfig.NetworkConfig.ClusterNetworks {
		reportOutput.Reports = append(reportOutput.Reports,
			Report{
				Name:       n.CIDR,
				Kind:       "ClusterNetwork",
				Supported:  true,
				Confidence: "yellow",
				Comment:    clusterNetworkComment,
			})
		reportOutput.Reports = append(reportOutput.Reports,
			Report{
				Name:       strconv.Itoa(int(n.HostSubnetLength)),
				Kind:       "ClusterNetwork",
				Supported:  false,
				Confidence: "red",
				Comment:    "hostSubnetLength is not supported in OCP4",
			})
	}

	reportOutput.Reports = append(reportOutput.Reports,
		Report{
			Name:       e.MasterConfig.NetworkConfig.ServiceNetworkCIDR,
			Kind:       "ServiceNetwork",
			Supported:  true,
			Confidence: "yellow",
			Comment:    "Networks must be configured during installation",
		})

	for _, externalCIDR := range e.MasterConfig.NetworkConfig.ExternalIPNetworkCIDRs {
		reportOutput.Reports = append(reportOutput.Reports,
			Report{
				Name:       externalCIDR,
				Kind:       "ExternalIPNetworkCIDRs",
				Supported:  false,
				Confidence: "red",
				Comment:    "Configuration of ExternalIPNetworkCIDRs is not supported in OCP4",
			})
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
	return "SDN"
}

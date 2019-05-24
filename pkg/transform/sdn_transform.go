package transform

import (
	"github.com/fusor/cpma/pkg/config"
	"github.com/fusor/cpma/pkg/config/decode"
	"github.com/fusor/cpma/pkg/env"
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
	Config *config.Config
}

// Transform convers OCP3 data to configuration useful for OCP4
func (e SDNExtraction) Transform() (Output, error) {
	logrus.Info("SDNTransform::Transform")

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

// Extract collects SDN configuration information from an OCP3 cluster
func (e SDNTransform) Extract() (Extraction, error) {
	logrus.Info("SDNTransform::Extract")

	content, err := e.Config.Fetch(env.Config().GetString("MasterConfigFile"))
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

package transform

import (
	"encoding/json"
	"errors"

	"github.com/fusor/cpma/pkg/decode"
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/oauth"
	"github.com/sirupsen/logrus"
)

// OAuthExtraction holds OAuth data extracted from OCP3
type OAuthExtraction struct {
	IdentityProviders []oauth.IdentityProvider
}

// OAuthTransform is an OAuth specific transform
type OAuthTransform struct {
}

// Transform converts data collected from an OCP3 into a useful output
func (e OAuthExtraction) Transform() (Output, error) {
	logrus.Info("OAuthTransform::Transform")

	switch env.Config().Get("mode") {
	case ReportOutputType:
		return e.buildReportOutput()
	case ConvertOutputType:
		return e.buildManifestOutput()
	}
	return nil, errors.New("Unsupported Output Type")
}

func (e OAuthExtraction) buildManifestOutput() (Output, error) {
	var ocp4Cluster Cluster

	oauth, secrets, configMaps, err := oauth.Translate(e.IdentityProviders)
	if err != nil {
		return nil, errors.New("Unable to generate OAuth CRD")
	}

	ocp4Cluster.Master.OAuth = *oauth
	ocp4Cluster.Master.Secrets = secrets
	ocp4Cluster.Master.ConfigMaps = configMaps

	var manifests []Manifest
	if ocp4Cluster.Master.OAuth.Kind != "" {
		oauthCRD, err := GenYAML(ocp4Cluster.Master.OAuth)
		if err != nil {
			return nil, err
		}

		manifest := Manifest{Name: "100_CPMA-cluster-config-oauth.yaml", CRD: oauthCRD}
		manifests = append(manifests, manifest)

		for _, secret := range ocp4Cluster.Master.Secrets {
			secretCR, err := GenYAML(secret)
			if err != nil {
				return nil, err
			}

			filename := "100_CPMA-cluster-config-secret-" + secret.Metadata.Name + ".yaml"
			m := Manifest{Name: filename, CRD: secretCR}
			manifests = append(manifests, m)
		}

		for _, configMap := range ocp4Cluster.Master.ConfigMaps {
			configMapYAML, err := GenYAML(configMap)
			if err != nil {
				return nil, err
			}

			filename := "100_CPMA-cluster-config-configmap-" + configMap.Metadata.Name + ".yaml"
			m := Manifest{Name: filename, CRD: configMapYAML}
			manifests = append(manifests, m)
		}
	}

	return ManifestOutput{Manifests: manifests}, nil
}

func (e OAuthExtraction) buildReportOutput() (Output, error) {
	reportOutput := ReportOutput{
		Component: OAuthComponentName,
	}

	for _, p := range e.IdentityProviders {
		switch p.Kind {
		case "GitHubIdentityProvider",
			"GitLabIdentityProvider",
			"GoogleIdentityProvider",
			"HTPasswdPasswordIdentityProvider",
			"OpenIDIdentityProvider",
			"RequestHeaderIdentityProvider",
			"LDAPPasswordIdentityProvider",
			"KeystonePasswordIdentityProvider",
			"BasicAuthPasswordIdentityProvider":
			reportOutput.Reports = append(reportOutput.Reports,
				Report{
					Name:       p.Name,
					Kind:       p.Kind,
					Supported:  true,
					Confidence: "green",
				})
		default:
			reportOutput.Reports = append(reportOutput.Reports,
				Report{
					Name:       p.Name,
					Kind:       p.Kind,
					Supported:  false,
					Confidence: "red",
				})
		}
	}

	return reportOutput, nil
}

// Extract collects OAuth configuration from an OCP3 cluster
func (e OAuthTransform) Extract() (Extraction, error) {
	logrus.Info("OAuthTransform::Extract")
	content, err := io.FetchFile(env.Config().GetString("MasterConfigFile"))
	if err != nil {
		return nil, err
	}

	masterConfig, err := decode.MasterConfig(content)
	if err != nil {
		return nil, err
	}

	var extraction OAuthExtraction
	var htContent, caContent, crtContent, keyContent []byte

	if masterConfig.OAuthConfig != nil {
		for _, identityProvider := range masterConfig.OAuthConfig.IdentityProviders {

			providerJSON, err := identityProvider.Provider.MarshalJSON()
			if err != nil {
				return nil, err
			}

			provider := oauth.Provider{}
			err = json.Unmarshal(providerJSON, &provider)
			if err != nil {
				return nil, err
			}

			if provider.File != "" {
				htContent, err = io.FetchFile(provider.File)
				if err != nil {
					return nil, err
				}
			}
			if provider.CA != "" {
				caContent, err = io.FetchFile(provider.CA)
				if err != nil {
					return nil, err
				}
			}
			if provider.CertFile != "" {
				crtContent, err = io.FetchFile(provider.CertFile)
				if err != nil {
					return nil, err
				}
			}
			if provider.KeyFile != "" {
				keyContent, err = io.FetchFile(provider.KeyFile)
				if err != nil {
					return nil, err
				}
			}

			extraction.IdentityProviders = append(extraction.IdentityProviders,
				oauth.IdentityProvider{
					Kind:            provider.Kind,
					APIVersion:      provider.APIVersion,
					MappingMethod:   identityProvider.MappingMethod,
					Name:            identityProvider.Name,
					Provider:        identityProvider.Provider,
					HTFileName:      provider.File,
					HTFileData:      htContent,
					CAData:          caContent,
					CrtData:         crtContent,
					KeyData:         keyContent,
					UseAsChallenger: identityProvider.UseAsChallenger,
					UseAsLogin:      identityProvider.UseAsLogin,
				})
		}
	}

	return extraction, nil
}

// Validate confirms we have recieved good OAuth configuration data during Extract
func (e OAuthExtraction) Validate() error {
	err := oauth.Validate(e.IdentityProviders)
	if err != nil {
		return err
	}

	return nil
}

// Name returns a human readable name for the transform
func (e OAuthTransform) Name() string {
	return "OAuth"
}

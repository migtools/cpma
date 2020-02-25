package transform

import (
	"encoding/json"
	"fmt"

	"github.com/konveyor/cpma/pkg/decode"
	"github.com/konveyor/cpma/pkg/env"
	"github.com/konveyor/cpma/pkg/io"
	"github.com/konveyor/cpma/pkg/transform/oauth"
	"github.com/konveyor/cpma/pkg/transform/reportoutput"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// OAuthComponentName is the OAuth component string
const OAuthComponentName = "OAuth"

// OAuthExtraction holds OAuth data extracted from OCP3
type OAuthExtraction struct {
	IdentityProviders []oauth.IdentityProvider
	TokenConfig       oauth.TokenConfig
	Templates         legacyconfigv1.OAuthTemplates
}

// OAuthTransform is an OAuth specific transform
type OAuthTransform struct {
}

// Transform converts data collected from an OCP3 into a useful output
func (e OAuthExtraction) Transform() ([]Output, error) {
	outputs := []Output{}

	if env.Config().GetBool("Manifests") {
		logrus.Info("OAuthTransform::Transform:Manifests")
		manifests, err := e.buildManifestOutput()
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, manifests)
	}

	if env.Config().GetBool("Reporting") {
		logrus.Info("OAuthTransform::Transform:Reports")
		e.buildReportOutput()
	}

	return outputs, nil
}

func (e OAuthExtraction) buildManifestOutput() (Output, error) {
	var ocp4Cluster Cluster

	oauthResources, err := oauth.Translate(e.IdentityProviders, e.TokenConfig, e.Templates)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to generate OAuth CRD")
	}

	ocp4Cluster.Master.OAuth = *oauthResources.OAuthCRD
	ocp4Cluster.Master.Secrets = oauthResources.Secrets
	ocp4Cluster.Master.ConfigMaps = oauthResources.ConfigMaps

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

			filename := "100_CPMA-cluster-config-secret-" + secret.Name + ".yaml"
			m := Manifest{Name: filename, CRD: secretCR}
			manifests = append(manifests, m)
		}

		for _, configMap := range ocp4Cluster.Master.ConfigMaps {
			configMapYAML, err := GenYAML(configMap)
			if err != nil {
				return nil, err
			}

			filename := "100_CPMA-cluster-config-configmap-" + configMap.ObjectMeta.Name + ".yaml"
			m := Manifest{Name: filename, CRD: configMapYAML}
			manifests = append(manifests, m)
		}
	}

	return ManifestOutput{Manifests: manifests}, nil
}

func (e OAuthExtraction) buildReportOutput() {
	componentReport := reportoutput.ComponentReport{
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
			componentReport.Reports = append(componentReport.Reports,
				reportoutput.Report{
					Name:       p.Kind,
					Kind:       "IdentityProviders",
					Supported:  true,
					Confidence: HighConfidence,
					Comment:    fmt.Sprintf("Identity provider %s is supported in OCP4", p.Name),
				})
			if p.Kind == "OpenIDIdentityProvider" {
				componentReport.Reports = append(componentReport.Reports,
					reportoutput.Report{
						Name:       "Issuer",
						Kind:       fmt.Sprintf("IdentityProviders:%s", p.Kind),
						Supported:  true,
						Confidence: HighConfidence,
						Comment:    "OCP4 requires an 'issuer' URL, please edit OAuth manifest file and configure this field",
					})
			}
		default:
			componentReport.Reports = append(componentReport.Reports,
				reportoutput.Report{
					Name:       p.Kind,
					Kind:       "IdentityProviders",
					Supported:  false,
					Confidence: NoConfidence,
					Comment:    fmt.Sprintf("Identity provider %s is not supported in OCP4", p.Name),
				})
		}
	}

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "AccessTokenMaxAgeSeconds",
			Kind:       "TokenConfig",
			Supported:  true,
			Confidence: HighConfidence,
		})

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "AuthorizeTokenMaxAgeSeconds",
			Kind:       "TokenConfig",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Translation of AuthorizeTokenMaxAgeSeconds is not supported, it's value is 5 minutes in OCP4",
		})

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "",
			Kind:       "AssetPublicURL",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Translation of AssetPublicURL is not supported",
		})

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "",
			Kind:       "MasterPublicURL",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Translation of MasterPublicURL is not supported",
		})

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "",
			Kind:       "MasterCA",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Translation of MasterCA is not supported",
		})

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "",
			Kind:       "MasterURL",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Translation of MasterURL is not supported",
		})

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "",
			Kind:       "GrantConfig",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Translation of GrantConfig is not supported",
		})

	componentReport.Reports = append(componentReport.Reports,
		reportoutput.Report{
			Name:       "",
			Kind:       "SessionConfig",
			Supported:  false,
			Confidence: NoConfidence,
			Comment:    "Translation of SessionConfig is not supported",
		})

	FinalReportOutput.Report.ComponentReports = append(FinalReportOutput.Report.ComponentReports, componentReport)
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

	// Parse extracted oAuth providers and fetch their configuration file contents
	var extraction OAuthExtraction
	var htContent, caContent, crtContent, keyContent []byte

	if masterConfig.OAuthConfig != nil {
		for _, identityProvider := range masterConfig.OAuthConfig.IdentityProviders {

			providerJSON, err := identityProvider.Provider.MarshalJSON()
			if err != nil {
				return nil, err
			}

			provider := oauth.Provider{}
			if err := json.Unmarshal(providerJSON, &provider); err != nil {
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
					Kind:          provider.Kind,
					APIVersion:    provider.APIVersion,
					MappingMethod: identityProvider.MappingMethod,
					Name:          identityProvider.Name,
					Provider:      identityProvider.Provider,
					HTFileName:    provider.File,
					HTFileData:    htContent,
					CAData:        caContent,
					CrtData:       crtContent,
					KeyData:       keyContent,
				})
		}
	}

	// Get extracted token config
	tokenConfig := masterConfig.OAuthConfig.TokenConfig
	extraction.TokenConfig = oauth.TokenConfig{
		AuthorizeTokenMaxAgeSeconds: tokenConfig.AuthorizeTokenMaxAgeSeconds,
		AccessTokenMaxAgeSeconds:    tokenConfig.AccessTokenMaxAgeSeconds,
	}

	// Get templates
	templates := masterConfig.OAuthConfig.Templates
	if templates != nil {
		extraction.Templates = *templates
	}

	return extraction, nil
}

// Validate confirms we have recieved good OAuth configuration data during Extract
func (e OAuthExtraction) Validate() error {
	if err := oauth.Validate(e.IdentityProviders); err != nil {
		return err
	}

	return nil
}

// Name returns a human readable name for the transform
func (e OAuthTransform) Name() string {
	return OAuthComponentName
}

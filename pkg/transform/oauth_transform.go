package transform

import (
	"encoding/json"
	"errors"

	"github.com/fusor/cpma/pkg/config"
	"github.com/fusor/cpma/pkg/config/decode"
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/transform/oauth"
	configv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/sirupsen/logrus"
)

// OAuthExtraction holds OAuth data extracted from OCP3
type OAuthExtraction struct {
	IdentityProviders []oauth.IdentityProvider
}

// OAuthTransform is an OAuth specific transform
type OAuthTransform struct {
	Config *config.Config
}

// Transform converts data collected from an OCP3 cluster to OCP4 CR's
func (e OAuthExtraction) Transform() (Output, error) {
	logrus.Info("OAuthTransform::Transform")

	var ocp4Cluster Cluster

	oauth, secrets, err := oauth.Translate(e.IdentityProviders)
	if err != nil {
		return nil, errors.New("Unable to generate OAuth CRD")
	}

	ocp4Cluster.Master.OAuth = *oauth
	ocp4Cluster.Master.Secrets = secrets

	var manifests []Manifest
	if ocp4Cluster.Master.OAuth.Kind != "" {
		oauthCRD, err := ocp4Cluster.Master.OAuth.GenYAML()
		if err != nil {
			return nil, err
		}

		manifest := Manifest{Name: "100_CPMA-cluster-config-oauth.yaml", CRD: oauthCRD}
		manifests = append(manifests, manifest)

		for _, secret := range ocp4Cluster.Master.Secrets {
			secretCR, err := secret.GenYAML()
			if err != nil {
				return nil, err
			}

			filename := "100_CPMA-cluster-config-secret-" + secret.Metadata.Name + ".yaml"
			m := Manifest{Name: filename, CRD: secretCR}
			manifests = append(manifests, m)
		}
	}

	return ManifestOutput{Manifests: manifests}, nil
}

// Extract collects OAuth configuration from an OCP3 cluster
func (e OAuthTransform) Extract() (Extraction, error) {
	logrus.Info("OAuthTransform::Extract")
	content, err := e.Config.Fetch(env.Config().GetString("MasterConfigFile"))
	if err != nil {
		return nil, err
	}

	var masterConfig *configv1.MasterConfig
	masterConfig, err = decode.MasterConfig(content)
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
				htContent, err = e.Config.Fetch(provider.File)
				if err != nil {
					return nil, err
				}
			}
			if provider.CA != "" {
				caContent, err = e.Config.Fetch(provider.CA)
				if err != nil {
					return nil, err
				}
			}
			if provider.CertFile != "" {
				crtContent, err = e.Config.Fetch(provider.CertFile)
				if err != nil {
					return nil, err
				}
			}
			if provider.KeyFile != "" {
				keyContent, err = e.Config.Fetch(provider.KeyFile)
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
	logrus.Warn("Oauth Transform Validation Not Implmeneted")
	return nil // Simulate fine
}

// Name returns a human readable name for the transform
func (e OAuthTransform) Name() string {
	return "OAuth"
}

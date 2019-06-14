package image

import (
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
)

// Translate ImagePolicyConfig definitions
func Translate(imageCR *configv1.Image, imagePolicyConfig legacyconfigv1.ImagePolicyConfig) error {
	if imagePolicyConfig.AllowedRegistriesForImport != nil {
		for _, r := range *imagePolicyConfig.AllowedRegistriesForImport {
			imageCR.Spec.AllowedRegistriesForImport = append(imageCR.Spec.AllowedRegistriesForImport,
				configv1.RegistryLocation{
					DomainName: r.DomainName,
					Insecure:   r.Insecure,
				})
		}
	}

	if imagePolicyConfig.AdditionalTrustedCA != "" {
		imageCR.Spec.AdditionalTrustedCA.Name = imagePolicyConfig.AdditionalTrustedCA
	}

	return nil
}

// Validate registry data collected from an OCP3 cluster
func Validate(masterConfig legacyconfigv1.MasterConfig) error {
	return nil
}

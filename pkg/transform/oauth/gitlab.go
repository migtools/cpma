package oauth

import (
	"encoding/base64"

	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/configmaps"
	"github.com/fusor/cpma/pkg/transform/secrets"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func buildGitLabIP(serializer *json.Serializer, p IdentityProvider) (*configv1.IdentityProvider, *secrets.Secret, *configmaps.ConfigMap, error) {
	var (
		err         error
		idP         = &configv1.IdentityProvider{}
		secret      *secrets.Secret
		caConfigmap *configmaps.ConfigMap
		gitlab      legacyconfigv1.GitLabIdentityProvider
	)

	if _, _, err = serializer.Decode(p.Provider.Raw, nil, &gitlab); err != nil {
		return nil, nil, nil, errors.Wrap(err, "Failed to decode gitlab, see error")
	}

	idP.Type = "GitLab"
	idP.Name = p.Name
	idP.MappingMethod = configv1.MappingMethodType(p.MappingMethod)
	idP.GitLab.URL = gitlab.URL
	idP.GitLab.ClientID = gitlab.ClientID

	if gitlab.CA != "" {
		caConfigmap = configmaps.GenConfigMap("gitlab-configmap", OAuthNamespace, p.CAData)
		idP.GitLab.CA = configv1.ConfigMapNameReference{Name: caConfigmap.Metadata.Name}
	}

	secretName := p.Name + "-secret"
	idP.GitLab.ClientSecret.Name = secretName
	secretContent, err := io.FetchStringSource(gitlab.ClientSecret)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "Failed to fetch client secret for gitlab, see error")
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(secretContent))
	if secret, err = secrets.GenSecret(secretName, encoded, OAuthNamespace, secrets.LiteralSecretType); err != nil {
		return nil, nil, nil, errors.Wrap(err, "Failed to generate secret for gitlab, see error")
	}

	return idP, secret, caConfigmap, nil
}

func validateGitLabProvider(serializer *json.Serializer, p IdentityProvider) error {
	var gitlab legacyconfigv1.GitLabIdentityProvider

	if _, _, err := serializer.Decode(p.Provider.Raw, nil, &gitlab); err != nil {
		return errors.Wrap(err, "Failed to decode gitlab, see error")
	}

	if p.Name == "" {
		return errors.New("Name can't be empty")
	}

	if err := validateMappingMethod(p.MappingMethod); err != nil {
		return err
	}

	if gitlab.URL == "" {
		return errors.New("URL can't be empty")
	}

	if gitlab.ClientSecret.KeyFile != "" {
		return errors.New("Usage of encrypted files as secret value is not supported")
	}

	if err := validateClientData(gitlab.ClientID, gitlab.ClientSecret); err != nil {
		return err
	}

	return nil
}

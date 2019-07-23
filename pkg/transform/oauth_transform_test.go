package transform_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform"
	cpmatest "github.com/fusor/cpma/pkg/transform/internal/test"
	"github.com/fusor/cpma/pkg/transform/oauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuthExtractionTransform(t *testing.T) {
	env.Config().Set("Manifests", true)
	env.Config().Set("Reporting", true)

	var expectedManifests []transform.Manifest

	expectedOAuthCRYAML, err := ioutil.ReadFile("testdata/expected-CR-oauth.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-oauth.yaml", CRD: expectedOAuthCRYAML})

	expectedSecretBasicAuthProviderClientCertCRYAML, err := ioutil.ReadFile("testdata/expected-CR-secret-basicauth-client-cert-secret.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-basicauth-client-cert-secret.yaml", CRD: expectedSecretBasicAuthProviderClientCertCRYAML})

	expectedSecretBasicAuthProviderClientKeyCRYAML, err := ioutil.ReadFile("testdata/expected-CR-secret-basicauth-client-key-secret.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-basicauth-client-key-secret.yaml", CRD: expectedSecretBasicAuthProviderClientKeyCRYAML})

	expectedSecretGithubProvider, err := ioutil.ReadFile("testdata/expected-CR-secret-github-secret.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-github-secret.yaml", CRD: expectedSecretGithubProvider})

	expectedSecretGitlabProvider, err := ioutil.ReadFile("testdata/expected-CR-secret-gitlab-secret.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-gitlab-secret.yaml", CRD: expectedSecretGitlabProvider})

	expectedSecretGoogleProvider, err := ioutil.ReadFile("testdata/expected-CR-secret-google.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-google-secret.yaml", CRD: expectedSecretGoogleProvider})

	expectedSecretHtpasswdProvider, err := ioutil.ReadFile("testdata/expected-CR-secret-htpasswd.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-htpasswd-secret.yaml", CRD: expectedSecretHtpasswdProvider})

	expectedSecretKeystoneProviderCert, err := ioutil.ReadFile("testdata/expected-CR-secret-keystone-client-cert.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-keystone-client-cert-secret.yaml", CRD: expectedSecretKeystoneProviderCert})

	expectedSecretKeystoneProviderKey, err := ioutil.ReadFile("testdata/expected-CR-secret-keystone-client-key.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-keystone-client-key-secret.yaml", CRD: expectedSecretKeystoneProviderKey})

	expectedSecretOpenidProvider, err := ioutil.ReadFile("testdata/expected-CR-secret-openid.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-openid-secret.yaml", CRD: expectedSecretOpenidProvider})

	expectedSecretTemplateLogin, err := ioutil.ReadFile("testdata/expected-CR-secret-templates-login-secret.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-templates-login-secret.yaml", CRD: expectedSecretTemplateLogin})

	expectedSecretTemplateError, err := ioutil.ReadFile("testdata/expected-CR-secret-templates-error-secret.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-templates-error-secret.yaml", CRD: expectedSecretTemplateError})

	expectedSecretTemplateSelect, err := ioutil.ReadFile("testdata/expected-CR-secret-templates-providerselect-secret.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-templates-providerselect-secret.yaml", CRD: expectedSecretTemplateSelect})

	expectedConfigmapBasicauthProvider, err := ioutil.ReadFile("testdata/expected-CR-configmap-basicauth.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-configmap-basicauth-configmap.yaml", CRD: expectedConfigmapBasicauthProvider})

	expectedConfigmapGithubProvider, err := ioutil.ReadFile("testdata/expected-CR-configmap-github.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-configmap-github-configmap.yaml", CRD: expectedConfigmapGithubProvider})

	expectedConfigmapGitlabProvider, err := ioutil.ReadFile("testdata/expected-CR-configmap-gitlab.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-configmap-gitlab-configmap.yaml", CRD: expectedConfigmapGitlabProvider})

	expectedConfigmapKeystoneProvider, err := ioutil.ReadFile("testdata/expected-CR-configmap-keystone.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-configmap-keystone-configmap.yaml", CRD: expectedConfigmapKeystoneProvider})

	expectedConfigmapLDAPProvider, err := ioutil.ReadFile("testdata/expected-CR-configmap-ldap.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-configmap-ldap-configmap.yaml", CRD: expectedConfigmapLDAPProvider})

	expectedConfigmapRequestheader, err := ioutil.ReadFile("testdata/expected-CR-configmap-requestheader.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-configmap-requestheader-configmap.yaml", CRD: expectedConfigmapRequestheader})

	expectedReport := transform.ReportOutput{}
	jsonData, err := io.ReadFile("testdata/expected-report-oauth.json")
	require.NoError(t, err)

	err = json.Unmarshal(jsonData, &expectedReport)
	require.NoError(t, err)

	testCases := []struct {
		name              string
		expectedManifests []transform.Manifest
		expectedReports   transform.ReportOutput
	}{
		{
			name:              "transform registries extraction",
			expectedManifests: expectedManifests,
			expectedReports:   expectedReport,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualManifestsChan := make(chan []transform.Manifest)
			actualReportsChan := make(chan transform.ReportOutput)

			// Override flush method
			transform.ManifestOutputFlush = func(manifests []transform.Manifest) error {
				actualManifestsChan <- manifests
				return nil
			}
			transform.ReportOutputFlush = func(reports transform.ReportOutput) error {
				actualReportsChan <- reports
				return nil
			}

			identityProviders, templates, err := cpmatest.LoadIPTestData("testdata/master_config-bulk.yaml")
			require.NoError(t, err)

			testExtraction := transform.OAuthExtraction{
				IdentityProviders: identityProviders,
				TokenConfig: oauth.TokenConfig{
					AccessTokenMaxAgeSeconds:    int32(86400),
					AuthorizeTokenMaxAgeSeconds: int32(500),
				},
				Templates: *templates,
			}

			go func() {
				transformOutput, err := testExtraction.Transform()
				if err != nil {
					t.Error(err)
				}
				for _, output := range transformOutput {
					output.Flush()
				}
			}()

			actualManifests := <-actualManifestsChan
			assert.Equal(t, tc.expectedManifests, actualManifests)
			actualReports := <-actualReportsChan
			assert.Equal(t, tc.expectedReports, actualReports)
		})
	}
}

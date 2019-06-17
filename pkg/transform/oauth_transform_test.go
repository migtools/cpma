package transform_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform"
	"github.com/fusor/cpma/pkg/transform/oauth"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuthExtractionTransform(t *testing.T) {
	var expectedManifests []transform.Manifest

	expectedOAuthCRYAML, err := ioutil.ReadFile("testdata/expected-CR-oauth.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-oauth.yaml", CRD: expectedOAuthCRYAML})

	expectedSecretBasicAuthProviderClientCertCRYAML, err := ioutil.ReadFile("testdata/expected-CR-secret-my_remote_basic_auth_provider-client-cert-secret.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-my_remote_basic_auth_provider-client-cert-secret.yaml", CRD: expectedSecretBasicAuthProviderClientCertCRYAML})

	expectedSecretBasicAuthProviderClientKeyCRYAML, err := ioutil.ReadFile("testdata/expected-CR-secret-my_remote_basic_auth_provider-client-key-secret.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-my_remote_basic_auth_provider-client-key-secret.yaml", CRD: expectedSecretBasicAuthProviderClientKeyCRYAML})

	expectedSecretGithubProvider, err := ioutil.ReadFile("testdata/expected-CR-secret-github123456789-secret.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-github123456789-secret.yaml", CRD: expectedSecretGithubProvider})

	expectedSecretGitlabProvider, err := ioutil.ReadFile("testdata/expected-CR-secret-gitlab123456789-secret.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-gitlab123456789-secret.yaml", CRD: expectedSecretGitlabProvider})

	expectedSecretGoogleProvider, err := ioutil.ReadFile("testdata/expected-CR-secret-google123456789123456789.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-google123456789123456789-secret.yaml", CRD: expectedSecretGoogleProvider})

	expectedSecretHtpasswdProvider, err := ioutil.ReadFile("testdata/expected-CR-secret-htpasswd_auth.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-htpasswd_auth-secret.yaml", CRD: expectedSecretHtpasswdProvider})

	expectedSecretKeystoneProviderCert, err := ioutil.ReadFile("testdata/expected-CR-secret-keystone_provider-client-cert.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-my_keystone_provider-client-cert-secret.yaml", CRD: expectedSecretKeystoneProviderCert})

	expectedSecretKeystoneProviderKey, err := ioutil.ReadFile("testdata/expected-CR-secret-my_keystone_provider-client-key.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-my_keystone_provider-client-key-secret.yaml", CRD: expectedSecretKeystoneProviderKey})

	expectedSecretOpenidProvider, err := ioutil.ReadFile("testdata/expected-CR-secret-my_openid_connect.yaml")
	require.NoError(t, err)
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-my_openid_connect-secret.yaml", CRD: expectedSecretOpenidProvider})

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

			identityProviders, err := cpmatest.LoadIPTestData("testdata/bulk-test-master-config.yaml")
			require.NoError(t, err)

			testExtraction := transform.OAuthExtraction{
				IdentityProviders: identityProviders,
				TokenConfig: oauth.TokenConfig{
					AccessTokenMaxAgeSeconds:    int32(86400),
					AuthorizeTokenMaxAgeSeconds: int32(500),
				},
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
			assert.Equal(t, actualManifests, tc.expectedManifests)
			actualReports := <-actualReportsChan
			assert.Equal(t, actualReports, tc.expectedReports)
		})
	}
}

package transform_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform"
	"github.com/fusor/cpma/pkg/transform/configmaps"
	"github.com/fusor/cpma/pkg/transform/oauth"
	"github.com/fusor/cpma/pkg/transform/secrets"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuthExtractionTransform(t *testing.T) {
	var expectedManifests []transform.Manifest

	var expectedCrd oauth.CRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = oauth.OAuthNamespace

	var basicAuthIDP oauth.IdentityProviderBasicAuth
	basicAuthIDP.Type = "BasicAuth"
	basicAuthIDP.Challenge = true
	basicAuthIDP.Login = true
	basicAuthIDP.Name = "my_remote_basic_auth_provider"
	basicAuthIDP.MappingMethod = "claim"
	basicAuthIDP.BasicAuth.URL = "https://www.example.com/"
	basicAuthIDP.BasicAuth.TLSClientCert = &oauth.TLSClientCert{Name: "my_remote_basic_auth_provider-client-cert-secret"}
	basicAuthIDP.BasicAuth.TLSClientKey = &oauth.TLSClientKey{Name: "my_remote_basic_auth_provider-client-key-secret"}
	basicAuthIDP.BasicAuth.CA = &oauth.CA{Name: "basicauth-configmap"}

	var basicAuthCrtSecretCrd secrets.Secret
	basicAuthCrtSecretCrd.APIVersion = "v1"
	basicAuthCrtSecretCrd.Kind = "Secret"
	basicAuthCrtSecretCrd.Type = "Opaque"
	basicAuthCrtSecretCrd.Metadata.Namespace = oauth.OAuthNamespace
	basicAuthCrtSecretCrd.Metadata.Name = "my_remote_basic_auth_provider-client-cert-secret"
	basicAuthCrtSecretCrd.Data = secrets.BasicAuthFileSecret{BasicAuth: base64.StdEncoding.EncodeToString([]byte(""))}

	var basicAuthKeySecretCrd secrets.Secret
	basicAuthKeySecretCrd.APIVersion = "v1"
	basicAuthKeySecretCrd.Kind = "Secret"
	basicAuthKeySecretCrd.Type = "Opaque"
	basicAuthKeySecretCrd.Metadata.Namespace = oauth.OAuthNamespace
	basicAuthKeySecretCrd.Metadata.Name = "my_remote_basic_auth_provider-client-key-secret"
	basicAuthKeySecretCrd.Data = secrets.BasicAuthFileSecret{BasicAuth: base64.StdEncoding.EncodeToString([]byte(""))}

	var basicAuthConfigMap configmaps.ConfigMap
	basicAuthConfigMap.APIVersion = "v1"
	basicAuthConfigMap.Kind = "ConfigMap"
	basicAuthConfigMap.Metadata.Name = "basicauth-configmap"
	basicAuthConfigMap.Metadata.Namespace = oauth.OAuthNamespace
	basicAuthConfigMap.Data.CAData = ""

	var githubIDP oauth.IdentityProviderGitHub
	githubIDP.Type = "GitHub"
	githubIDP.Challenge = false
	githubIDP.Login = true
	githubIDP.MappingMethod = "claim"
	githubIDP.Name = "github123456789"
	githubIDP.GitHub.HostName = "test.example.com"
	githubIDP.GitHub.CA = &oauth.CA{Name: "github-configmap"}
	githubIDP.GitHub.ClientID = "2d85ea3f45d6777bffd7"
	githubIDP.GitHub.Organizations = []string{"myorganization1", "myorganization2"}
	githubIDP.GitHub.Teams = []string{"myorganization1/team-a", "myorganization2/team-b"}
	githubIDP.GitHub.ClientSecret.Name = "github123456789-secret"

	var githubSecretCrd secrets.Secret
	githubSecretCrd.APIVersion = "v1"
	githubSecretCrd.Kind = "Secret"
	githubSecretCrd.Type = "Opaque"
	githubSecretCrd.Metadata.Namespace = oauth.OAuthNamespace
	githubSecretCrd.Metadata.Name = "github123456789-secret"
	githubSecretCrd.Data = secrets.LiteralSecret{ClientSecret: base64.StdEncoding.EncodeToString([]byte("fake-secret"))}

	var githubConfigMap configmaps.ConfigMap
	githubConfigMap.APIVersion = "v1"
	githubConfigMap.Kind = "ConfigMap"
	githubConfigMap.Metadata.Name = "github-configmap"
	githubConfigMap.Metadata.Namespace = oauth.OAuthNamespace
	githubConfigMap.Data.CAData = ""

	var gitlabIDP oauth.IdentityProviderGitLab
	gitlabIDP.Name = "gitlab123456789"
	gitlabIDP.Type = "GitLab"
	gitlabIDP.Challenge = true
	gitlabIDP.Login = true
	gitlabIDP.MappingMethod = "claim"
	gitlabIDP.GitLab.URL = "https://gitlab.com/"
	gitlabIDP.GitLab.CA = &oauth.CA{Name: "gitlab-configmap"}
	gitlabIDP.GitLab.ClientID = "fake-id"
	gitlabIDP.GitLab.ClientSecret.Name = "gitlab123456789-secret"

	var gitlabSecretCrd secrets.Secret
	gitlabSecretCrd.APIVersion = "v1"
	gitlabSecretCrd.Kind = "Secret"
	gitlabSecretCrd.Type = "Opaque"
	gitlabSecretCrd.Metadata.Namespace = oauth.OAuthNamespace
	gitlabSecretCrd.Metadata.Name = "gitlab123456789-secret"
	gitlabSecretCrd.Data = secrets.LiteralSecret{ClientSecret: base64.StdEncoding.EncodeToString([]byte("fake-secret"))}

	var gitlabConfigMap configmaps.ConfigMap
	gitlabConfigMap.APIVersion = "v1"
	gitlabConfigMap.Kind = "ConfigMap"
	gitlabConfigMap.Metadata.Name = "gitlab-configmap"
	gitlabConfigMap.Metadata.Namespace = oauth.OAuthNamespace
	gitlabConfigMap.Data.CAData = ""

	var googleIDP oauth.IdentityProviderGoogle
	googleIDP.Type = "Google"
	googleIDP.Challenge = false
	googleIDP.Login = true
	googleIDP.MappingMethod = "claim"
	googleIDP.Name = "google123456789123456789"
	googleIDP.Google.ClientID = "82342890327-tf5lqn4eikdf4cb4edfm85jiqotvurpq.apps.googleusercontent.com"
	googleIDP.Google.ClientSecret.Name = "google123456789123456789-secret"
	googleIDP.Google.HostedDomain = "test.example.com"

	var googleSecretCrd secrets.Secret
	googleSecretCrd.APIVersion = "v1"
	googleSecretCrd.Kind = "Secret"
	googleSecretCrd.Type = "Opaque"
	googleSecretCrd.Metadata.Namespace = oauth.OAuthNamespace
	googleSecretCrd.Metadata.Name = "google123456789123456789-secret"
	googleSecretCrd.Data = secrets.LiteralSecret{ClientSecret: base64.StdEncoding.EncodeToString([]byte("fake-secret"))}

	var keystoneIDP oauth.IdentityProviderKeystone
	keystoneIDP.Type = "Keystone"
	keystoneIDP.Challenge = true
	keystoneIDP.Login = true
	keystoneIDP.Name = "my_keystone_provider"
	keystoneIDP.MappingMethod = "claim"
	keystoneIDP.Keystone.DomainName = "default"
	keystoneIDP.Keystone.URL = "http://fake.url:5000"
	keystoneIDP.Keystone.CA = &oauth.CA{Name: "keystone-configmap"}
	keystoneIDP.Keystone.TLSClientCert = &oauth.TLSClientCert{Name: "my_keystone_provider-client-cert-secret"}
	keystoneIDP.Keystone.TLSClientKey = &oauth.TLSClientKey{Name: "my_keystone_provider-client-key-secret"}

	var keystoneConfigMap configmaps.ConfigMap
	keystoneConfigMap.APIVersion = "v1"
	keystoneConfigMap.Kind = "ConfigMap"
	keystoneConfigMap.Metadata.Name = "keystone-configmap"
	keystoneConfigMap.Metadata.Namespace = oauth.OAuthNamespace
	keystoneConfigMap.Data.CAData = ""

	var htpasswdIDP oauth.IdentityProviderHTPasswd
	htpasswdIDP.Name = "htpasswd_auth"
	htpasswdIDP.Type = "HTPasswd"
	htpasswdIDP.Challenge = true
	htpasswdIDP.Login = true
	htpasswdIDP.MappingMethod = "claim"
	htpasswdIDP.HTPasswd.FileData.Name = "htpasswd_auth-secret"

	var htpasswdSecretCrd secrets.Secret
	htpasswdSecretCrd.APIVersion = "v1"
	htpasswdSecretCrd.Kind = "Secret"
	htpasswdSecretCrd.Type = "Opaque"
	htpasswdSecretCrd.Metadata.Namespace = oauth.OAuthNamespace
	htpasswdSecretCrd.Metadata.Name = "htpasswd_auth-secret"
	htpasswdSecretCrd.Data = secrets.HTPasswdFileSecret{HTPasswd: ""}

	var keystoneCrtSecretCrd secrets.Secret
	keystoneCrtSecretCrd.APIVersion = "v1"
	keystoneCrtSecretCrd.Kind = "Secret"
	keystoneCrtSecretCrd.Type = "Opaque"
	keystoneCrtSecretCrd.Metadata.Namespace = oauth.OAuthNamespace
	keystoneCrtSecretCrd.Metadata.Name = "my_keystone_provider-client-cert-secret"
	keystoneCrtSecretCrd.Data = secrets.KeystoneFileSecret{Keystone: ""}

	var keystoneKeySecretCrd secrets.Secret
	keystoneKeySecretCrd.APIVersion = "v1"
	keystoneKeySecretCrd.Kind = "Secret"
	keystoneKeySecretCrd.Type = "Opaque"
	keystoneKeySecretCrd.Metadata.Namespace = oauth.OAuthNamespace
	keystoneKeySecretCrd.Metadata.Name = "my_keystone_provider-client-key-secret"
	keystoneKeySecretCrd.Data = secrets.KeystoneFileSecret{Keystone: ""}

	var ldapIDP oauth.IdentityProviderLDAP
	ldapIDP.Name = "my_ldap_provider"
	ldapIDP.Type = "LDAP"
	ldapIDP.Challenge = true
	ldapIDP.Login = true
	ldapIDP.MappingMethod = "claim"
	ldapIDP.LDAP.Attributes.ID = []string{"dn"}
	ldapIDP.LDAP.Attributes.Email = []string{"mail"}
	ldapIDP.LDAP.Attributes.Name = []string{"cn"}
	ldapIDP.LDAP.Attributes.PreferredUsername = []string{"uid"}
	ldapIDP.LDAP.BindDN = "123"
	ldapIDP.LDAP.BindPassword = "321"
	ldapIDP.LDAP.CA = &oauth.CA{Name: "ldap-configmap"}
	ldapIDP.LDAP.Insecure = false
	ldapIDP.LDAP.URL = "ldap://ldap.example.com/ou=users,dc=acme,dc=com?uid"

	var ldapConfigMap configmaps.ConfigMap
	ldapConfigMap.APIVersion = "v1"
	ldapConfigMap.Kind = "ConfigMap"
	ldapConfigMap.Metadata.Name = "ldap-configmap"
	ldapConfigMap.Metadata.Namespace = oauth.OAuthNamespace
	ldapConfigMap.Data.CAData = ""

	var requestHeaderIDP oauth.IdentityProviderRequestHeader
	requestHeaderIDP.Type = "RequestHeader"
	requestHeaderIDP.Name = "my_request_header_provider"
	requestHeaderIDP.Challenge = true
	requestHeaderIDP.Login = true
	requestHeaderIDP.MappingMethod = "claim"
	requestHeaderIDP.RequestHeader.ChallengeURL = "https://example.com"
	requestHeaderIDP.RequestHeader.LoginURL = "https://example.com"
	requestHeaderIDP.RequestHeader.CA = &oauth.CA{Name: "requestheader-configmap"}
	requestHeaderIDP.RequestHeader.ClientCommonNames = []string{"my-auth-proxy"}
	requestHeaderIDP.RequestHeader.Headers = []string{"X-Remote-User", "SSO-User"}
	requestHeaderIDP.RequestHeader.EmailHeaders = []string{"X-Remote-User-Email"}
	requestHeaderIDP.RequestHeader.NameHeaders = []string{"X-Remote-User-Display-Name"}
	requestHeaderIDP.RequestHeader.PreferredUsernameHeaders = []string{"X-Remote-User-Login"}

	var requestheaderConfigMap configmaps.ConfigMap
	requestheaderConfigMap.APIVersion = "v1"
	requestheaderConfigMap.Kind = "ConfigMap"
	requestheaderConfigMap.Metadata.Name = "requestheader-configmap"
	requestheaderConfigMap.Metadata.Namespace = oauth.OAuthNamespace
	requestheaderConfigMap.Data.CAData = ""

	var openidIDP oauth.IdentityProviderOpenID
	openidIDP.Type = "OpenID"
	openidIDP.Challenge = false
	openidIDP.Login = true
	openidIDP.MappingMethod = "claim"
	openidIDP.Name = "my_openid_connect"
	openidIDP.OpenID.ClientID = "testid"
	openidIDP.OpenID.Claims.PreferredUsername = []string{"preferred_username", "email"}
	openidIDP.OpenID.Claims.Name = []string{"nickname", "given_name", "name"}
	openidIDP.OpenID.Claims.Email = []string{"custom_email_claim", "email"}
	openidIDP.OpenID.URLs.Authorize = "https://myidp.example.com/oauth2/authorize"
	openidIDP.OpenID.URLs.Token = "https://myidp.example.com/oauth2/token"
	openidIDP.OpenID.ClientSecret.Name = "my_openid_connect-secret"

	var openidSecretCrd secrets.Secret
	openidSecretCrd.APIVersion = "v1"
	openidSecretCrd.Kind = "Secret"
	openidSecretCrd.Type = "Opaque"
	openidSecretCrd.Metadata.Namespace = oauth.OAuthNamespace
	openidSecretCrd.Metadata.Name = "my_openid_connect-secret"
	openidSecretCrd.Data = secrets.LiteralSecret{ClientSecret: base64.StdEncoding.EncodeToString([]byte("testsecret"))}

	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, basicAuthIDP)
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, githubIDP)
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, gitlabIDP)
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, googleIDP)
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, htpasswdIDP)
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, keystoneIDP)
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, ldapIDP)
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, requestHeaderIDP)
	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, openidIDP)

	expectedCrd.Spec.TokenConfig.AccessTokenMaxAgeSeconds = int32(86400)

	expectedManifest, err := transform.GenYAML(expectedCrd)
	require.NoError(t, err)
	basicAuthCrtSecretManifest, err := transform.GenYAML(basicAuthCrtSecretCrd)
	require.NoError(t, err)
	basicAuthKeySecretManifest, err := transform.GenYAML(basicAuthKeySecretCrd)
	require.NoError(t, err)
	githubSecretManifest, err := transform.GenYAML(githubSecretCrd)
	require.NoError(t, err)
	gitlabSecretManifest, err := transform.GenYAML(gitlabSecretCrd)
	require.NoError(t, err)
	googleSecretManifest, err := transform.GenYAML(googleSecretCrd)
	require.NoError(t, err)
	htpasswdSecretManifest, err := transform.GenYAML(htpasswdSecretCrd)
	require.NoError(t, err)
	keystoneCrtSecretManifest, err := transform.GenYAML(keystoneCrtSecretCrd)
	require.NoError(t, err)
	keystoneKeySecretManifest, err := transform.GenYAML(keystoneKeySecretCrd)
	require.NoError(t, err)
	openidSecretManifest, err := transform.GenYAML(openidSecretCrd)
	require.NoError(t, err)

	basicAuthConfigMapManifest, err := transform.GenYAML(basicAuthConfigMap)
	require.NoError(t, err)
	githubConfigMapManifest, err := transform.GenYAML(githubConfigMap)
	require.NoError(t, err)
	gitlabConfigMapManifest, err := transform.GenYAML(gitlabConfigMap)
	require.NoError(t, err)
	keystoneConfigMapManifest, err := transform.GenYAML(keystoneConfigMap)
	require.NoError(t, err)
	ldapConfigMapManifest, err := transform.GenYAML(ldapConfigMap)
	require.NoError(t, err)
	requestheaderConfigMapManifest, err := transform.GenYAML(requestheaderConfigMap)
	require.NoError(t, err)

	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-oauth.yaml", CRD: expectedManifest})
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-my_remote_basic_auth_provider-client-cert-secret.yaml", CRD: basicAuthCrtSecretManifest})
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-my_remote_basic_auth_provider-client-key-secret.yaml", CRD: basicAuthKeySecretManifest})
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-github123456789-secret.yaml", CRD: githubSecretManifest})
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-gitlab123456789-secret.yaml", CRD: gitlabSecretManifest})
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-google123456789123456789-secret.yaml", CRD: googleSecretManifest})
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-htpasswd_auth-secret.yaml", CRD: htpasswdSecretManifest})
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-my_keystone_provider-client-cert-secret.yaml", CRD: keystoneCrtSecretManifest})
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-my_keystone_provider-client-key-secret.yaml", CRD: keystoneKeySecretManifest})
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-secret-my_openid_connect-secret.yaml", CRD: openidSecretManifest})
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-configmap-basicauth-configmap.yaml", CRD: basicAuthConfigMapManifest})
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-configmap-github-configmap.yaml", CRD: githubConfigMapManifest})
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-configmap-gitlab-configmap.yaml", CRD: gitlabConfigMapManifest})
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-configmap-keystone-configmap.yaml", CRD: keystoneConfigMapManifest})
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-configmap-ldap-configmap.yaml", CRD: ldapConfigMapManifest})
	expectedManifests = append(expectedManifests,
		transform.Manifest{Name: "100_CPMA-cluster-config-configmap-requestheader-configmap.yaml", CRD: requestheaderConfigMapManifest})

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

package oauth_test

import (
	"errors"
	"testing"

	"github.com/fusor/cpma/pkg/transform/oauth"
	cpmatest "github.com/fusor/cpma/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformMasterConfigKeystone(t *testing.T) {
	identityProviders, err := cpmatest.LoadIPTestData("testdata/keystone/test-master-config.yaml")
	require.NoError(t, err)

	var expectedCrd oauth.CRD
	expectedCrd.APIVersion = "config.openshift.io/v1"
	expectedCrd.Kind = "OAuth"
	expectedCrd.Metadata.Name = "cluster"
	expectedCrd.Metadata.NameSpace = oauth.OAuthNamespace

	var keystoneIDP = &oauth.IdentityProviderKeystone{}
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

	expectedCrd.Spec.IdentityProviders = append(expectedCrd.Spec.IdentityProviders, keystoneIDP)

	testCases := []struct {
		name        string
		expectedCrd *oauth.CRD
	}{
		{
			name:        "build keystone provider",
			expectedCrd: &expectedCrd,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resCrd, _, _, err := oauth.Translate(identityProviders)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedCrd, resCrd)
		})
	}
}

func TestKeystoneValidation(t *testing.T) {
	testCases := []struct {
		name         string
		requireError bool
		inputFile    string
		expectedErr  error
	}{
		{
			name:         "validate keystone provider",
			requireError: false,
			inputFile:    "testdata/keystone/test-master-config.yaml",
		},
		{
			name:         "fail on invalid name in keystone provider",
			requireError: true,
			inputFile:    "testdata/keystone/invalid-name-master-config.yaml",
			expectedErr:  errors.New("Name can't be empty"),
		},
		{
			name:         "fail on invalid mapping method in keystone provider",
			requireError: true,
			inputFile:    "testdata/keystone/invalid-mapping-master-config.yaml",
			expectedErr:  errors.New("Not valid mapping method"),
		},
		{
			name:         "fail on invalid url in keystone provider",
			requireError: true,
			inputFile:    "testdata/keystone/invalid-url-master-config.yaml",
			expectedErr:  errors.New("URL can't be empty"),
		},
		{
			name:         "fail on invalid domain name in keystone provider",
			requireError: true,
			inputFile:    "testdata/keystone/invalid-domainname-master-config.yaml",
			expectedErr:  errors.New("Domain name can't be empty"),
		},
		{
			name:         "fail on invalid key file in keystone provider",
			requireError: true,
			inputFile:    "testdata/keystone/invalid-keyfile-master-config.yaml",
			expectedErr:  errors.New("Key file can't be empty if cert file is specified"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			identityProvider, err := cpmatest.LoadIPTestData(tc.inputFile)
			require.NoError(t, err)

			err = oauth.Validate(identityProvider)

			if tc.requireError {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

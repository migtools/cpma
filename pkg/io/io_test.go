package io

import (
	"testing"

	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/io/remotehost"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Save before overriding
var _RunCMD = remotehost.RunCMD

// Overriding
func mockRunCMD(hostname, cmd string) (string, error) {
	return "remote value", nil
}

func TestFetchFile(t *testing.T) {
	testCases := []struct {
		name     string
		remote   bool
		expected []byte
		filename string
	}{
		{
			name:     "Fetch from remote",
			remote:   true,
			filename: "remote-file",
			expected: []byte{0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x20, 0x76, 0x61, 0x6c, 0x75, 0x65},
		},
		{
			name:     "Fetch from local",
			remote:   false,
			filename: "testdata/local-file",
			expected: []byte{0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x20, 0x76, 0x61, 0x6c, 0x75, 0x65, 0xa},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			env.Config().Set("FetchFromRemote", tc.remote)
			if tc.remote {
				defer func() { remotehost.RunCMD = _RunCMD }()
				remotehost.RunCMD = mockRunCMD

				f, err := FetchFile(tc.filename)
				require.NoError(t, err)
				assert.Equal(t, f, tc.expected)
			} else {
				f, err := FetchFile(tc.filename)
				require.NoError(t, err)
				assert.Equal(t, f, tc.expected)
			}
		})
	}
}

func TestFetchEnv(t *testing.T) {
	testCases := []struct {
		name     string
		host     string
		env      string
		expected string
	}{
		{
			name:     "Fetch remote ENV variable",
			host:     "remote.test.com",
			env:      "CPMA_TEST_ENV",
			expected: "remote value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() { remotehost.RunCMD = _RunCMD }()
			remotehost.RunCMD = mockRunCMD

			env, err := FetchEnv(tc.host, tc.env)
			require.NoError(t, err)
			assert.Equal(t, env, tc.expected)
		})
	}
}

func TestStringSource(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		expected string
	}{
		{
			name:     "Fetch String Source ENV variable",
			filename: "testdata/local-file",
			expected: "local value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			env.Config().Set("FetchFromRemote", false)
			stringSource := legacyconfigv1.StringSource{}
			stringSource.File = "testdata/local-file"
			f, err := FetchStringSource(stringSource)
			require.NoError(t, err)
			assert.Equal(t, f, tc.expected)
		})
	}
}

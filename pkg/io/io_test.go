package io

import (
	"os"
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
		expected string
		filename string
	}{
		{
			name:     "Fetch from remote",
			remote:   true,
			filename: "testdata/remote-file",
			expected: "remote value",
		},
		{
			name:     "Fetch from local",
			remote:   false,
			filename: "testdata/local-file",
			expected: "local value\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			env.Config().Set("FetchFromRemote", tc.remote)
			if tc.remote {
				defer func() { remotehost.RunCMD = _RunCMD }()
				defer func() {
					if err := os.Remove(tc.filename); err != nil {
						t.Logf("Cannot remove file: %s", tc.filename)
					}
				}()
				remotehost.RunCMD = mockRunCMD
			}
			f, err := FetchFile(tc.filename)
			require.NoError(t, err)
			assert.Equal(t, f, []byte(tc.expected))
		})
	}
}

func TestFetchEnv(t *testing.T) {
	testCases := []struct {
		name     string
		host     string
		env      string
		expected string
		remote   bool
	}{
		{
			name:     "Fetch remote ENV variable",
			host:     "remote.test.com",
			env:      "CPMA_TEST_ENV",
			expected: "remote value",
			remote:   true,
		},
		{
			name:     "Fetch local ENV variable",
			host:     "",
			env:      "CPMA_TEST_ENV",
			expected: "local value",
			remote:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.remote {
				env.Config().Set("FetchFromRemote", true)
			} else {
				env.Config().Set("FetchFromRemote", false)
				os.Setenv("CPMA_TEST_ENV", "local value")
			}

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
			name:     "Fetch String Source from value",
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

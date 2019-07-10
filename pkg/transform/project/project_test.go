package project

import (
	"encoding/json"
	"testing"

	"github.com/fusor/cpma/pkg/io"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTranslate(t *testing.T) {
	projectConfig := legacyconfigv1.ProjectConfig{
		DefaultNodeSelector:    "node-role.kubernetes.io/compute=true",
		ProjectRequestMessage:  "To provision Projects you must request access in https://labs.opentlc.com or https://rhpds.redhat.com",
		ProjectRequestTemplate: "default/project-request",
		SecurityAllocator: &legacyconfigv1.SecurityAllocator{
			UIDAllocatorRange:   "1000000000-1999999999/10000",
			MCSAllocatorRange:   "s0:/2",
			MCSLabelsPerProject: 5,
		},
	}

	f := "testdata/expected-project.json"
	content, err := io.ReadFile(f)
	if err != nil {
		t.Fatalf("Cannot read file: %s", f)
	}
	expected := &configv1.Project{}
	if err = json.Unmarshal(content, &expected); err != nil {
		t.Fatalf("Error Unmarshalling %s", f)
	}

	t.Run("Translate ProjectConfig", func(t *testing.T) {
		projectConfig, err := Translate(projectConfig)
		require.NoError(t, err)
		assert.Equal(t, projectConfig, expected)
	})
}

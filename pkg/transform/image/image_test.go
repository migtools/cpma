package image_test

import (
	"encoding/json"
	"testing"

	"github.com/fusor/cpma/pkg/io"
	"github.com/fusor/cpma/pkg/transform/image"
	configv1 "github.com/openshift/api/config/v1"
	legacyconfigv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestTranslate(t *testing.T) {
	t.Parallel()

	imagePolicyConfig := legacyconfigv1.ImagePolicyConfig{
		MaxImagesBulkImportedPerRepository:         0,
		DisableScheduledImport:                     false,
		ScheduledImageImportMinimumIntervalSeconds: 0,
		MaxScheduledImageImportsPerMinute:          0,
		AllowedRegistriesForImport: &legacyconfigv1.AllowedRegistries{
			{
				DomainName: "registry1.test.com",
				Insecure:   true,
			},
			{
				DomainName: "registry2.test.com",
				Insecure:   false,
			},
		},
		InternalRegistryHostname: "docker-registry.default.svc:5000",
		ExternalRegistryHostname: "external-registry.example.com",
		AdditionalTrustedCA:      "",
	}

	imageCR := &configv1.Image{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "config.openshift.io/v1",
			Kind:       "Image",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        "cluster",
			Annotations: map[string]string{"release.openshift.io/create-only": "true"},
		},
		Spec: configv1.ImageSpec{
			RegistrySources: configv1.RegistrySources{
				BlockedRegistries:  []string{"block.test.com"},
				InsecureRegistries: []string{"insecure.test.com"},
				AllowedRegistries:  []string{"allow1.test.com", "allow2.test.com"},
			},
		},
	}

	f := "testdata/expected-image.json"
	content, err := io.ReadFile(f)
	if err != nil {
		t.Fatalf("Cannot read file: %s", f)
	}
	expected := &configv1.Image{}
	if err = json.Unmarshal(content, &expected); err != nil {
		t.Fatalf("Error Unmarshalling %s", f)
	}

	t.Run("Translate ProjectConfig", func(t *testing.T) {
		err := image.Translate(imageCR, imagePolicyConfig)
		require.NoError(t, err)
		assert.Equal(t, imageCR, expected)
	})
}

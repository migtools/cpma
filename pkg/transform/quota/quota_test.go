package quota_test

import (
	"testing"

	cpmatest "github.com/fusor/cpma/pkg/transform/internal/test"
	"github.com/fusor/cpma/pkg/transform/quota"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	k8sapicore "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestTranslate(t *testing.T) {
	quotaList := cpmatest.CreateTestResourceQuotaList()

	configmaps := resource.Quantity{
		Format: resource.DecimalSI,
	}
	configmaps.Set(int64(10))

	t.Run("Translate Resource Quota", func(t *testing.T) {
		quota, err := quota.Translate(quotaList.Items[0])
		require.NoError(t, err)
		assert.Equal(t, "ResourceQuota", quota.TypeMeta.Kind)
		assert.Equal(t, "quota.openshift.io/v1", quota.TypeMeta.APIVersion)
		assert.Equal(t, "resourcequota1", quota.ObjectMeta.Name)
		assert.Equal(t, "namespacetest1", quota.ObjectMeta.Namespace)
		assert.Equal(t, "true", quota.ObjectMeta.Annotations["release.openshift.io/create-only"])
		assert.Equal(t, configmaps, quota.Spec.Hard["configmaps"])
		assert.Equal(t, k8sapicore.ResourceList(nil), quota.Status.Hard)
		assert.Equal(t, k8sapicore.ResourceList(nil), quota.Status.Used)
	})
}

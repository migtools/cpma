package clusterquota_test

import (
	"testing"

	"github.com/konveyor/cpma/pkg/transform/clusterquota"
	cpmatest "github.com/konveyor/cpma/pkg/transform/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	k8sapicore "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestTranslate(t *testing.T) {
	clusterQuotaList := cpmatest.CreateTestClusterQuotaList()

	testkey := resource.Quantity{
		Format: resource.DecimalSI,
	}
	testkey.Set(int64(99))

	t.Run("Translate Cluster Resource Quota", func(t *testing.T) {
		quota, err := clusterquota.Translate(clusterQuotaList.Items[0])
		require.NoError(t, err)
		assert.Equal(t, "ClusterResourceQuota", quota.TypeMeta.Kind)
		assert.Equal(t, "quota.openshift.io/v1", quota.TypeMeta.APIVersion)
		assert.Equal(t, "test-quota1", quota.ObjectMeta.Name)
		assert.Equal(t, "true", quota.ObjectMeta.Annotations["release.openshift.io/create-only"])
		assert.Equal(t, testkey, quota.Spec.Quota.Hard["testkey"])
		assert.Equal(t, k8sapicore.ResourceList(nil), quota.Status.Total.Hard)
		assert.Equal(t, k8sapicore.ResourceList(nil), quota.Status.Total.Used)
	})
}

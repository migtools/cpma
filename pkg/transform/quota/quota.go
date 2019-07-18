package quota

import (
	k8sapicore "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BuildManifest definitions
func BuildManifest(quota k8sapicore.ResourceQuota) (*k8sapicore.ResourceQuota, error) {
	const (
		annokey = "release.openshift.io/create-only"
		annoval = "true"
	)

	quota.APIVersion = "quota.openshift.io/v1"
	quota.Kind = "ResourceQuota"
	quota.ObjectMeta = metav1.ObjectMeta{
		Name:        quota.Name,
		Namespace:   quota.Namespace,
		Annotations: map[string]string{annokey: annoval},
	}
	quota.Status = k8sapicore.ResourceQuotaStatus{}

	return &quota, nil
}

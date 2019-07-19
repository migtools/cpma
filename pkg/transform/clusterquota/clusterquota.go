package clusterquota

import (
	o7tapiquota "github.com/openshift/api/quota/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Translate definitions
func Translate(quota o7tapiquota.ClusterResourceQuota) (*o7tapiquota.ClusterResourceQuota, error) {
	const (
		annokey = "release.openshift.io/create-only"
		annoval = "true"
	)

	quota.APIVersion = "quota.openshift.io/v1"
	quota.Kind = "ClusterResourceQuota"
	quota.ObjectMeta = metav1.ObjectMeta{
		Name:        quota.Name,
		Annotations: map[string]string{annokey: annoval},
	}
	quota.Status = o7tapiquota.ClusterResourceQuotaStatus{}

	return &quota, nil
}

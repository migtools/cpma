package api

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// InitK8SOrDie init k8s client or panic
func InitK8SOrDie(config *rest.Config) *kubernetes.Clientset {
	once.K8S.Do(func() {
		k8sClient, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic("Error in creating K8S API client")
		}
		instances.K8S = k8sClient
	})
	return instances.K8S
}

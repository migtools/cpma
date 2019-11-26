package api

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewDiscoveryOrDie init k8s client or panic
func NewDiscoveryOrDie(config *rest.Config) *discovery.DiscoveryClient {
	return discovery.NewDiscoveryClientForConfigOrDie(config)
}

// NewK8SOrDie init k8s client or panic
func NewK8SOrDie(config *rest.Config) *kubernetes.Clientset {
	return kubernetes.NewForConfigOrDie(config)
}

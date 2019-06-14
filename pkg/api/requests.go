package api

import (
	k8sapi "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var listOptions metav1.ListOptions

// ListNamespaces list all namespaces, wrapper around cliet-go
func ListNamespaces() (*k8sapi.NamespaceList, error) {
	return Client.CoreV1().Namespaces().List(listOptions)
}

// ListPods list all pods in namespace, wrapper around cliet-go
func ListPods(namespace string) (*k8sapi.PodList, error) {
	return Client.CoreV1().Pods(namespace).List(listOptions)
}

// ListPVs list all pods in namespace, wrapper around cliet-go
func ListPVs() (*k8sapi.PersistentVolumeList, error) {
	return Client.CoreV1().PersistentVolumes().List(listOptions)
}

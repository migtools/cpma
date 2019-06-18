package api

import (
	k8sapicore "k8s.io/api/core/v1"
	k8sapistorage "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Resources represent api resources used in report
type Resources struct {
	PersistentVolumeList *k8sapicore.PersistentVolumeList
	NodeList             *k8sapicore.NodeList
	StorageClassList     *k8sapistorage.StorageClassList
	NamespaceMap         map[string]*NamespaceResources
}

// NamespaceResources holds all resources that belong to a namespace
type NamespaceResources struct {
	PodList *k8sapicore.PodList
}

var listOptions metav1.ListOptions

// ListNamespaces list all namespaces, wrapper around client-go
func ListNamespaces() (*k8sapicore.NamespaceList, error) {
	return Client.CoreV1().Namespaces().List(listOptions)
}

// ListPods list all pods in namespace, wrapper around client-go
func ListPods(namespace string) (*k8sapicore.PodList, error) {
	return Client.CoreV1().Pods(namespace).List(listOptions)
}

// ListPVs list all PVs, wrapper around client-go
func ListPVs() (*k8sapicore.PersistentVolumeList, error) {
	return Client.CoreV1().PersistentVolumes().List(listOptions)
}

// ListNodes list all nodes, wrapper around client-go
func ListNodes() (*k8sapicore.NodeList, error) {
	return Client.CoreV1().Nodes().List(listOptions)
}

// ListStorageClasses list all storage classes, wrapper around client-go
func ListStorageClasses() (*k8sapistorage.StorageClassList, error) {
	return Client.StorageV1().StorageClasses().List(listOptions)
}

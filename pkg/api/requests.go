package api

import (
	O7tapiroute "github.com/openshift/api/route/v1"

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
	PodList   *k8sapicore.PodList
	RouteList *O7tapiroute.RouteList
}

var listOptions metav1.ListOptions

// ListNamespaces list all namespaces, wrapper around client-go
func ListNamespaces() (*k8sapicore.NamespaceList, error) {
	return K8sClient.CoreV1().Namespaces().List(listOptions)
}

// ListPods list all pods in namespace, wrapper around client-go
func ListPods(namespace string) (*k8sapicore.PodList, error) {
	return K8sClient.CoreV1().Pods(namespace).List(listOptions)
}

// ListPVs list all PVs, wrapper around client-go
func ListPVs() (*k8sapicore.PersistentVolumeList, error) {
	return K8sClient.CoreV1().PersistentVolumes().List(listOptions)
}

// ListNodes list all nodes, wrapper around client-go
func ListNodes() (*k8sapicore.NodeList, error) {
	return K8sClient.CoreV1().Nodes().List(listOptions)
}

// ListStorageClasses list all storage classes, wrapper around client-go
func ListStorageClasses() (*k8sapistorage.StorageClassList, error) {
	return K8sClient.StorageV1().StorageClasses().List(listOptions)
}

// ListRoutes list all storage classes, wrapper around client-go
func ListRoutes(namespace string) (*O7tapiroute.RouteList, error) {
	return O7tClient.routeClient.Routes(namespace).List(listOptions)
}

package api

import (
	O7tapiauth "github.com/openshift/api/authorization/v1"
	o7tapiroute "github.com/openshift/api/route/v1"
	O7tapisecurity "github.com/openshift/api/security/v1"
	O7tapiuser "github.com/openshift/api/user/v1"
	k8sapiapps "k8s.io/api/apps/v1"
	k8sapicore "k8s.io/api/core/v1"
	k8sapistorage "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Resources represent api resources used in report
type Resources struct {
	PersistentVolumeList *k8sapicore.PersistentVolumeList
	NodeList             *k8sapicore.NodeList
	StorageClassList     *k8sapistorage.StorageClassList
	NamespaceList        []NamespaceResources
	RBACResources        RBACResources
}

// RBACResources contains all resources related to RBAC report
type RBACResources struct {
	UsersList                      *O7tapiuser.UserList
	GroupsList                     *O7tapiuser.GroupList
	ClusterRolesList               *O7tapiauth.ClusterRoleList
	ClusterRolesBindingsList       *O7tapiauth.ClusterRoleBindingList
	SecurityContextConstraintsList *O7tapisecurity.SecurityContextConstraintsList
}

// NamespaceResources holds all resources that belong to a namespace
type NamespaceResources struct {
	NamespaceName  string
	PodList        *k8sapicore.PodList
	RouteList      *o7tapiroute.RouteList
	DaemonSetList  *k8sapiapps.DaemonSetList
	DeploymentList *k8sapiapps.DeploymentList
	RolesList      *O7tapiauth.RoleList
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
func ListRoutes(namespace string) (*o7tapiroute.RouteList, error) {
	return O7tClient.routeClient.Routes(namespace).List(listOptions)
}

// ListDeployments will list all deployments seeding in the selected namespace
func ListDeployments(namespace string) (*k8sapiapps.DeploymentList, error) {
	return K8sClient.AppsV1().Deployments(namespace).List(listOptions)
}

// ListDaemonSets will collect all DS from specific namespace
func ListDaemonSets(namespace string) (*k8sapiapps.DaemonSetList, error) {
	return K8sClient.AppsV1().DaemonSets(namespace).List(listOptions)
}

// ListUsers list all users, wrapper around client-go
func ListUsers() (*O7tapiuser.UserList, error) {
	return O7tClient.userClient.Users().List(listOptions)
}

// ListGroups list all users, wrapper around client-go
func ListGroups() (*O7tapiuser.GroupList, error) {
	return O7tClient.userClient.Groups().List(listOptions)
}

// ListRoles list all storage classes, wrapper around client-go
func ListRoles(namespace string) (*O7tapiauth.RoleList, error) {
	return O7tClient.authClient.Roles(namespace).List(listOptions)
}

// ListClusterRoles list all storage classes, wrapper around client-go
func ListClusterRoles() (*O7tapiauth.ClusterRoleList, error) {
	return O7tClient.authClient.ClusterRoles().List(listOptions)
}

// ListClusterRolesBindings list all storage classes, wrapper around client-go
func ListClusterRolesBindings() (*O7tapiauth.ClusterRoleBindingList, error) {
	return O7tClient.authClient.ClusterRoleBindings().List(listOptions)
}

// ListSCC list all security context constraints, wrapper around client-go
func ListSCC() (*O7tapisecurity.SecurityContextConstraintsList, error) {
	return O7tClient.securityClient.SecurityContextConstraints().List(listOptions)
}

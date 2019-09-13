package api

import (
	o7tapiauth "github.com/openshift/api/authorization/v1"
	o7tapiquota "github.com/openshift/api/quota/v1"
	o7tapiroute "github.com/openshift/api/route/v1"
	o7tapisecurity "github.com/openshift/api/security/v1"
	o7tapiuser "github.com/openshift/api/user/v1"
	"github.com/sirupsen/logrus"

	"k8s.io/api/apps/v1beta1"
	k8sapicore "k8s.io/api/core/v1"
	extv1b1 "k8s.io/api/extensions/v1beta1"
	k8sapistorage "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Resources represent api resources used in report
type Resources struct {
	QuotaList            *o7tapiquota.ClusterResourceQuotaList
	NodeList             *k8sapicore.NodeList
	PersistentVolumeList *k8sapicore.PersistentVolumeList
	StorageClassList     *k8sapistorage.StorageClassList
	NamespaceList        []NamespaceResources
	RBACResources        RBACResources
}

// RBACResources contains all resources related to RBAC report
type RBACResources struct {
	UsersList                      *o7tapiuser.UserList
	GroupList                      *o7tapiuser.GroupList
	ClusterRolesList               *o7tapiauth.ClusterRoleList
	ClusterRolesBindingsList       *o7tapiauth.ClusterRoleBindingList
	SecurityContextConstraintsList *o7tapisecurity.SecurityContextConstraintsList
}

// NamespaceResources holds all resources that belong to a namespace
type NamespaceResources struct {
	NamespaceName     string
	DaemonSetList     *extv1b1.DaemonSetList
	DeploymentList    *v1beta1.DeploymentList
	PodList           *k8sapicore.PodList
	ResourceQuotaList *k8sapicore.ResourceQuotaList
	RolesList         *o7tapiauth.RoleList
	RouteList         *o7tapiroute.RouteList
	PVCList           *k8sapicore.PersistentVolumeClaimList
}

var listOptions metav1.ListOptions

// ListNamespaces list all namespaces, wrapper around client-go
func ListNamespaces(ch chan<- *k8sapicore.NamespaceList) {
	namespaces, err := K8sClient.CoreV1().Namespaces().List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- namespaces
}

// ListPods list all pods in namespace, wrapper around client-go
func ListPods(namespace string, ch chan<- *k8sapicore.PodList) {
	pods, err := K8sClient.CoreV1().Pods(namespace).List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- pods
}

// ListPVs list all PVs, wrapper around client-go
func ListPVs(ch chan<- *k8sapicore.PersistentVolumeList) {
	pvs, err := K8sClient.CoreV1().PersistentVolumes().List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- pvs
}

// ListNodes list all nodes, wrapper around client-go
func ListNodes(ch chan<- *k8sapicore.NodeList) {
	nodes, err := K8sClient.CoreV1().Nodes().List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- nodes
}

// ListQuotas list all cluster quotas classes, wrapper around client-go
func ListQuotas(ch chan<- *o7tapiquota.ClusterResourceQuotaList) {
	quotas, err := O7tClient.quotaClient.ClusterResourceQuotas().List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- quotas
}

// ListResourceQuotas list all quotas classes, wrapper around client-go
func ListResourceQuotas(namespace string, ch chan<- *k8sapicore.ResourceQuotaList) {
	quotas, err := K8sClient.CoreV1().ResourceQuotas(namespace).List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- quotas
}

// ListRoutes list all routes classes, wrapper around client-go
func ListRoutes(namespace string, ch chan<- *o7tapiroute.RouteList) {
	routes, err := O7tClient.routeClient.Routes(namespace).List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- routes
}

// ListStorageClasses list all storage classes, wrapper around client-go
func ListStorageClasses(ch chan<- *k8sapistorage.StorageClassList) {
	sc, err := K8sClient.StorageV1().StorageClasses().List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- sc
}

// ListDeployments will list all deployments seeding in the selected namespace
func ListDeployments(namespace string, ch chan<- *v1beta1.DeploymentList) {
	deployments, err := K8sClient.AppsV1beta1().Deployments(namespace).List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- deployments
}

// ListDaemonSets will collect all DS from specific namespace
func ListDaemonSets(namespace string, ch chan<- *extv1b1.DaemonSetList) {
	daemonSets, err := K8sClient.ExtensionsV1beta1().DaemonSets(namespace).List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- daemonSets
}

// ListUsers list all users, wrapper around client-go
func ListUsers(ch chan<- *o7tapiuser.UserList) {
	users, err := O7tClient.userClient.Users().List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- users
}

// ListGroups list all users, wrapper around client-go
func ListGroups(ch chan<- *o7tapiuser.GroupList) {
	groups, err := O7tClient.userClient.Groups().List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- groups
}

// ListRoles list all storage classes, wrapper around client-go
func ListRoles(namespace string, ch chan<- *o7tapiauth.RoleList) {
	roles, err := O7tClient.authClient.Roles(namespace).List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- roles
}

// ListClusterRoles list all storage classes, wrapper around client-go
func ListClusterRoles(ch chan<- *o7tapiauth.ClusterRoleList) {
	clusterRoles, err := O7tClient.authClient.ClusterRoles().List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- clusterRoles
}

// ListClusterRolesBindings list all storage classes, wrapper around client-go
func ListClusterRolesBindings(ch chan<- *o7tapiauth.ClusterRoleBindingList) {
	clusterRolesBindings, err := O7tClient.authClient.ClusterRoleBindings().List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- clusterRolesBindings
}

// ListSCC list all security context constraints, wrapper around client-go
func ListSCC(ch chan<- *o7tapisecurity.SecurityContextConstraintsList) {
	scc, err := O7tClient.securityClient.SecurityContextConstraints().List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- scc
}

// ListPVCs list all PVs, wrapper around client-go
func ListPVCs(namespace string, ch chan<- *k8sapicore.PersistentVolumeClaimList) {
	pvcs, err := K8sClient.CoreV1().PersistentVolumeClaims(namespace).List(listOptions)
	if err != nil {
		logrus.Fatal(err)
	}
	ch <- pvcs
}

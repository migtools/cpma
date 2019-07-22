package cluster

import (
	"github.com/fusor/cpma/pkg/api"
	o7tapiauth "github.com/openshift/api/authorization/v1"
	o7tapiquota "github.com/openshift/api/quota/v1"
	o7tapiroute "github.com/openshift/api/route/v1"
	"github.com/sirupsen/logrus"

	k8sapiapps "k8s.io/api/apps/v1"
	k8sapicore "k8s.io/api/core/v1"
	k8scorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	k8sMeta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Report represents json report of k8s resources
type Report struct {
	Nodes          []NodeReport         `json:"nodes"`
	Quotas         []QuotaReport        `json:"quotas"`
	Namespaces     []NamespaceReport    `json:"namespaces,omitempty"`
	PVs            []PVReport           `json:"pvs,omitempty"`
	StorageClasses []StorageClassReport `json:"storageClasses,omitempty"`
	RBACReport     RBACReport           `json:"rbacreport,omitempty"`
}

// NodeReport represents json report of k8s nodes
type NodeReport struct {
	Name       string        `json:"name"`
	MasterNode bool          `json:"masterNode"`
	Resources  NodeResources `json:"resources"`
}

// NodeResources represents a json report of Node resources
type NodeResources struct {
	CPU            *resource.Quantity `json:"cpu"`
	MemoryConsumed *resource.Quantity `json:"memoryConsumed"`
	MemoryCapacity *resource.Quantity `json:"memoryCapacity"`
	RunningPods    *resource.Quantity `json:"runningPods"`
	PodCapacity    *resource.Quantity `json:"podCapacity"`
}

// NamespaceReport represents json report of k8s namespaces
type NamespaceReport struct {
	Name         string                   `json:"name"`
	LatestChange k8sMeta.Time             `json:"latestChange,omitempty"`
	Resources    ContainerResourcesReport `json:"resources,omitempty"`
	Pods         []PodReport              `json:"pods,omitempty"`
	Routes       []RouteReport            `json:"routes,omitempty"`
	DaemonSets   []DaemonSetReport        `json:"daemonSets,omitempty"`
	Deployments  []DeploymentReport       `json:"deployments,omitempty"`
	Quotas       []ResourceQuotaReport    `json:"quotas,omitempty"`
}

// PodReport represents json report of k8s pods
type PodReport struct {
	Name string `json:"name"`
}

// QuotaReport represents json report of o7t cluster quotas
type QuotaReport struct {
	Name     string                                   `json:"name"`
	Quota    k8sapicore.ResourceQuotaSpec             `json:"quota,omitempty"`
	Selector o7tapiquota.ClusterResourceQuotaSelector `json:"selector,omitempty"`
}

// ResourceQuotaReport represents json report of Quota resources
type ResourceQuotaReport struct {
	Name          string                          `json:"name"`
	Hard          k8scorev1.ResourceList          `json:"hard,omitempty"`
	ScopeSelector *k8sapicore.ScopeSelector       `json:"selector,omitempty"`
	Scopes        []k8sapicore.ResourceQuotaScope `json:"scopes,omitempty"`
}

// RouteReport represents json report of k8s pods
type RouteReport struct {
	Name              string                             `json:"name"`
	Host              string                             `json:"host"`
	Path              string                             `json:"path,omitempty"`
	AlternateBackends []o7tapiroute.RouteTargetReference `json:"alternateBackends,omitempty"`
	TLS               *o7tapiroute.TLSConfig             `json:"tls,omitempty"`
	To                o7tapiroute.RouteTargetReference   `json:"to"`
	WildcardPolicy    o7tapiroute.WildcardPolicyType     `json:"wildcardPolicy"`
}

// DaemonSetReport represents json report of k8s DaemonSet relevant information
type DaemonSetReport struct {
	Name         string       `json:"name"`
	LatestChange k8sMeta.Time `json:"latestChange,omitempty"`
}

// DeploymentReport represents json report of DeploymentReport resources
type DeploymentReport struct {
	Name         string       `json:"name"`
	LatestChange k8sMeta.Time `json:"latestChange,omitempty"`
}

// ContainerResourcesReport represents json report for aggregated container resources
type ContainerResourcesReport struct {
	ContainerCount int                `json:"containerCount"`
	CPUTotal       *resource.Quantity `json:"cpuTotal"`
	MemoryTotal    *resource.Quantity `json:"memoryTotal"`
}

// PVReport represents json report of k8s PVs
type PVReport struct {
	Name         string                            `json:"name"`
	Driver       k8sapicore.PersistentVolumeSource `json:"driver"`
	StorageClass string                            `json:"storageClass,omitempty"`
	Capacity     k8sapicore.ResourceList           `json:"capacity,omitempty"`
	Phase        k8sapicore.PersistentVolumePhase  `json:"phase,omitempty"`
}

// StorageClassReport represents json report of k8s storage classes
type StorageClassReport struct {
	Name        string `json:"name"`
	Provisioner string `json:"provisioner"`
}

// RBACReport contains RBAC report
type RBACReport struct {
	Users                      []OpenshiftUser                       `json:"users"`
	Groups                     []OpenshiftGroup                      `json:"group"`
	Roles                      []OpenshiftNamespaceRole              `json:"roles"`
	ClusterRoles               []OpenshiftClusterRole                `json:"clusterRoles"`
	ClusterRoleBinding         []OpenshiftClusterRoleBinding         `json:"clusterRoleBindings"`
	SecurityContextConstraints []OpenshiftSecurityContextConstraints `json:"securityContextConstraints"`
}

// OpenshiftUser wrapper around openshift user
type OpenshiftUser struct {
	Name       string   `json:"name"`
	FullName   string   `json:"fullName,omitempty" protobuf:"bytes,2,opt,name=fullName"`
	Identities []string `json:"identities" protobuf:"bytes,3,rep,name=identities"`
	Groups     []string `json:"groups" protobuf:"bytes,4,rep,name=groups"`
}

// OpenshiftGroup wrapper around openshift group
type OpenshiftGroup struct {
	Name  string   `json:"name"`
	Users []string `json:"users" protobuf:"bytes,2,rep,name=users"`
}

// OpenshiftNamespaceRole represent roles mapped to namespaces
type OpenshiftNamespaceRole struct {
	Namespace string          `json:"namespace"`
	Roles     []OpenshiftRole `json:"roles"`
}

// OpenshiftRole wrapper around openshift role
type OpenshiftRole struct {
	Name  string                  `json:"name"`
	Rules []o7tapiauth.PolicyRule `json:"rules,omitempty" protobuf:"bytes,2,rep,name=rules"`
}

// OpenshiftClusterRole wrapper around cluster role
type OpenshiftClusterRole struct {
	Name  string                  `json:"name"`
	Rules []o7tapiauth.PolicyRule `json:"rules,omitempty" protobuf:"bytes,2,rep,name=rules"`
}

// OpenshiftClusterRoleBinding wrapper around openshift cluster role bindings
type OpenshiftClusterRoleBinding struct {
	Name       string                      `json:"name"`
	UserNames  o7tapiauth.OptionalNames    `json:"userNames" protobuf:"bytes,2,rep,name=userNames"`
	GroupNames o7tapiauth.OptionalNames    `json:"groupNames" protobuf:"bytes,3,rep,name=groupNames"`
	Subjects   []k8scorev1.ObjectReference `json:"subjects" protobuf:"bytes,4,rep,name=subjects"`
	RoleRef    k8scorev1.ObjectReference   `json:"roleRef" protobuf:"bytes,5,opt,name=roleRef"`
}

// OpenshiftSecurityContextConstraints wrapper aroung opeshift scc
type OpenshiftSecurityContextConstraints struct {
	Name   string   `json:"name"`
	Users  []string `json:"users" protobuf:"bytes,18,rep,name=users"`
	Groups []string `json:"groups" protobuf:"bytes,19,rep,name=groups"`
}

// GenClusterReport inserts report values into structures for json output
func GenClusterReport(apiResources api.Resources) (clusterReport Report) {
	clusterReport.ReportQuotas(apiResources)
	clusterReport.ReportNamespaces(apiResources)
	clusterReport.ReportNodes(apiResources)
	clusterReport.ReportPVs(apiResources)
	clusterReport.ReportRBAC(apiResources)
	clusterReport.ReportStorageClasses(apiResources)
	return
}

// ReportContainerResources create report about container resources
func ReportContainerResources(reportedNamespace *NamespaceReport, pod *k8sapicore.Pod) {
	cpuTotal := reportedNamespace.Resources.CPUTotal.Value()
	memoryTotal := reportedNamespace.Resources.MemoryTotal.Value()

	for _, container := range pod.Spec.Containers {
		cpuTotal += container.Resources.Requests.Cpu().Value()
		memoryTotal += container.Resources.Requests.Memory().Value()
	}
	reportedNamespace.Resources.CPUTotal.Set(cpuTotal)
	reportedNamespace.Resources.MemoryTotal.Set(memoryTotal)
	reportedNamespace.Resources.ContainerCount += len(pod.Spec.Containers)
}

// ReportDaemonSets generate DaemonSet report
func ReportDaemonSets(reporeportedNamespace *NamespaceReport, dsList *k8sapiapps.DaemonSetList) {
	for _, ds := range dsList.Items {
		reportedDS := DaemonSetReport{
			Name:         ds.Name,
			LatestChange: ds.ObjectMeta.CreationTimestamp,
		}

		reporeportedNamespace.DaemonSets = append(reporeportedNamespace.DaemonSets, reportedDS)
	}
}

// ReportDeployments generate Deployments report
func ReportDeployments(reportedNamespace *NamespaceReport, deploymentList *k8sapiapps.DeploymentList) {
	for _, deployment := range deploymentList.Items {
		reportedDeployment := DeploymentReport{
			Name:         deployment.Name,
			LatestChange: deployment.ObjectMeta.CreationTimestamp,
		}

		reportedNamespace.Deployments = append(reportedNamespace.Deployments, reportedDeployment)
	}
}

// ReportNamespaces fills in information about Namespaces
func (clusterReport *Report) ReportNamespaces(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportNamespaces")

	for _, resources := range apiResources.NamespaceList {
		reportedNamespace := NamespaceReport{Name: resources.NamespaceName}

		ReportResourceQuotas(&reportedNamespace, resources.ResourceQuotaList)
		ReportPods(&reportedNamespace, resources.PodList)
		ReportResources(&reportedNamespace, resources.PodList)
		ReportRoutes(&reportedNamespace, resources.RouteList)
		ReportDeployments(&reportedNamespace, resources.DeploymentList)
		ReportDaemonSets(&reportedNamespace, resources.DaemonSetList)
		clusterReport.Namespaces = append(clusterReport.Namespaces, reportedNamespace)
	}
}

// ReportNodes fills in information about nodes
func (clusterReport *Report) ReportNodes(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportNodes")

	for _, node := range apiResources.NodeList.Items {
		nodeReport := NodeReport{
			Name: node.ObjectMeta.Name,
		}

		isMaster, ok := node.ObjectMeta.Labels["node-role.kubernetes.io/master"]
		nodeReport.MasterNode = ok && isMaster == "true"

		ReportNodeResources(&nodeReport, node.Status, apiResources)
		clusterReport.Nodes = append(clusterReport.Nodes, nodeReport)
	}
}

// ReportNodeResources parse and insert info about consumed resources
func ReportNodeResources(repotedNode *NodeReport, nodeStatus k8sapicore.NodeStatus, apiResources api.Resources) {
	repotedNode.Resources.CPU = nodeStatus.Capacity.Cpu()

	repotedNode.Resources.MemoryCapacity = nodeStatus.Capacity.Memory()

	memConsumed := new(resource.Quantity)
	memCapacity, _ := nodeStatus.Capacity.Memory().AsInt64()
	memAllocatable, _ := nodeStatus.Allocatable.Memory().AsInt64()
	memConsumed.Set(memCapacity - memAllocatable)
	memConsumed.Format = resource.BinarySI
	repotedNode.Resources.MemoryConsumed = memConsumed

	var runningPodsCount int64
	for _, resources := range apiResources.NamespaceList {
		for _, pod := range resources.PodList.Items {
			if pod.Spec.NodeName == repotedNode.Name {
				runningPodsCount++
			}
		}
	}
	podsRunning := new(resource.Quantity)
	podsRunning.Set(runningPodsCount)
	podsRunning.Format = resource.DecimalSI
	repotedNode.Resources.RunningPods = podsRunning

	repotedNode.Resources.PodCapacity = nodeStatus.Capacity.Pods()
}

// ReportQuotas creates report about cluster quotas
func (clusterReport *Report) ReportQuotas(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportQuotas")

	for _, quota := range apiResources.QuotaList.Items {
		quotaReport := QuotaReport{
			Name:     quota.ObjectMeta.Name,
			Quota:    quota.Spec.Quota,
			Selector: quota.Spec.Selector,
		}

		clusterReport.Quotas = append(clusterReport.Quotas, quotaReport)
	}
}

// ReportPods creates info about cluster pods
func ReportPods(reportedNamespace *NamespaceReport, podList *k8sapicore.PodList) {
	for _, pod := range podList.Items {
		reportedPod := PodReport{Name: pod.Name}
		reportedNamespace.Pods = append(reportedNamespace.Pods, reportedPod)

		// Update namespace touch timestamp
		if pod.ObjectMeta.CreationTimestamp.Time.Unix() > reportedNamespace.LatestChange.Time.Unix() {
			reportedNamespace.LatestChange = pod.ObjectMeta.CreationTimestamp
		}
	}
}

// ReportResourceQuotas creates report about quotas
func ReportResourceQuotas(reportedNamespace *NamespaceReport, quotaList *k8sapicore.ResourceQuotaList) {
	for _, quota := range quotaList.Items {
		reportedQuota := ResourceQuotaReport{
			Name:          quota.ObjectMeta.Name,
			Hard:          quota.Spec.Hard,
			ScopeSelector: quota.Spec.ScopeSelector,
			Scopes:        quota.Spec.Scopes,
		}
		reportedNamespace.Quotas = append(reportedNamespace.Quotas, reportedQuota)
	}
}

// ReportResources create report about namespace resources
func ReportResources(reportedNamespace *NamespaceReport, podList *k8sapicore.PodList) {
	resources := ContainerResourcesReport{
		CPUTotal:    &resource.Quantity{Format: resource.DecimalSI},
		MemoryTotal: &resource.Quantity{Format: resource.BinarySI},
	}
	reportedNamespace.Resources = resources

	for _, pod := range podList.Items {
		ReportContainerResources(reportedNamespace, &pod)
	}
}

// ReportRoutes create report about routes
func ReportRoutes(reportedNamespace *NamespaceReport, routeList *o7tapiroute.RouteList) {
	for _, route := range routeList.Items {
		reportedRoute := RouteReport{
			Name:              route.Name,
			AlternateBackends: route.Spec.AlternateBackends,
			Host:              route.Spec.Host,
			Path:              route.Spec.Path,
			To:                route.Spec.To,
			TLS:               route.Spec.TLS,
			WildcardPolicy:    route.Spec.WildcardPolicy,
		}

		reportedNamespace.Routes = append(reportedNamespace.Routes, reportedRoute)
	}
}

// ReportPVs create report oabout pvs
func (clusterReport *Report) ReportPVs(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportPVs")
	pvList := apiResources.PersistentVolumeList

	// Go through all PV and save required information to report
	for _, pv := range pvList.Items {
		reportedPV := PVReport{
			Name:         pv.Name,
			Driver:       pv.Spec.PersistentVolumeSource,
			StorageClass: pv.Spec.StorageClassName,
			Capacity:     pv.Spec.Capacity,
			Phase:        pv.Status.Phase,
		}

		clusterReport.PVs = append(clusterReport.PVs, reportedPV)
	}
}

// ReportRBAC create report about RBAC policy
func (clusterReport *Report) ReportRBAC(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportRBAC")

	clusterReport.RBACReport.Users = make([]OpenshiftUser, 0)
	for _, user := range apiResources.RBACResources.UsersList.Items {
		reportedUser := OpenshiftUser{
			Name:       user.Name,
			FullName:   user.FullName,
			Identities: user.Identities,
			Groups:     user.Groups,
		}
		clusterReport.RBACReport.Users = append(clusterReport.RBACReport.Users, reportedUser)
	}

	clusterReport.RBACReport.Groups = make([]OpenshiftGroup, 0)
	for _, group := range apiResources.RBACResources.GroupList.Items {
		reportedGroup := OpenshiftGroup{
			Name:  group.Name,
			Users: group.Users,
		}

		clusterReport.RBACReport.Groups = append(clusterReport.RBACReport.Groups, reportedGroup)
	}

	clusterReport.RBACReport.Roles = make([]OpenshiftNamespaceRole, 0)
	for _, namespace := range apiResources.NamespaceList {
		reportedNamespaceRoles := OpenshiftNamespaceRole{Namespace: namespace.NamespaceName}

		reportedNamespaceRoles.Roles = make([]OpenshiftRole, 0)
		for _, role := range namespace.RolesList.Items {
			reportedRole := OpenshiftRole{
				Name: role.Name,
			}
			reportedNamespaceRoles.Roles = append(reportedNamespaceRoles.Roles, reportedRole)
		}

		clusterReport.RBACReport.Roles = append(clusterReport.RBACReport.Roles, reportedNamespaceRoles)
	}

	clusterReport.RBACReport.ClusterRoles = make([]OpenshiftClusterRole, 0)
	for _, clusterRole := range apiResources.RBACResources.ClusterRolesList.Items {
		reportedClusterRole := OpenshiftClusterRole{
			Name: clusterRole.Name,
		}

		clusterReport.RBACReport.ClusterRoles = append(clusterReport.RBACReport.ClusterRoles, reportedClusterRole)
	}

	clusterReport.RBACReport.ClusterRoleBinding = make([]OpenshiftClusterRoleBinding, 0)
	for _, clusterRoleBinding := range apiResources.RBACResources.ClusterRolesBindingsList.Items {
		reportedClusterRoleBinding := OpenshiftClusterRoleBinding{
			Name:       clusterRoleBinding.Name,
			UserNames:  clusterRoleBinding.UserNames,
			GroupNames: clusterRoleBinding.GroupNames,
			Subjects:   clusterRoleBinding.Subjects,
			RoleRef:    clusterRoleBinding.RoleRef,
		}

		clusterReport.RBACReport.ClusterRoleBinding = append(clusterReport.RBACReport.ClusterRoleBinding, reportedClusterRoleBinding)
	}

	clusterReport.RBACReport.SecurityContextConstraints = make([]OpenshiftSecurityContextConstraints, 0)

	for _, scc := range apiResources.RBACResources.SecurityContextConstraintsList.Items {
		reportedSCC := OpenshiftSecurityContextConstraints{
			Name:   scc.Name,
			Users:  scc.Users,
			Groups: scc.Groups,
		}

		clusterReport.RBACReport.SecurityContextConstraints = append(clusterReport.RBACReport.SecurityContextConstraints, reportedSCC)
	}
}

// ReportStorageClasses create report about storage classes
func (clusterReport *Report) ReportStorageClasses(apiResources api.Resources) {
	logrus.Debug("ClusterReport::ReportStorageClasses")
	// Go through all storage classes and save required information to report
	storageClassList := apiResources.StorageClassList
	for _, storageClass := range storageClassList.Items {
		reportedStorageClass := StorageClassReport{
			Name:        storageClass.Name,
			Provisioner: storageClass.Provisioner,
		}

		clusterReport.StorageClasses = append(clusterReport.StorageClasses, reportedStorageClass)
	}
}
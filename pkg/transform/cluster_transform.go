package transform

import (
	"fmt"

	"github.com/fusor/cpma/pkg/api"
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/transform/cluster"
	"github.com/fusor/cpma/pkg/transform/clusterquota"
	"github.com/fusor/cpma/pkg/transform/quota"
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
)

// ClusterTransformName is the cluster report name
const ClusterTransformName = "Cluster"

// ClusterExtraction holds data extracted from k8s API resources
type ClusterExtraction struct {
	api.Resources
}

// ClusterTransform reprents transform for k8s API resources
type ClusterTransform struct {
}

// Transform converts data collected from an OCP3 API into a useful output
func (e ClusterExtraction) Transform() ([]Output, error) {
	outputs := []Output{}

	if env.Config().GetBool("Manifests") {
		logrus.Info("ClusterTransform::Transform:Manifests")
		manifests, err := e.buildManifestOutput()
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, manifests)
	}

	if env.Config().GetBool("Reporting") {
		logrus.Info("ClusterTransform::Transform:Reports")

		clusterReport := cluster.GenClusterReport(api.Resources{
			QuotaList:            e.QuotaList,
			NamespaceList:        e.NamespaceList,
			NodeList:             e.NodeList,
			PersistentVolumeList: e.PersistentVolumeList,
			RBACResources:        e.RBACResources,
			StorageClassList:     e.StorageClassList,
		})

		FinalReportOutput.Report.ClusterReport = clusterReport
	}

	return outputs, nil
}

func (e ClusterExtraction) buildManifestOutput() (Output, error) {
	var manifests []Manifest

	for _, clusterQuota := range e.QuotaList.Items {
		clusterQuotaCR, err := clusterquota.Translate(clusterQuota)
		quotaCRYAML, err := GenYAML(clusterQuotaCR)
		if err != nil {
			return nil, err
		}
		name := fmt.Sprintf("100_CPMA-cluster-quota-resource-%s.yaml", clusterQuota.Name)
		manifest := Manifest{Name: name, CRD: quotaCRYAML}
		manifests = append(manifests, manifest)
	}

	for _, clusterNamespace := range e.NamespaceList {
		for _, resourceQuota := range clusterNamespace.ResourceQuotaList.Items {
			quotaCR, err := quota.Translate(resourceQuota)
			quotaCRYAML, err := GenYAML(quotaCR)
			if err != nil {
				return nil, err
			}
			name := fmt.Sprintf("100_CPMA-%s-resource-quota-%s.yaml", resourceQuota.Namespace, resourceQuota.Name)
			manifest := Manifest{Name: name, CRD: quotaCRYAML}
			manifests = append(manifests, manifest)
		}
	}

	return ManifestOutput{
		Manifests: manifests,
	}, nil
}

// Validate no need to validate it, data is exctracted from API
func (e ClusterExtraction) Validate() (err error) { return }

// Extract collects data for cluster report
func (e ClusterTransform) Extract() (Extraction, error) {
	chanNodes := make(chan *k8sapicore.NodeList)
	chanClusterQuotas := make(chan *o7tapiquota.ClusterResourceQuotaList)
	chanNamespaces := make(chan *k8sapicore.NamespaceList)
	chanPVs := make(chan *k8sapicore.PersistentVolumeList)
	chanUsers := make(chan *o7tapiuser.UserList)
	chanGroups := make(chan *o7tapiuser.GroupList)
	chanClusterRoles := make(chan *o7tapiauth.ClusterRoleList)
	chanClusterRolesListBindings := make(chan *o7tapiauth.ClusterRoleBindingList)
	chanStorageClassList := make(chan *k8sapistorage.StorageClassList)
	chanSecurityContextConstraints := make(chan *o7tapisecurity.SecurityContextConstraintsList)

	go api.ListNamespaces(chanNamespaces)
	go api.ListNodes(chanNodes)
	go api.ListQuotas(chanClusterQuotas)
	go api.ListPVs(chanPVs)
	go api.ListUsers(chanUsers)
	go api.ListGroups(chanGroups)
	go api.ListClusterRoles(chanClusterRoles)
	go api.ListClusterRolesBindings(chanClusterRolesListBindings)
	go api.ListSCC(chanSecurityContextConstraints)
	go api.ListStorageClasses(chanStorageClassList)

	extraction := &ClusterExtraction{}

	// Map all namespaces to their resources
	namespacesList := <-chanNamespaces
	namespaceListSize := len(namespacesList.Items)
	extraction.NamespaceList = make([]api.NamespaceResources, namespaceListSize, namespaceListSize)
	for i, namespace := range namespacesList.Items {
		namespaceResources := api.NamespaceResources{NamespaceName: namespace.Name}

		chanQuotas := make(chan *k8sapicore.ResourceQuotaList)
		chanPods := make(chan *k8sapicore.PodList)
		chanRoutes := make(chan *o7tapiroute.RouteList)
		chanDeployments := make(chan *v1beta1.DeploymentList)
		chanDaemonSets := make(chan *extv1b1.DaemonSetList)
		chanRoles := make(chan *o7tapiauth.RoleList)
		chanPVCs := make(chan *k8sapicore.PersistentVolumeClaimList)

		go api.ListResourceQuotas(namespace.Name, chanQuotas)
		go api.ListPods(namespace.Name, chanPods)
		go api.ListRoutes(namespace.Name, chanRoutes)
		go api.ListDeployments(namespace.Name, chanDeployments)
		go api.ListDaemonSets(namespace.Name, chanDaemonSets)
		go api.ListRoles(namespace.Name, chanRoles)
		go api.ListPVCs(namespace.Name, chanPVCs)

		namespaceResources.ResourceQuotaList = <-chanQuotas
		namespaceResources.PodList = <-chanPods
		namespaceResources.RouteList = <-chanRoutes
		namespaceResources.DeploymentList = <-chanDeployments
		namespaceResources.DaemonSetList = <-chanDaemonSets
		namespaceResources.RolesList = <-chanRoles
		namespaceResources.PVCList = <-chanPVCs

		extraction.NamespaceList[i] = namespaceResources
	}

	extraction.NodeList = <-chanNodes
	extraction.QuotaList = <-chanClusterQuotas
	extraction.PersistentVolumeList = <-chanPVs
	extraction.RBACResources.UsersList = <-chanUsers
	extraction.RBACResources.GroupList = <-chanGroups
	extraction.RBACResources.ClusterRolesList = <-chanClusterRoles
	extraction.RBACResources.ClusterRolesBindingsList = <-chanClusterRolesListBindings
	extraction.RBACResources.SecurityContextConstraintsList = <-chanSecurityContextConstraints
	extraction.StorageClassList = <-chanStorageClassList

	return *extraction, nil
}

// Name returns a human readable name for the transform
func (e ClusterTransform) Name() string {
	return ClusterTransformName
}

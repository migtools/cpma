package transform

import (
	"fmt"

	"github.com/fusor/cpma/pkg/api"
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/transform/cluster"
	"github.com/fusor/cpma/pkg/transform/clusterquota"
	"github.com/fusor/cpma/pkg/transform/quota"
	"github.com/sirupsen/logrus"
)

// ClusterReportName is the cluster report name
const ClusterReportName = "ClusterReport"

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
	extraction := &ClusterExtraction{}

	nodeList, err := api.ListNodes()
	if err != nil {
		return nil, err
	}
	extraction.NodeList = nodeList

	quotaList, err := api.ListQuotas()
	if err != nil {
		return nil, err
	}
	extraction.QuotaList = quotaList

	namespacesList, err := api.ListNamespaces()
	if err != nil {
		return nil, err
	}

	// Map all namespaces to their resources
	namespaceListSize := len(namespacesList.Items)
	extraction.NamespaceList = make([]api.NamespaceResources, namespaceListSize, namespaceListSize)
	for i, namespace := range namespacesList.Items {
		namespaceResources := api.NamespaceResources{NamespaceName: namespace.Name}

		quotaList, err := api.ListResourceQuotas(namespace.Name)
		if err != nil {
			return nil, err
		}
		namespaceResources.ResourceQuotaList = quotaList

		podsList, err := api.ListPods(namespace.Name)
		if err != nil {
			return nil, err
		}
		namespaceResources.PodList = podsList

		routesList, err := api.ListRoutes(namespace.Name)
		if err != nil {
			return nil, err
		}
		namespaceResources.RouteList = routesList

		deploymentList, err := api.ListDeployments(namespace.Name)
		if err != nil {
			return nil, err
		}
		namespaceResources.DeploymentList = deploymentList

		daemonSetList, err := api.ListDaemonSets(namespace.Name)
		if err != nil {
			return nil, err
		}
		namespaceResources.DaemonSetList = daemonSetList

		rolesList, err := api.ListRoles(namespace.Name)
		if err != nil {
			return nil, err
		}
		namespaceResources.RolesList = rolesList

		extraction.NamespaceList[i] = namespaceResources
	}

	pvList, err := api.ListPVs()
	if err != nil {
		return nil, err
	}
	extraction.PersistentVolumeList = pvList

	storageClassList, err := api.ListStorageClasses()
	if err != nil {
		return nil, err
	}
	extraction.StorageClassList = storageClassList

	userList, err := api.ListUsers()
	if err != nil {
		return nil, err
	}
	extraction.RBACResources.UsersList = userList

	groupList, err := api.ListGroups()
	if err != nil {
		return nil, err
	}
	extraction.RBACResources.GroupList = groupList

	clusterRolesList, err := api.ListClusterRoles()
	if err != nil {
		return nil, err
	}
	extraction.RBACResources.ClusterRolesList = clusterRolesList

	clusterRolesListBindings, err := api.ListClusterRolesBindings()
	if err != nil {
		return nil, err
	}
	extraction.RBACResources.ClusterRolesBindingsList = clusterRolesListBindings

	securityContextConstraints, err := api.ListSCC()
	if err != nil {
		return nil, err
	}
	extraction.RBACResources.SecurityContextConstraintsList = securityContextConstraints

	return *extraction, nil
}

// Name returns a human readable name for the transform
func (e ClusterTransform) Name() string {
	return ClusterReportName
}

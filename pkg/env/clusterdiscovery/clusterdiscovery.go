package clusterdiscovery

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/fusor/cpma/pkg/api"
	"github.com/pkg/errors"

	k8sapicore "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// DiscoverCluster Get kubeconfig using $KUBECONFIG, if not try ~/.kube/config
// parse kubeconfig and select source cluster from available contexts
// query k8s api for nodes, get node urls from api response and survey master node
func DiscoverCluster() (string, string, error) {
	selectedCluster := SurveyClusters()
	if err := api.CreateK8sClient(selectedCluster); err != nil {
		return "", "", errors.Wrap(err, "k8s api client failed to create")
	}

	clusterNodes, err := queryNodes(api.K8sClient)
	if err != nil {
		return "", "", errors.Wrap(err, "cluster node query failed")
	}

	selectedNode := surveyNodes(clusterNodes)

	return selectedNode, selectedCluster, nil
}

// SurveyClusters list clusters from kubeconfig
func SurveyClusters() string {
	// Survey options should be an array
	clusters := make([]string, 0, len(api.ClusterNames))
	// It's better to have current context's cluster first, because
	// it will be easier to select it using survey
	currentContext := api.KubeConfig.CurrentContext

	currentContextCluster := api.KubeConfig.Contexts[currentContext].Cluster
	clusters = append(clusters, currentContextCluster)

	for cluster := range api.ClusterNames {
		if cluster != currentContextCluster {
			clusters = append(clusters, cluster)
		}
	}

	selectedCluster := ""
	prompt := &survey.Select{
		Message: "Select cluster obtained from KUBECONFIG contexts",
		Options: clusters,
	}
	survey.AskOne(prompt, &selectedCluster)

	return selectedCluster
}

func queryNodes(client *kubernetes.Clientset) ([]string, error) {
	chanNodes := make(chan *k8sapicore.NodeList)
	go api.ListNodes(client, chanNodes)
	nodeList := <-chanNodes

	nodes := make([]string, 0, len(nodeList.Items))
	for _, node := range nodeList.Items {
		nodes = append(nodes, node.Name)
	}

	return nodes, nil
}

func surveyNodes(nodes []string) string {
	selectedNode := ""
	prompt := &survey.Select{
		Message: "Select master node",
		Options: nodes,
	}
	survey.AskOne(prompt, &selectedNode)

	return selectedNode
}

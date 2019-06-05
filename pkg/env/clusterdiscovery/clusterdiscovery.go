package clusterdiscovery

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	kubeConfig   *clientcmdapi.Config
	clusterNames = make(map[string]string)
)

// DiscoverCluster Get kubeconfig using $KUBECONFIG, if not try ~/.kube/config
// Parse kubeconfig and select cluster from available contexts, then get server url from context
// query k8s api for nodes, get node urls from api response and survey master node
func DiscoverCluster() (string, error) {
	err := parseKubeConfig()
	if err != nil {
		return "", errors.Wrap(err, "kubeconfig parsing failed")
	}

	selectedCluster := surveyClusters()

	apiClient, err := createAPIClient(selectedCluster)
	if err != nil {
		return "", errors.Wrap(err, "k8s api client failed to create")
	}

	clusterNodes, err := queryNodes(apiClient)
	if err != nil {
		return "", errors.Wrap(err, "cluster node query failed")
	}

	selectedNode := surveyNodes(clusterNodes)

	return selectedNode, nil
}

func parseKubeConfig() error {
	kubeConfigPath, err := getKubeConfigPath()
	if err != nil {
		return err
	}

	kubeConfigFile, err := ioutil.ReadFile(kubeConfigPath)
	if err != nil {
		return err
	}

	kubeConfig, err = clientcmd.Load(kubeConfigFile)
	if err != nil {
		return err
	}

	// Map context clusters and name for easier access in future
	for name, context := range kubeConfig.Contexts {
		clusterNames[context.Cluster] = name
	}

	return nil
}

func getKubeConfigPath() (string, error) {
	// Get kubeconfig using $KUBECONFIG, if not try ~/.kube/config
	var kubeConfigPath string

	kubeconfigEnv := os.Getenv("KUBECONFIG")
	if kubeconfigEnv != "" {
		kubeConfigPath = kubeconfigEnv
	} else {
		home, err := homedir.Dir()
		if err != nil {
			return "", errors.Wrap(err, "Can't detect home user directory")
		}
		kubeConfigPath = fmt.Sprintf("%s/.kube/config", home)
	}

	return kubeConfigPath, nil
}

func surveyClusters() string {
	// Survey options should be an array
	clusters := make([]string, 0, len(clusterNames))
	// It's better to have current context's cluster first, because
	// it will be easier to select it using survey
	currentContext := kubeConfig.CurrentContext
	currentContextCluster := kubeConfig.Contexts[currentContext].Cluster
	clusters = append(clusters, currentContextCluster)

	for cluster := range clusterNames {
		if cluster != currentContextCluster {
			clusters = append(clusters, cluster)
		}
	}

	selectedCluster := ""
	prompt := &survey.Select{
		Message: "Select cluster obtained from KUBECONFIG contexts",
		Options: clusters,
	}
	survey.AskOne(prompt, &selectedCluster, nil)

	return selectedCluster
}

func createAPIClient(selectedCluster string) (corev1.CoreV1Interface, error) {
	// set current context to selected cluster for connecting to cluster using client-go
	kubeConfig.CurrentContext = clusterNames[selectedCluster]

	var kubeConfigGetter = func() (*clientcmdapi.Config, error) {
		return kubeConfig, nil
	}
	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", kubeConfigGetter)

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1(), nil
}

func queryNodes(apiClient corev1.CoreV1Interface) ([]string, error) {
	listOptions := metav1.ListOptions{}

	nodeList, err := apiClient.Nodes().List(listOptions)
	if err != nil {
		return nil, err
	}

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
	survey.AskOne(prompt, &selectedNode, nil)

	return selectedNode
}

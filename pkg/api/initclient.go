package api

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	// KubeConfig represents kubeconfig
	KubeConfig *clientcmdapi.Config
	// ClusterNames contains names of contexts and cluster
	ClusterNames = make(map[string]string)
	// Client api client used for connecting to k8s api
	Client *kubernetes.Clientset
)

// ParseKubeConfig parse kubeconfig
func ParseKubeConfig() error {
	kubeConfigPath, err := getKubeConfigPath()
	if err != nil {
		return err
	}

	kubeConfigFile, err := ioutil.ReadFile(kubeConfigPath)
	if err != nil {
		return err
	}

	KubeConfig, err = clientcmd.Load(kubeConfigFile)
	if err != nil {
		return err
	}

	// Map context clusters and name for easier access in future
	for name, context := range KubeConfig.Contexts {
		ClusterNames[context.Cluster] = name
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

// CreateAPIClient create api client using cluster from kubeconfig context
func CreateAPIClient(contextCluster string) error {
	// set current context to selected cluster for connecting to cluster using client-go
	KubeConfig.CurrentContext = ClusterNames[contextCluster]

	var kubeConfigGetter = func() (*clientcmdapi.Config, error) {
		return KubeConfig, nil
	}
	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", kubeConfigGetter)

	Client, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	logrus.Debugf("API client initialized for %s", contextCluster)

	return nil
}

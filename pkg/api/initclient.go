package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	// KubeConfig represents kubeconfig
	KubeConfig *clientcmdapi.Config
	// ClusterNames contains names of contexts and cluster
	ClusterNames = make(map[string]string)
	// K8sClient api client used for connecting to k8s api
	K8sClient *kubernetes.Clientset
	// O7tClient api client used for connecting to Openshift api
	O7tClient *OpenshiftClient

	kubeConfigGetter = func() (*clientcmdapi.Config, error) {
		return KubeConfig, nil
	}
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

// CreateK8sClient create api client using cluster from kubeconfig context
func CreateK8sClient(contextCluster string) error {
	if K8sClient == nil {
		config, err := buildConfig(contextCluster)
		if err != nil {
			return err
		}

		K8sClient = NewK8SOrDie(config)
		logrus.Debugf("Kubernetes API client initialized for %s", contextCluster)
	}

	return nil
}

// CreateO7tClient create api client using cluster from kubeconfig context
func CreateO7tClient(contextCluster string) error {
	if O7tClient == nil {
		config, err := buildConfig(contextCluster)
		if err != nil {
			return err
		}

		O7tClient = NewO7tOrDie(config)
		logrus.Debugf("Openshift API client initialized for %s", contextCluster)
	}

	return nil
}

func buildConfig(contextCluster string) (*rest.Config, error) {
	// Check if context is present in kubeconfig
	if err := validateConfig(contextCluster); err != nil {
		return nil, err
	}

	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", kubeConfigGetter)
	if err != nil {
		return nil, errors.Wrap(err, "Error in KUBECONFIG")
	}

	config.AcceptContentTypes = "application/vnd.kubernetes.protobuf,application/json"
	config.UserAgent = fmt.Sprintf(
		"cpma/v1.0 (%s/%s) kubernetes/v1.0",
		runtime.GOOS, runtime.GOARCH,
	)

	return config, nil
}

func validateConfig(contextCluster string) error {
	for context := range ClusterNames {
		if context == contextCluster {
			return nil
		}
	}

	return errors.New("Not valid context")
}

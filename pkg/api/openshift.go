package api

import (
	"sync"

	authv1 "github.com/openshift/client-go/authorization/clientset/versioned/typed/authorization/v1"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// OpenshiftClient - Client to interact with openshift api
type OpenshiftClient struct {
	authClient  authv1.AuthorizationV1Interface
	routeClient routev1.RouteV1Interface
}

var instances struct {
	Openshift *OpenshiftClient
}

var once struct {
	Openshift sync.Once
}

func createClientConfigFromFile(configPath string) (*rest.Config, error) {
	kubeConfigPath, err := getKubeConfigPath()
	if err != nil {
		return nil, err
	}

	clientConfig, err := clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.NewDefaultClientConfig(*clientConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

// Openshift - Create a new openshift client if needed, returns reference
func Openshift() (*OpenshiftClient, error) {
	errMsg := "Something went wrong while initializing openshift client!\n"
	once.Openshift.Do(func() {
		client, err := newOpenshift()
		if err != nil {
			log.Error(errMsg)
			// NOTE: Looking to leverage panic recovery to gracefully handle this
			// with things like retries or better intelligence, but the environment
			// is probably in a unrecoverable state as far as the broker is concerned,
			// and demands the attention of an operator.
			panic(err.Error())
		}
		instances.Openshift = client
	})
	if instances.Openshift == nil {
		return nil, errors.New("OpenShift client instance is nil")
	}
	return instances.Openshift, nil
}

func newOpenshift() (*OpenshiftClient, error) {
	// NOTE: Both the external and internal client object are using the same
	// clientset library. Internal clientset normally uses a different
	// library
	clientConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Debug("Checking for a local Cluster Config")
		clientConfig, err = createClientConfigFromFile(homedir.HomeDir() + "/.kube/config")
		if err != nil {
			log.Error("Failed to create LocalClientSet")
			return nil, err
		}
	}

	clientset, err := newForConfig(clientConfig)
	if err != nil {
		log.Error("Failed to create LocalClientSet")
		return nil, err
	}

	return clientset, err
}

func newForConfig(c *rest.Config) (*OpenshiftClient, error) {
	authClient, err := authv1.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	routeClient, err := routev1.NewForConfig(c)
	if err != nil {
		return nil, err
	}
	return &OpenshiftClient{authClient: authClient, routeClient: routeClient}, nil
}

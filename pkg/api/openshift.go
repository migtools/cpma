package api

import (
	"sync"

	authv1 "github.com/openshift/client-go/authorization/clientset/versioned/typed/authorization/v1"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"github.com/pkg/errors"
	"k8s.io/client-go/rest"
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

// Openshift - Create a new openshift client if needed, returns reference
func Openshift(config *rest.Config) (*OpenshiftClient, error) {
	once.Openshift.Do(func() {
		client, _ := newOpenshift(config)
		instances.Openshift = client
	})
	if instances.Openshift == nil {
		return nil, errors.New("OpenShift client instance is nil")
	}
	return instances.Openshift, nil
}

func newOpenshift(config *rest.Config) (*OpenshiftClient, error) {
	authClient, err := authv1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	routeClient, err := routev1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &OpenshiftClient{authClient: authClient, routeClient: routeClient}, nil
}

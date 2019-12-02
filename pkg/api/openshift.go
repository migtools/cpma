package api

import (
	authv1 "github.com/openshift/client-go/authorization/clientset/versioned/typed/authorization/v1"
	quotav1 "github.com/openshift/client-go/quota/clientset/versioned/typed/quota/v1"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	security1 "github.com/openshift/client-go/security/clientset/versioned/typed/security/v1"
	userv1 "github.com/openshift/client-go/user/clientset/versioned/typed/user/v1"

	"k8s.io/client-go/rest"
)

// OpenshiftClient - Client to interact with openshift api
type OpenshiftClient struct {
	authClient     authv1.AuthorizationV1Interface
	quotaClient    quotav1.QuotaV1Interface
	routeClient    routev1.RouteV1Interface
	securityClient security1.SecurityV1Interface
	userClient     userv1.UserV1Interface
}

// NewO7tOrDie - Create a new openshift client if needed, returns reference
func NewO7tOrDie(config *rest.Config) *OpenshiftClient {
	return &OpenshiftClient{
		authClient:     authv1.NewForConfigOrDie(config),
		quotaClient:    quotav1.NewForConfigOrDie(config),
		routeClient:    routev1.NewForConfigOrDie(config),
		securityClient: security1.NewForConfigOrDie(config),
		userClient:     userv1.NewForConfigOrDie(config),
	}
}

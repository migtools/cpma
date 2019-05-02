package ocp4

import (
	"github.com/fusor/cpma/pkg/ocp4/oauth"
	"github.com/fusor/cpma/pkg/ocp4/secrets"
)

type Cluster struct {
	Master Master
}

type Master struct {
	OAuth   oauth.OAuthCRD
	Secrets []secrets.Secret
}

type Manifests []Manifest

// Manifest holds a CRD object
type Manifest struct {
	Name string
	CRD  []byte
}

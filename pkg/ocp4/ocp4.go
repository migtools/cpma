package ocp4

import (
	"github.com/fusor/cpma/pkg/ocp4/oauth"
	"github.com/fusor/cpma/pkg/ocp4/sdn"
	"github.com/fusor/cpma/pkg/ocp4/secrets"
	configv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/sirupsen/logrus"
)

type Master struct {
	OAuth   oauth.OAuthCRD
	Secrets []secrets.Secret
	Network sdn.NetworkCR
}

type Node struct {
	Secrets []secrets.Secret
}

type Manifests []Manifest

// Manifest holds a CRD object
type Manifest struct {
	Name string
	CRD  []byte
}

// OCP4InstallMsg message about using generated manifests
const OCP4InstallMsg = `To install OCP4 run the installer as follow in order to add CRDs:
' /openshift-install --dir $INSTALL_DIR create install-config'
'./openshift-install --dir $INSTALL_DIR create manifests'
# Copy generated CRD manifest files  to '$INSTALL_DIR/openshift/'
# Edit them if needed, then run installation:
'./openshift-install --dir $INSTALL_DIR  create cluster'`

const manifestPrefix = "100_CPMA-cluster-"

func (ocp4Node *Node) Transform(ocp3Node configv1.NodeConfig) {
}

func OAuthManifest(oauthCRKind string, crd []byte, manifests Manifests) Manifests {
	if oauthCRKind != "" {
		filename := manifestPrefix + "config-oauth.yaml"
		m := Manifest{Name: filename, CRD: crd}
		manifests = append(manifests, m)
	} else {
		logrus.Debugln("Skipping oauth, no manifests found")
	}
	return manifests
}

func SecretsManifest(secret secrets.Secret, crd []byte, manifests Manifests) Manifests {
	filename := manifestPrefix + "config-secret-" + secret.Metadata.Name + ".yaml"
	m := Manifest{Name: filename, CRD: crd}
	manifests = append(manifests, m)
	return manifests
}

func SDNManifest(networkCR []byte, manifests Manifests) Manifests {
	filename := manifestPrefix + "config-sdn.yaml"
	m := Manifest{Name: filename, CRD: networkCR}
	manifests = append(manifests, m)
	return manifests
}

// GenYAML returns the list of translated CRDs
func (ocp4Node *Node) GenYAML() []Manifest {
	return nil
}

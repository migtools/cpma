package ocp4

import (
	"github.com/fusor/cpma/pkg/ocp4/secrets"
	"github.com/sirupsen/logrus"
)

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

func OAuthManifest(oauthCRKind string, crd []byte, manifests Manifests) Manifests {
	if oauthCRKind != "" {
		filename := manifestPrefix + "config-oauth.yaml"
		manifest := Manifest{Name: filename, CRD: crd}
		manifests = append(manifests, manifest)
	} else {
		logrus.Debugln("Skipping oauth, no manifests found")
	}
	return manifests
}

func SecretsManifest(secret secrets.Secret, crd []byte, manifests Manifests) Manifests {
	filename := manifestPrefix + "config-secret-" + secret.Metadata.Name + ".yaml"
	manifest := Manifest{Name: filename, CRD: crd}
	manifests = append(manifests, manifest)
	return manifests
}

func SDNManifest(networkCR []byte, manifests Manifests) Manifests {
	filename := manifestPrefix + "config-sdn.yaml"
	manifest := Manifest{Name: filename, CRD: networkCR}
	manifests = append(manifests, manifest)
	return manifests
}

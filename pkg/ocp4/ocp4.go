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

// Translate translate OCP3 configs
func (ocp4Master *Master) Translate(ocp3Master configv1.MasterConfig) {
	if ocp3Master.OAuthConfig != nil {
		logrus.Debugln("Translating oauth config")
		oauth, secrets, err := oauth.Translate(ocp3Master.OAuthConfig)

		if err != nil {
			logrus.WithError(err).Fatalf("Unable to generate OAuth CRD from %+v", ocp3Master.OAuthConfig)
		}
		ocp4Master.OAuth = *oauth
		ocp4Master.Secrets = secrets
	}

	if &ocp3Master.NetworkConfig != nil {
		logrus.Debugln("Translating SDN config")
		networkCR := sdn.Translate(ocp3Master.NetworkConfig)
		ocp4Master.Network = *networkCR
	}
}

func (ocp4Node *Node) Translate(ocp3Node configv1.NodeConfig) {
}

// GenYAML returns the list of translated CRDs
func (ocp4Master *Master) GenYAML() []Manifest {
	var manifests []Manifest

	// Generate yaml for oauth config
	crd := ocp4Master.OAuth.GenYAML()
	manifests = oauthManifests(ocp4Master.OAuth.Kind, crd, manifests)

	// Generate yaml for oauth secrets
	manifests = oauthSecrets(ocp4Master.Secrets, manifests)

	// Generate yaml for SDN config
	networkCR := ocp4Master.Network.GenYAML()
	manifests = sdnManifest(networkCR, manifests)

	return manifests
}

func oauthManifests(oauthCRKind string, crd []byte, manifests []Manifest) []Manifest {
	if oauthCRKind != "" {
		filename := manifestPrefix + "config-oauth.yaml"
		manifest := Manifest{Name: filename, CRD: crd}
		manifests = append(manifests, manifest)
	} else {
		logrus.Debugln("Skipping oauth, no manifests found")
	}

	return manifests
}

func oauthSecrets(ocp4MasterSecrets []secrets.Secret, manifests []Manifest) []Manifest {
	for _, secret := range ocp4MasterSecrets {
		filename := manifestPrefix + "config-secret-" + secret.Metadata.Name + ".yaml"
		m := Manifest{Name: filename, CRD: secret.GenYAML()}
		manifests = append(manifests, m)
	}

	return manifests
}

func sdnManifest(networkCR []byte, manifests []Manifest) []Manifest {
	filename := manifestPrefix + "config-sdn.yaml"
	manifest := Manifest{Name: filename, CRD: networkCR}
	manifests = append(manifests, manifest)

	return manifests
}

// GenYAML returns the list of translated CRDs
func (ocp4Node *Node) GenYAML() []Manifest {
	return nil
}

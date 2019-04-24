package ocp4

import (
	"github.com/fusor/cpma/ocp3"
	"github.com/fusor/cpma/ocp4/oauth"
	"github.com/fusor/cpma/ocp4/secrets"
	configv1 "github.com/openshift/api/legacyconfig/v1"
	"github.com/sirupsen/logrus"
)

type Cluster struct {
	Master Master
}

type Master struct {
	OAuth   oauth.OAuthCRD
	Secrets []secrets.Secret
}

// Manifest holds a CRD object
type Manifest struct {
	Name string
	CRD  []byte
}

const OCP4InstallMsg = `To install OCP4 run the installer as follow in order to add CRDs:
' /openshift-install --dir $INSTALL_DIR create install-config'
'./openshift-install --dir $INSTALL_DIR create manifests'
# Copy generated CRD manifest files  to '$INSTALL_DIR/openshift/'
# Edit them if needed, then run installation:
'./openshift-install --dir $INSTALL_DIR  create cluster'`

func (clusterOCP4 *Cluster) Translate(cluster ocp3.Cluster) {
	clusterOCP4.Master.TranslateMaster(cluster.Master.Config)
}

func (masterOCP4 *Master) TranslateMaster(masterOCP3 configv1.MasterConfig) {
	if masterOCP3.OAuthConfig != nil {
		oauth, secrets, err := oauth.Translate(masterOCP3.OAuthConfig)
		if err != nil {
			logrus.WithError(err).Fatalf("Unable to generate OAuth CRD from %+v", masterOCP3.OAuthConfig)
		}
		masterOCP4.OAuth = *oauth
		masterOCP4.Secrets = secrets
	}
}

// GenYAML returns the list of translated CRDs
func (clusterOCP4 *Cluster) GenYAML() []Manifest {
	var manifests []Manifest

	masterManifests := clusterOCP4.Master.MasterGenYAML()
	for _, manifest := range masterManifests {
		manifests = append(manifests, manifest)
	}
	return manifests
}

// GenYAML returns the list of translated CRDs for the Master configuration
func (master *Master) MasterGenYAML() []Manifest {
	var manifests []Manifest
	manifest := Manifest{Name: "100_CPMA-cluster-config-oauth.yaml", CRD: master.OAuth.GenYAML()}
	manifests = append(manifests, manifest)

	for _, secret := range master.Secrets {
		filename := "100_CPMA-cluster-config-secret-" + secret.Metadata.Name + ".yaml"
		m := Manifest{Name: filename, CRD: secret.GenYAML()}
		manifests = append(manifests, m)
	}
	return manifests
}

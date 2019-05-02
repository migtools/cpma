package ocp4

import (
	"errors"

	"github.com/fusor/cpma/pkg/ocp4/oauth"
	"github.com/fusor/cpma/pkg/ocp4/secrets"
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

type Manifests []Manifest

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

func (ocp4Master *Master) Translate(ocp3Master configv1.MasterConfig) {
	if ocp3Master.OAuthConfig != nil {
		oauth, secrets, err := oauth.Translate(ocp3Master.OAuthConfig)
		if err != nil {
			logrus.WithError(err).Fatalf("Unable to generate OAuth CRD from %+v", ocp3Master.OAuthConfig)
		}
		ocp4Master.OAuth = *oauth
		ocp4Master.Secrets = secrets
	}
}

// GenYAML returns the list of translated CRDs
func (ocp4Master *Master) GenYAML() ([]Manifest, error) {
	var manifests []Manifest
	if ocp4Master.OAuth.Kind != "" {
		manifest := Manifest{Name: "100_CPMA-cluster-config-oauth.yaml", CRD: ocp4Master.OAuth.GenYAML()}
		manifests = append(manifests, manifest)

		for _, secret := range ocp4Master.Secrets {
			filename := "100_CPMA-cluster-config-secret-" + secret.Metadata.Name + ".yaml"
			m := Manifest{Name: filename, CRD: secret.GenYAML()}
			manifests = append(manifests, m)
		}
		return manifests, nil
	}
	return nil, errors.New("No manifests")
}

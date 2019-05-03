package ocp4

import (
	"github.com/fusor/cpma/pkg/ocp3"
	"github.com/fusor/cpma/pkg/ocp4/oauth"
	"github.com/sirupsen/logrus"
)

const OCP4InstallMsg = `To install OCP4 run the installer as follow in order to add CRDs:
' /openshift-install --dir $INSTALL_DIR create install-config'
'./openshift-install --dir $INSTALL_DIR create manifests'
# Copy generated CRD manifest files  to '$INSTALL_DIR/openshift/'
# Edit them if needed, then run installation:
'./openshift-install --dir $INSTALL_DIR  create cluster'`

func (ocp4Master *Master) Translate(cluster ocp3.Cluster) {
	if cluster.MasterConfig.OAuthConfig != nil {
		oauth, secrets, err := oauth.Translate(cluster)
		if err != nil {
			logrus.WithError(err).Fatalf("Unable to generate OAuth CRD from %+v", cluster.MasterConfig.OAuthConfig)
		}
		ocp4Master.OAuth = *oauth
		ocp4Master.Secrets = secrets
	}
}

package ocp4

import (
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

func (cluster *Cluster) Translate(masterconfig configv1.MasterConfig) {
	// TODO: Pass just a part of required structure instead of full master config.
	oauth, secrets, err := oauth.Translate(masterconfig)
	if err != nil {
		logrus.WithError(err).Fatalf("Unable to generate OAuth CRD from %+v", masterconfig.OAuthConfig)
	}
	cluster.Master.OAuth = *oauth
	cluster.Master.Secrets = secrets
}

// PrintCRD Prints translated CRDs
func (cluster *Cluster) PrintCRD() {
	oauthCRD := cluster.Master.OAuth.PrintCRD()
	logrus.Print(oauthCRD)

	for _, secret := range cluster.Master.Secrets {
		logrus.Print(secret.PrintCRD())
	}
}

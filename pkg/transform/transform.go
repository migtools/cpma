package transform

import (
	"github.com/fusor/cpma/pkg/etl"
)

//Start generating manifests to be used with Openshift 4
func Start() {
	config := etl.LoadConfig()
	runner := etl.NewRunner(config)

	runner.Transform([]etl.Transform{
		OAuthTransform{
			Config: &config,
		},
		RegistriesTransform{
			Config: &config,
		},
	})
}

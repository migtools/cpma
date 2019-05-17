package report

import (
	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/transform"
)

//Start generating reports about what we can migrate to OCP4
func Start() {
	env.Config().Set("mode", "report")
	config := transform.LoadConfig()
	runner := transform.NewRunner(config)

	runner.Transform([]transform.Transform{
		transform.OAuthTransform{
			Config: &config,
		},
		transform.SDNTransform{
			Config: &config,
		},
		transform.RegistriesTransform{
			Config: &config,
		},
	})
}

package report

import "github.com/fusor/cpma/pkg/etl"

//Start generating a component transform confidence report
func Start() {
	config := etl.LoadConfig()
	runner := etl.NewRunner(config)

	runner.Transform([]etl.Transform{
		SDNReport{
			Config: &config,
		},
	})
}

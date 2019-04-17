// Copyright © 2019 Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"path"

	"github.com/fusor/cpma/env"
	ocp3 "github.com/fusor/cpma/ocp3config"
	"github.com/fusor/cpma/ocp4crd/oauth"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var debugLogLevel bool

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringVar(&env.ConfigFile, "config", "", "config file (default is $HOME/.cpma.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debugLogLevel, "debug", false, "show debug ouput")

	rootCmd.Flags().StringP("output-dir", "o", "", "set the directory to store extracted configuration.")
	env.Config().BindPFlag("OutputDir", rootCmd.Flags().Lookup("output-dir"))
	env.Config().SetDefault("OutputDir", path.Dir(""))
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cpma",
	Short: "Helps migration cluster configuration of a OCP 3.x cluster to OCP 4.x",
	Long:  `Helps migration cluster configuration of a OCP 3.x cluster to OCP 4.x`,
	Run: func(cmd *cobra.Command, args []string) {
		if debugLogLevel {
			log.SetLevel(log.DebugLevel)
		}
		env.InitConfig()

		// TODO: Passing *e.Info here is not exactly nice. Fix?
		ocp3config := ocp3.New()
		ocp3config.Fetch()

		m := ocp3config.ParseMaster()

		crd, err := oauth.Generate(m)

		if err != nil {
			log.Fatal(err)
		}

		oauth.PrintCRD(crd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// It only needs to happen once.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

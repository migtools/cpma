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
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/fusor/cpma/env"
	"github.com/fusor/cpma/ocp3"
	"github.com/fusor/cpma/ocp4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringVar(&env.ConfigFile, "config", "", "config file (default is $HOME/.cpma.yaml)")

	rootCmd.Flags().Bool("debug", false, "show debug ouput")
	env.Config().BindPFlag("Debug", rootCmd.Flags().Lookup("debug"))

	rootCmd.Flags().StringP("output-dir", "o", path.Dir(""), "set the directory to store extracted configuration.")
	env.Config().BindPFlag("OutputDir", rootCmd.Flags().Lookup("output-dir"))

	// Default timeout is 10s
	rootCmd.Flags().DurationP("timeout", "t", 10000000000, "Set timeout, unit must be provided, i.e. '-t 20s'.")
	env.Config().BindPFlag("TimeOut", rootCmd.Flags().Lookup("timeout"))
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cpma",
	Short: "Helps migration cluster configuration of a OCP 3.x cluster to OCP 4.x",
	Long:  `Helps migration cluster configuration of a OCP 3.x cluster to OCP 4.x`,
	Run: func(cmd *cobra.Command, args []string) {
		env.InitConfig()
		env.InitLogger()

		// TODO: add survey to handle UI
		ocp3config := ocp3.New()
		ocp3config.Fetch()
		mc := ocp3config.ParseMaster()
		clusterV4 := ocp4.Cluster{}
		clusterV4.Translate(mc)
		manifests := clusterV4.GenYAML()

		// TODO: Add to pipeline as exit channel
		for _, manifest := range manifests {
			maniftestfile := filepath.Join(env.Config().GetString("OutputDir"), "manifests", manifest.Name)
			os.MkdirAll(path.Dir(maniftestfile), 0755)
			err := ioutil.WriteFile(maniftestfile, manifest.CRD, 0644)
			logrus.Printf("CR manifest created: %s", maniftestfile)
			if err != nil {
				logrus.Panic(err)
			}
		}
		fmt.Println(ocp4.OCP4InstallMsg)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// It only needs to happen once.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("Something went terribly wrong!")
	}
}

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
	//Workaround go mod vendor issue 27063
	_ "github.com/shurcooL/vfsgen"

	"github.com/fusor/cpma/pkg/env"
	"github.com/fusor/cpma/pkg/transform"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize()

	// Get config file from CLI argument an save it to viper config
	rootCmd.PersistentFlags().StringVar(&env.ConfigFile, "config", "", "config file (Default searches ./cpma.yaml, $HOME/cpma.yml)")

	// Allow insecure host key if true
	rootCmd.PersistentFlags().BoolP("allow-insecure-host", "i", false, "allow insecure ssh host key ")
	env.Config().BindPFlag("InsecureHostKey", rootCmd.PersistentFlags().Lookup("allow-insecure-host"))

	// Set configuration source
	rootCmd.PersistentFlags().String("config-source", "", "source for OCP3 config files, accepted values: remote or local")
	env.Config().BindPFlag("ConfigSource", rootCmd.PersistentFlags().Lookup("config-source"))

	// Get OCP3 source cluster name that is used in kubeconfig and save it to viper config
	rootCmd.PersistentFlags().StringP("cluster-name", "c", "", "OCP3 cluster kubeconfig name")
	env.Config().BindPFlag("ClusterName", rootCmd.PersistentFlags().Lookup("cluster-name"))

	// Get crio config file location
	rootCmd.PersistentFlags().String("crio-config", "", "path to crio config file")
	env.Config().BindPFlag("CrioConfigFile", rootCmd.PersistentFlags().Lookup("crio-config"))

	// Set log level
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "show debug ouput")
	env.Config().BindPFlag("Debug", rootCmd.PersistentFlags().Lookup("debug"))

	// Target cluster
	rootCmd.PersistentFlags().Bool("target-cluster", false, "target cluster")
	env.Config().BindPFlag("TargetCluster", rootCmd.PersistentFlags().Lookup("target-cluster"))

	// Target cluster name for Kubeconfig context
	rootCmd.PersistentFlags().StringP("target-cluster-name", "t", "", "target cluster name")
	env.Config().BindPFlag("TargetClusterName", rootCmd.PersistentFlags().Lookup("target-cluster-name"))

	// Get etcd config file location
	rootCmd.PersistentFlags().String("etcd-config", "", "path to etcd config file")
	env.Config().BindPFlag("ETCDConfigFile", rootCmd.PersistentFlags().Lookup("etcd-config"))

	// Get OCP3 source cluster and save it to viper config
	rootCmd.PersistentFlags().StringP("hostname", "n", "", "OCP3 cluster hostname")
	env.Config().BindPFlag("Hostname", rootCmd.PersistentFlags().Lookup("hostname"))

	// Get node config file location
	rootCmd.PersistentFlags().String("node-config", "", "path to node config file")
	env.Config().BindPFlag("NodeConfigFile", rootCmd.PersistentFlags().Lookup("node-config"))

	// Flag to generate manifests
	rootCmd.PersistentFlags().BoolP("manifests", "m", true, "Generate manifests")
	env.Config().BindPFlag("Manifests", rootCmd.PersistentFlags().Lookup("manifests"))

	// Get master config file location
	rootCmd.PersistentFlags().String("master-config", "", "path to master config file")
	env.Config().BindPFlag("MasterConfigFile", rootCmd.PersistentFlags().Lookup("master-config"))

	// Get registries config file location
	rootCmd.PersistentFlags().String("registries-config", "", "path to registries config file")
	env.Config().BindPFlag("RegistriesConfigFile", rootCmd.PersistentFlags().Lookup("registries-config"))

	// Flag to generate reporting
	rootCmd.PersistentFlags().BoolP("reporting", "r", true, "Generate reporting ")
	env.Config().BindPFlag("Reporting", rootCmd.PersistentFlags().Lookup("reporting"))

	rootCmd.PersistentFlags().StringP("ssh-keyfile", "k", "", "OCP3 ssh keyfile path")
	env.Config().BindPFlag("SSHPrivateKey", rootCmd.PersistentFlags().Lookup("ssh-keyfile"))

	// Get ssh config values
	rootCmd.PersistentFlags().StringP("ssh-login", "l", "", "OCP3 ssh login")
	env.Config().BindPFlag("SSHLogin", rootCmd.PersistentFlags().Lookup("ssh-login"))

	rootCmd.PersistentFlags().Int16P("ssh-port", "p", 0, "OCP3 ssh port")
	env.Config().BindPFlag("SSHPort", rootCmd.PersistentFlags().Lookup("ssh-port"))

	// Don't output logs to console if true
	rootCmd.PersistentFlags().BoolP("silent", "s", false, "silent mode, disable logging output to console")
	env.Config().BindPFlag("Silent", rootCmd.PersistentFlags().Lookup("silent"))

	// Get config file from an save to viper config
	rootCmd.PersistentFlags().StringP("work-dir", "w", "", "set application data working directory (Default \".\")")
	env.Config().BindPFlag("WorkDir", rootCmd.PersistentFlags().Lookup("work-dir"))
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cpma",
	Short: "Helps migration cluster configuration of a OCP 3.x cluster to OCP 4.x",
	Long:  `Helps migration cluster configuration of a OCP 3.x cluster to OCP 4.x`,
	Run: func(cmd *cobra.Command, args []string) {
		env.InitLogger()

		if err := env.InitConfig(); err != nil {
			logrus.Fatal(err)
		}

		transform.Start()
	},
	Args: cobra.MaximumNArgs(0),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// It only needs to happen once.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
	}
}

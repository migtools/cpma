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

	//"github.com/fusor/network"

	"github.com/fusor/cpma/env"
	"github.com/fusor/cpma/ocp3"
	"github.com/fusor/cpma/ocp4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type File struct {
	URL string
}

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
// The workflow is organized using pipelines chaining queues.
// Queues are channels starting from files to manifests
var rootCmd = &cobra.Command{
	Use:   "cpma",
	Short: "Helps migration cluster configuration of a OCP 3.x cluster to OCP 4.x",
	Long:  `Helps migration cluster configuration of a OCP 3.x cluster to OCP 4.x`,
	Run: func(cmd *cobra.Command, args []string) {
		startTime := time.Now()
		env.InitConfig()
		env.InitLogger()

		filestoFetch := make(chan File)
		fetchedFiles := make(chan File)
		toTranslate := make(chan File)
		manifests := make(chan File)

		logrus.Printf("Timeout in %s (Ctrl+C to stop)", env.Config().GetDuration("TimeOut"))

		// Communication Broker
		go func() {
			for {
				file := <-filestoFetch
				fmt.Printf("broker: %v\n", file.URL)
				// use stmpclient and save file locally
				fetchedFiles <- file
			}
		}()

		// OCP4 translater
		go func() {
			for {
				file := <-toTranslate
				fmt.Printf("ocp4 translator: %v\n", file.URL)
				// Generate manifest
				manifests <- file
			}
		}()

		// OCP3 decoder
		go func() {
			for {
				file := <-fetchedFiles
				fmt.Printf("ocp3 decoder: %v\n", file.URL)
				toTranslate <- file
			}
		}()

		// Dir Monitor
		go func(dir string) {
			processedFiles := make([]File, 20)

			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
					return err
				}

				// TODO: check file has already processed and loop around Dir
				if path != dir && !info.IsDir() {
					file := File{URL: path}
					processedFiles = append(processedFiles, file)
					fmt.Printf("monitor: %v\n", path)
					filestoFetch <- file
				}
				return nil
			})
			if err != nil {
				fmt.Printf("error walking the path %q: %v\n", dir, err)
				return
			}
		}(env.Config().GetString("OutputDir"))

		for {
			elapsed := time.Now().Sub(startTime)
			select {
			case msg := <-manifests:
				fmt.Println("manifests:", msg)
			case <-time.After(env.Config().GetDuration("TimeOut") - elapsed):
				fmt.Println("timeout")
				os.Exit(0)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// It only needs to happen once.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("Something went terribly wrong!")
	}
}

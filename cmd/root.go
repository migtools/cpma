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
	"strings"
	"time"

	//"github.com/fusor/network"

	"github.com/fusor/cpma/env"
	"github.com/fusor/cpma/internal/sftpclient"
	"github.com/fusor/cpma/internal/types"
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

		filestoFetch := make(chan types.File)
		fetchedFiles := make(chan types.File)
		toTranslate := make(chan types.File)
		manifests := make(chan ocp4.Manifest)

		logrus.Printf("Timeout in %s (Ctrl+C to stop)", env.Config().GetDuration("TimeOut"))

		outputDir := env.Config().GetString("OutputDir")

		// Communication Broker
		// TODO: Ideally create a channel per host
		go func() {
			for {
				file := <-filestoFetch
				logrus.Debugf("broker hostname: %v, file: %v\n", file.Hostname, file.Path)

				dst := filepath.Join(env.Config().GetString("OutputDir"), file.Hostname, file.Path)
				src := filepath.Join("/", file.Path)
				sftpclient.Fetch(file.Hostname, src, dst)

				f, err := ioutil.ReadFile(filepath.Join(outputDir, file.Hostname, file.Path))
				if err != nil {
					logrus.Fatal(err)
				}
				file.Content = f
				fetchedFiles <- file
			}
		}()

		// OCP3 decoder
		go func() {
			for {
				file := <-fetchedFiles
				logrus.Debugf("ocp3 decoder: %v\n", file.Path)
				err := file.OCP3.Decode(filepath.Join(outputDir, file.Hostname, file.Path))
				if err == nil {
					toTranslate <- file
				}
			}
		}()

		// OCP4 translator generates manifests
		go func() {
			for {
				file := <-toTranslate
				logrus.Debugf("ocp4 translator: %v\n", file.Path)
				file.OCP4.Translate(file.OCP3)
				crds := file.OCP4.GenYAML()
				for _, crd := range crds {
					manifests <- crd
				}
			}
		}()

		// Directory Monitor
		go func(dir string) {
			processedFiles := make([]types.File, 20)

			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
					return err
				}

				fullpath := strings.Join(strings.Split(path, "/")[1:], "/")
				hostname := strings.Split(fullpath, "/")[0]
				filepath := strings.Join(strings.Split(fullpath, "/")[1:], "/")

				// TODO: check file has already processed and add wrap filepath.Walk with loop
				if path != dir && !info.IsDir() {
					file := types.File{Hostname: hostname, Path: filepath}
					processedFiles = append(processedFiles, file)
					filestoFetch <- file
				}
				return nil
			})
			if err != nil {
				fmt.Printf("Error walking the path %q: %v\n", dir, err)
				return
			}
		}(env.Config().GetString("OutputDir"))

		for {
			elapsed := time.Now().Sub(startTime)
			select {
			case manifest := <-manifests:
				logrus.Debugf(fmt.Sprintf("manifests: %s", manifest.Name))

				maniftestfile := filepath.Join("manifests", manifest.Name)
				os.MkdirAll(path.Dir(maniftestfile), 0755)
				err := ioutil.WriteFile(maniftestfile, manifest.CRD, 0644)
				logrus.Printf("CR manifest created: %s", maniftestfile)
				if err != nil {
					logrus.Panic(err)
				}

			case <-time.After(env.Config().GetDuration("TimeOut") - elapsed):
				fmt.Println("timeout")
				fmt.Println(ocp4.OCP4InstallMsg)
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

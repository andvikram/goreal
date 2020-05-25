// Copyright © 2018 Vikram Anand <vikram983@outlook.com>
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
	"strings"

	"github.com/andvikram/goreal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	env           string
	configDirPath string
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Command to run GoReal",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("ENVIRONMENT", env)
		viper.Set("CONFIGDIRPATH", configDirPath)
		server.Start(strings.ToLower(env))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(
		&env,
		"environment",
		"e", "development",
		"Specify environment to run",
	)

	runCmd.Flags().StringVarP(
		&configDirPath,
		"config-file",
		"c", homeDir,
		"Specify YAML config directory path. Defaults to Home directory.",
	)
}

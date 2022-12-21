/*
Copyright © 2021 ant1k9 <ant1k9@protonmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ant1k9/auto-launcher/cmd/auto-launcher/edit"
	"github.com/ant1k9/auto-launcher/cmd/auto-launcher/rm"
	"github.com/ant1k9/auto-launcher/internal/config"
	"github.com/ant1k9/auto-launcher/internal/pkg/discover"
	"github.com/ant1k9/auto-launcher/internal/pkg/utils"
)

// nolint: gochecknoglobals
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "auto-launcher",
	Short: "Auto discover and launch executable files",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(discover.RunFile)
		if os.IsNotExist(err) {
			err = discover.ChooseExecutable(config.GetConfig())
		}
		utils.Must(err)

		utils.Must(utils.RunCommand(
			"/usr/bin/env",
			append([]string{"bash", discover.RunFile}, args...)...),
		)
	},
}

func main() {
	rootCmd.AddCommand(edit.Cmd)
	rootCmd.AddCommand(rm.Cmd)
	cobra.CheckErr(rootCmd.Execute())
}

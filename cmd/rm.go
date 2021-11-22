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
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ant1k9/auto-launcher/internal/pkg/discover"
)

// editCmd represents the edit command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove launcher command",
	Run: func(_ *cobra.Command, args []string) {
		_ = os.Remove(discover.RunFile)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}

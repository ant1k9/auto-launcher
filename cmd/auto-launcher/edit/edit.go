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
package edit

import (
	"github.com/spf13/cobra"

	"github.com/ant1k9/auto-launcher/internal/pkg/discover"
	"github.com/ant1k9/auto-launcher/internal/pkg/utils"
)

// nolint: gochecknoglobals
// editCmd represents the edit command
var Cmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit launcher command",
	Run: func(_ *cobra.Command, args []string) {
		utils.FatalIfErr(utils.RunCommand("/usr/bin/env", "vim", discover.RunFile))
	},
}

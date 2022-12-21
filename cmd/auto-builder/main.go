/*
Copyright © 2022 ant1k9 <ant1k9@protonmail.com>
*/
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ant1k9/auto-launcher/internal/config"
	"github.com/ant1k9/auto-launcher/internal/pkg/discover"
	"github.com/ant1k9/auto-launcher/internal/pkg/utils"
	"github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

// nolint: gochecknoglobals
// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "auto-builder",
	Short: "Auto download and build executable files",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Print("provide a github link to install executable")
			return
		}

		_, name := filepath.Split(args[0])
		cloneDirectory := filepath.Join(os.TempDir(), name)

		utils.Must(os.RemoveAll(cloneDirectory))
		_, err := git.PlainClone(cloneDirectory, false, &git.CloneOptions{URL: args[0]})
		utils.Must(err)
		utils.Must(os.Chdir(cloneDirectory))

		buildCommand, err := discover.ChooseBuildCommand(name, config.GetConfig())
		utils.Must(err)

		utils.Must(utils.RunCommand(buildCommand[0], buildCommand[1:]...))

		home, err := os.UserHomeDir()
		utils.Must(err)
		utils.Must(copy.Copy(name, filepath.Join(home, "bin", name)))
	},
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}

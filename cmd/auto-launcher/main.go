package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/ant1k9/auto-launcher/internal/pkg/discover"
	"github.com/ant1k9/auto-launcher/internal/pkg/utils"
)

var (
	app  = kingpin.New("auto-launcher", "Auto discover and launch executable files")
	post = app.Command("edit", "Edit launch command")
	rm   = app.Command("rm", "Remove launch command")
)

func main() {
	// pass arguments to exec command
	var extraArgs []string
	if len(os.Args) > 1 && os.Args[1] == "--" {
		extraArgs = os.Args[2:]
		os.Args = os.Args[:1]
	}

	if len(os.Args) > 1 {
		switch kingpin.MustParse(app.Parse(os.Args[1:])) {
		case post.FullCommand():
			utils.FatalIfErr(utils.RunCommand("/usr/bin/env", "vim", discover.RunFile))
		case rm.FullCommand():
			_ = os.Remove(discover.RunFile)
		}
		return
	}

	// by default run command
	{
		_, err := os.Stat(discover.RunFile)
		if os.IsNotExist(err) {
			err = discover.ChooseExecutable()
		}
		utils.FatalIfErr(err)

		utils.FatalIfErr(utils.RunCommand(
			"/usr/bin/env",
			append([]string{"bash", discover.RunFile}, extraArgs...)...),
		)
	}
}

package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/ant1k9/auto-launcher/internal/pkg/discover"
)

func main() {
	_, err := os.Stat(".run")

	switch {
	case os.IsNotExist(err):
		err = discover.ChooseExecutable()
	}

	if err == nil {
		cmd := exec.Command("/usr/bin/env", "bash", discover.RunFile)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err = cmd.Run()
	}

	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/ant1k9/auto-launcher/internal/pkg/discover"
)

func main() {
	_, err := os.Stat(".run")
	if os.IsNotExist(err) {
		err = discover.ChooseExecutable()
	}

	if err == nil {
		cmd := exec.Command("/usr/bin/env", "bash", discover.RunFile)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		err = cmd.Run()
	}

	if err != nil {
		log.Fatal(err)
	}
}

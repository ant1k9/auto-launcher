package utils

import (
	"log"
	"os"
	"os/exec"
)

func RunCommand(exectutable string, args ...string) error {
	cmd := exec.Command(exectutable, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

func FatalIfErr(err error) {
	if err == nil {
		return
	}
	log.Fatal(err)
}

package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/gizak/termui/v3"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const RunFile = ".run"

var (
	CMainDecl    = regexp.MustCompile(`(?m)(void|int)\s+main`)
	GoMainDecl   = regexp.MustCompile(`func\s+main`)
	RustMainDecl = regexp.MustCompile(`fn\s+main`)
)

type (
	Extension = string
	Filename  = string
)

func isServicePath(path string) bool {
	switch path {
	case ".git", "target", "test":
		return true
	default:
		return false
	}
}

func hasMain(path string, mainDecl *regexp.Regexp) bool {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}
	return len(mainDecl.Find(content)) > 0
}

func isCExecutable(path string) bool    { return hasMain(path, CMainDecl) }
func isGoExecutable(path string) bool   { return hasMain(path, GoMainDecl) }
func isRustExecutable(path string) bool { return hasMain(path, RustMainDecl) }

func getExtension(path string) string {
	switch path {
	case "Makefile":
		return path
	default:
		return filepath.Ext(path)
	}
}

func isExecutable(path string, info fs.FileInfo) bool {
	if info.IsDir() {
		return false
	}

	switch getExtension(info.Name()) {
	case ".cpp", ".c":
		return isCExecutable(path)
	case ".rs":
		return isRustExecutable(path)
	case ".go":
		return isGoExecutable(path)
	case ".sh", ".py", ".js", ".mk", "Makefile":
		return true // we cannot say is it a script or a package
	default:
		return false
	}
}

func getExecutables() (map[Extension]Filename, error) {
	files := make(map[Extension]Filename)
	return files, filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if isServicePath(path) {
			return filepath.SkipDir
		}

		if isExecutable(path, info) {
			extension := getExtension(info.Name())
			if f, ok := files[getExtension(extension)]; ok {
				return fmt.Errorf(
					"several files with <%s> extension: %s, %s",
					extension, f, path,
				)
			}
			files[extension] = path
		}
		return err
	})
}

func saveExecutable(ext, path string) error {
	var command string
	switch ext {
	case ".c":
		command = "gcc -O2 -o main " + path + " && ./main"
	case ".cpp":
		command = "g++ -O2 -std=c++17 -o main " + path + " && ./main"
	case ".rs":
		command = "cargo run"
	case ".go":
		command = "go run " + path
	}

	return ioutil.WriteFile(RunFile, []byte(command), fs.ModePerm)
}

func chooseExecutable(executables map[Extension]Filename) error {
	if len(executables) == 1 {
		for ext, path := range executables {
			return saveExecutable(ext, path)
		}
	}

	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	extensions := make([]string, 0, len(executables))

	l := widgets.NewList()
	l.Title = "Choose executable to run further:\n"
	for ext, path := range executables {
		l.Rows = append(l.Rows, fmt.Sprintf("%d. %s", len(l.Rows)+1, path))
		extensions = append(extensions, ext)
	}

	l.SetRect(0, 0, 50, 10)
	l.SelectedRowStyle = termui.NewStyle(ui.ColorGreen)

	ui.Render(l)
	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return nil
		case "<Enter>":
			saveExecutable(
				extensions[l.SelectedRow],
				executables[extensions[l.SelectedRow]],
			)
			return nil
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		}
		ui.Render(l)
	}
}

func main() {
	_, err := os.Stat(".run")

	switch {
	case os.IsNotExist(err):
		var executables map[Extension]Filename
		executables, err = getExecutables()
		if err != nil {
			break
		}
		err = chooseExecutable(executables)
	}

	if err == nil {
		cmd := exec.Command("/usr/bin/env", "bash", RunFile)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err = cmd.Run()
	}

	if err != nil {
		log.Fatal(err)
	}
}

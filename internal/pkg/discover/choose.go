package discover

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"regexp"

	"github.com/ant1k9/auto-launcher/internal/config"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const RunFile = ".run"

var listedPathsPattern = regexp.MustCompile(`\d+\. (.*)`)

func prepareCommand(ext, path string) string {
	switch ext {
	case ".c":
		return "gcc -O2 -o main " + path + " && ./main"
	case ".cpp":
		return "g++ -O2 -std=c++17 -o main " + path + " && ./main"
	case ".rs":
		return "cargo run"
	case ".py":
		return "python " + path
	case ".js":
		return "node " + path
	case ".go":
		return "go run " + path
	case "Makefile", ".mk":
		return "make"
	default:
		return ""
	}
}

func saveExecutable(ext, path string) error {
	command := prepareCommand(ext, path) + " $*"
	return ioutil.WriteFile(RunFile, []byte(command), fs.ModePerm)
}

func chooseInteractive(executables map[Extension][]Filename) error {
	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	extensions := make([]string, 0, len(executables))
	l := widgets.NewList()
	l.Title = "Choose executable to run further:\n"
	for ext, paths := range executables {
		for _, path := range paths {
			l.Rows = append(l.Rows, fmt.Sprintf("%d. %s", len(l.Rows)+1, path))
			extensions = append(extensions, ext)
		}
	}

	l.SetRect(0, 0, 50, 10)
	l.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)

	ui.Render(l)
	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return nil
		case "<Enter>":
			if m := listedPathsPattern.FindStringSubmatch(l.Rows[l.SelectedRow]); len(m) > 0 {
				return saveExecutable(extensions[l.SelectedRow], m[1])
			}
			return fmt.Errorf("unexpected row data: %s", l.Rows[l.SelectedRow])
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		}
		ui.Render(l)
	}
}

func ChooseExecutable(cfg config.Config) error {
	executables, err := getExecutables(".", cfg)
	if err != nil {
		return err
	}

	if len(executables) == 1 {
		for ext, paths := range executables {
			if len(paths) == 1 {
				return saveExecutable(ext, paths[0])
			}
		}
	}

	return chooseInteractive(executables)
}

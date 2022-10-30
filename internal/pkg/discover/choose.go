package discover

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
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
		cFiles := filepath.Join(filepath.Dir(path), "*.c")
		return "gcc -O2 -o main " + cFiles + " && ./main"
	case ".cpp":
		cppFiles := filepath.Join(filepath.Dir(path), "*.cpp")
		return "g++ -O2 -std=c++17 -o main " + cppFiles + " && ./main"
	case ".rs":
		return "cargo run"
	case ".py":
		return "python " + path
	case ".js":
		return "node " + path
	case ".go":
		return "go run " + path
	case Makefile, ".mk":
		return "make"
	default:
		return ""
	}
}

func saveExecutable(ext, path string) error {
	command := prepareCommand(ext, path) + " $*"
	return ioutil.WriteFile(RunFile, []byte(command), fs.ModePerm)
}

// nolint: maintidx
func chooseInteractive(executables map[Extension][]Filename) error {
	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialize termui: %w", err)
	}
	defer ui.Close()

	extensions, widgetsList := prepareWidgetsList(executables)
	ui.Render(widgetsList)
	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return nil
		case "<Enter>":
			if m := listedPathsPattern.FindStringSubmatch(widgetsList.Rows[widgetsList.SelectedRow]); len(m) > 0 {
				return saveExecutable(extensions[widgetsList.SelectedRow], m[1])
			}
			return fmt.Errorf("unexpected row data: %s", widgetsList.Rows[widgetsList.SelectedRow])
		case "j", "<Down>":
			widgetsList.ScrollDown()
		case "k", "<Up>":
			widgetsList.ScrollUp()
		}
		ui.Render(widgetsList)
	}
}

func prepareWidgetsList(executables map[Extension][]Filename) ([]string, *widgets.List) {
	extensions := make([]string, 0, len(executables))

	widgetsList := widgets.NewList()
	widgetsList.Title = "Choose executable to run further:\n"
	for ext, paths := range executables {
		for _, path := range paths {
			widgetsList.Rows = append(widgetsList.Rows, fmt.Sprintf("%d. %s", len(widgetsList.Rows)+1, path))
			extensions = append(extensions, ext)
		}
	}

	widgetsList.SetRect(0, 0, 50, 10) // nolint: gomnd
	widgetsList.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)

	return extensions, widgetsList
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

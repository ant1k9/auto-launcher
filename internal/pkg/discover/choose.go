package discover

import (
	"fmt"
	"io/fs"
	"io/ioutil"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const RunFile = ".run"

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
	case "Makefile", ".mk":
		command = "make"
	}

	return ioutil.WriteFile(RunFile, []byte(command), fs.ModePerm)
}

func chooseInteractive(executables map[Extension]Filename) error {
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
	l.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)

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

func ChooseExecutable() error {
	executables, err := getExecutables()
	if err != nil {
		return err
	}

	if len(executables) == 1 {
		for ext, path := range executables {
			return saveExecutable(ext, path)
		}
	}

	return chooseInteractive(executables)
}

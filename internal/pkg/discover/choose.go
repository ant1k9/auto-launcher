package discover

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ant1k9/auto-launcher/internal/config"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	listedPathsPattern = regexp.MustCompile(`\d+\. (.*)`)

	ErrCommandNotFound = errors.New("command not found")
)

// nolint: maintidx
func prepareCommand(ext, path string) (string, error) {
	switch ext {
	case CExtension:
		cFiles := filepath.Join(filepath.Dir(path), "*.c")
		return "gcc -O2 -o main " + cFiles + " && ./main " + BashArgs, nil
	case CPPExtension:
		cppFiles := filepath.Join(filepath.Dir(path), "*.cpp")
		return "g++ -O2 -std=c++17 -o main " + cppFiles + " && ./main " + BashArgs, nil
	case RustExtension:
		return "cargo run " + BashArgs, nil
	case PythonExtension:
		return fmt.Sprintf("python %s %s", path, BashArgs), nil
	case JavaScriptExtension:
		return fmt.Sprintf("node %s %s", path, BashArgs), nil
	case GoExtension:
		return fmt.Sprintf("go run %s %s", path, BashArgs), nil
	case BashExtension:
		return fmt.Sprintf("bash %s %s", path, BashArgs), nil
	case FishExtension:
		return fmt.Sprintf("fish %s %s", path, BashArgs), nil
	case Makefile, MakeExtension:
		return "make " + BashArgs, nil
	case Dockerfile:
		return prepareDockerCommand()
	default:
		return "", ErrCommandNotFound
	}
}

func prepareBuildCommand(ext, path, name string) ([]string, error) {
	switch ext {
	case RustExtension:
		return prepareRustBuildCommand(path)
	case GoExtension:
		return prepareGoBuildCommand(path, name)
	case Makefile, MakeExtension:
		return []string{"make", "install"}, nil
	default:
		return nil, ErrCommandNotFound
	}
}

func prepareGoBuildCommand(path, name string) ([]string, error) {
	dir, _ := filepath.Split(path)
	switch {
	case len(dir) > 1:
		dir = "./" + filepath.Join(dir, "...")
	default:
		dir = "."
	}
	return []string{"go", "build", "-o", name, dir}, nil
}

func prepareRustBuildCommand(path string) ([]string, error) {
	pathParts := strings.Split(path, string(os.PathSeparator))
	for idx := range pathParts {
		fmt.Println(idx, pathParts[idx])
		if pathParts[idx] == "src" {
			return []string{"cargo", "install", "--path", filepath.Join(pathParts[:idx]...)}, nil
		}
	}
	return []string{"cargo", "install", "--path", "."}, nil
}

func prepareDockerCommand() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	baseDir := filepath.Base(dir)
	return fmt.Sprintf(
		"docker build -t %[1]s:local .\ndocker run --rm -ti %[2]s %[1]s:local",
		baseDir, BashArgs,
	), nil
}

func saveExecutable(ext, path string) error {
	command, err := prepareCommand(ext, path)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(RunFile, []byte(command), fs.ModePerm)
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

	_, err = chooseExecutableInteractively(executables)
	return err
}

func chooseExecutableInteractively(executables map[Extension][]Filename) (any, error) {
	return chooseInteractively(executables, func(ext, path string) (any, error) {
		return nil, saveExecutable(ext, path)
	})
}

func ChooseBuildCommand(name string, cfg config.Config) ([]string, error) {
	executables, err := getBuildExecutables(".", cfg)
	if err != nil {
		return nil, err
	}

	if len(executables) == 1 {
		for ext, paths := range executables {
			if len(paths) == 1 {
				return prepareBuildCommand(ext, paths[0], name)
			}
		}
	}

	return chooseBuilderInteractively(executables, name)
}

func chooseBuilderInteractively(executables map[Extension][]Filename, name string) ([]string, error) {
	return chooseInteractively(executables, func(ext, path string) ([]string, error) {
		return prepareBuildCommand(ext, path, name)
	})
}

// nolint: maintidx
func chooseInteractively[T any](
	executables map[Extension][]Filename,
	resultFn func(ext, path string) (T, error),
) (result T, err error) {
	if err = ui.Init(); err != nil {
		return result, fmt.Errorf("failed to initialize termui: %w", err)
	}
	defer ui.Close()

	extensions, widgetsList := prepareWidgetsList(executables)
	ui.Render(widgetsList)
	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return result, nil
		case "<Enter>":
			if m := listedPathsPattern.FindStringSubmatch(widgetsList.Rows[widgetsList.SelectedRow]); len(m) > 0 {
				return resultFn(extensions[widgetsList.SelectedRow], m[1])
			}
			return result, fmt.Errorf(
				"unexpected row data: %s", widgetsList.Rows[widgetsList.SelectedRow],
			)
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

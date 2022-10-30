package discover

import (
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/ant1k9/auto-launcher/internal/config"
)

var (
	CMainDecl      = regexp.MustCompile(`(?m)(void|int)\s+main`)
	GoMainDecl     = regexp.MustCompile(`func\s+main`)
	PythonMainDecl = regexp.MustCompile(`if\s+__name__\s*==\s*["']__main__["']`)
	RustMainDecl   = regexp.MustCompile(`fn\s+main`)
)

type (
	Extension = string
	Filename  = string
)

func isSkippedPath(path string, cfg config.Config) bool {
	for _, skipPath := range cfg.SkipPaths {
		if filepath.Base(path) == skipPath {
			return true
		}
	}
	return false
}

func hasMain(path string, mainDecl *regexp.Regexp) bool {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}
	return len(mainDecl.Find(content)) > 0
}

func isCExecutable(path string) bool      { return hasMain(path, CMainDecl) }
func isGoExecutable(path string) bool     { return hasMain(path, GoMainDecl) }
func isPythonExecutable(path string) bool { return hasMain(path, PythonMainDecl) }
func isRustExecutable(path string) bool   { return hasMain(path, RustMainDecl) }

func getExtension(path string) string {
	switch path {
	case Makefile, Dockerfile:
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
	case ".py":
		return isPythonExecutable(path)
	case ".fish", ".sh", ".js", ".mk", Makefile, Dockerfile:
		return true // we cannot say is it a script or a package
	default:
		return false
	}
}

func getExecutables(root string, cfg config.Config) (map[Extension][]Filename, error) {
	files := make(map[Extension][]Filename)
	return files, filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if isSkippedPath(path, cfg) {
			return filepath.SkipDir
		}

		if isExecutable(path, info) {
			extension := getExtension(info.Name())
			files[extension] = append(files[extension], path)
		}
		return err
	})
}

package discover

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

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
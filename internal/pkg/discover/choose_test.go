package discover

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/ant1k9/auto-launcher/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrepareCommand(t *testing.T) {
	dir, _ := os.Getwd()
	baseDir := filepath.Base(dir)

	tests := []struct {
		name string
		ext  string
		path string
		want string
		err  error
	}{
		{
			name: "go command",
			ext:  ".go",
			path: "main.go",
			want: "go run main.go $*",
		},
		{
			name: "rust command",
			ext:  ".rs",
			path: "main.rs",
			want: "cargo run $*",
		},
		{
			name: "c++ command",
			ext:  ".cpp",
			path: "main.cpp",
			want: "g++ -O2 -std=c++17 -o main *.cpp && ./main $*",
		},
		{
			name: "c command",
			ext:  ".c",
			path: "main.c",
			want: "gcc -O2 -o main *.c && ./main $*",
		},
		{
			name: "Makefile command",
			ext:  Makefile,
			path: Makefile,
			want: "make $*",
		},
		{
			name: "python command",
			ext:  ".py",
			path: "main.py",
			want: "python main.py $*",
		},
		{
			name: "javascript command",
			ext:  ".js",
			path: "main.js",
			want: "node main.js $*",
		},
		{
			name: "docker command",
			ext:  Dockerfile,
			path: Dockerfile,
			want: fmt.Sprintf(
				"docker build -t %[1]s:local .\ndocker run --rm -ti $* %[1]s:local",
				baseDir,
			),
		},
		{
			name: "bash command",
			ext:  ".sh",
			path: "main.sh",
			want: "bash main.sh $*",
		},
		{
			name: "fish command",
			ext:  ".fish",
			path: "main.fish",
			want: "fish main.fish $*",
		},
		{
			name: "Unknown command",
			ext:  ".abc",
			path: "file.abc",
			err:  ErrCommandNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prepareCommand(tt.ext, tt.path)
			if got != tt.want {
				t.Errorf("prepareCommand() = %v, want %v", got, tt.want)
			}
			if tt.err != err {
				t.Errorf("error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestPrepareBuildCommand(t *testing.T) {
	tests := []struct {
		name       string
		ext        string
		path       string
		executable string
		want       []string
		err        error
	}{
		{
			name:       "go command",
			ext:        ".go",
			path:       "cmd/main.go",
			executable: "executable",
			want:       []string{"go", "build", "-o", "executable", "./cmd/..."},
		},
		{
			name:       "go command nested",
			ext:        ".go",
			path:       "main.go",
			executable: "executable",
			want:       []string{"go", "build", "-o", "executable", "."},
		},
		{
			name:       "rust command",
			ext:        ".rs",
			path:       "main.rs",
			executable: "executable",
			want:       []string{"cargo", "install", "--path", "."},
		},
		{
			name:       "Makefile command",
			ext:        Makefile,
			path:       Makefile,
			executable: "executable",
			want:       []string{"make", "install"},
		},
		{
			name:       "Unknown command",
			ext:        ".abc",
			path:       "file.abc",
			executable: "executable",
			err:        ErrCommandNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prepareBuildCommand(tt.ext, tt.path, tt.executable)
			assert.EqualValues(t, tt.want, got)

			if tt.err != err {
				t.Errorf("error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestChooseExecutable(t *testing.T) {
	tests := []struct {
		name        string
		genFilename string
		genContent  string
		want        string
	}{
		{
			name:        "go executable",
			genFilename: "main.go",
			genContent: `
package main

func main()	{}
`,
			want: `go run main.go $*`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootPath, err := ioutil.TempDir("/tmp", "discover-test")
			require.NoError(t, err)
			defer os.RemoveAll(rootPath)

			err = os.Mkdir(path.Join(rootPath, "test"), 0755)
			require.NoError(t, err)

			require.NoError(t, ioutil.WriteFile(
				path.Join(rootPath, tt.genFilename),
				[]byte(tt.genContent),
				fs.ModePerm,
			))

			require.NoError(t, os.Chdir(rootPath))

			err = ChooseExecutable(config.Config{})
			require.NoError(t, err)

			content, err := ioutil.ReadFile(".run")
			require.NoError(t, err)
			assert.EqualValues(t, tt.want, string(content))
		})
	}
}

func TestChooseBuildCommand(t *testing.T) {
	tests := []struct {
		name        string
		genFilename string
		genContent  string
		executable  string
		want        []string
	}{
		{
			name:        "go executable",
			genFilename: "main.go",
			genContent: `
package main

func main()	{}
`,
			executable: "executable",
			want:       []string{"go", "build", "-o", "executable", "."},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootPath, err := ioutil.TempDir("/tmp", "discover-test")
			require.NoError(t, err)
			defer os.RemoveAll(rootPath)

			err = os.Mkdir(path.Join(rootPath, "test"), 0755)
			require.NoError(t, err)

			require.NoError(t, ioutil.WriteFile(
				path.Join(rootPath, tt.genFilename),
				[]byte(tt.genContent),
				fs.ModePerm,
			))

			require.NoError(t, os.Chdir(rootPath))

			buildCommand, err := ChooseBuildCommand(tt.executable, config.Config{})
			require.NoError(t, err)

			assert.EqualValues(t, tt.want, buildCommand)
		})
	}
}

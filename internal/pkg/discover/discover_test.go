package discover

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/ant1k9/auto-launcher/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getExecutables(t *testing.T) {
	tests := []struct {
		name        string
		genFilename string
		genContent  string
		want        map[Extension]Filename
	}{
		{
			name:        "go executable",
			genFilename: "main.go",
			genContent: `
package main

func main()	{}
`,
			want: map[string]string{
				".go": "main.go",
			},
		},
		{
			name:        "go executable no main",
			genFilename: "lib.go",
			genContent: `
package main

func init()	{}
`,
			want: map[string]string{},
		},
		{
			name:        "rust executable",
			genFilename: "main.rs",
			genContent: `
fn main() {
}
`,
			want: map[string]string{
				".rs": "main.rs",
			},
		},
		{
			name:        "c++ executable",
			genFilename: "main.cpp",
			genContent: `
int main() {
}
`,
			want: map[string]string{
				".cpp": "main.cpp",
			},
		},
		{
			name:        "bash executable",
			genFilename: "script.sh",
			genContent: `
echo
`,
			want: map[string]string{
				".sh": "script.sh",
			},
		},
		{
			name:        "Makefile",
			genFilename: "Makefile",
			genContent: `
.PHONY: all
all:
`,
			want: map[string]string{
				"Makefile": "Makefile",
			},
		},
		{
			name:        "python executable",
			genFilename: "exec.py",
			genContent: `
if __name__ == "__main__":
	print("Hello")
`,
			want: map[string]string{
				".py": "exec.py",
			},
		},
		{
			name:        "no executables",
			genFilename: "file.txt",
			genContent:  `Hello world!`,
			want:        map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootPath, err := ioutil.TempDir("/tmp", "discover-test")
			require.NoError(t, err)
			defer os.RemoveAll(rootPath)

			err = os.Mkdir(path.Join(rootPath, ".git"), 0755)
			require.NoError(t, err)

			require.NoError(t, ioutil.WriteFile(
				path.Join(rootPath, tt.genFilename),
				[]byte(tt.genContent),
				fs.ModePerm,
			))

			got, err := getExecutables(rootPath, config.Config{})
			require.NoError(t, err)

			assert.Len(t, got, len(tt.want))
			for k, v := range tt.want {
				assert.Equal(t, path.Join(rootPath, v), got[k])
			}
		})
	}
}

func Test_skipPaths(t *testing.T) {
	tests := []struct {
		name        string
		genFilename string
		genContent  string
		want        map[Extension]Filename
	}{
		{
			name:        "go executable",
			genFilename: "main.go",
			genContent: `
package main

func main()	{}
`,
			want: map[string]string{
				".go": "main.go",
			},
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

			// write to a skip path by confg
			require.NoError(t, ioutil.WriteFile(
				path.Join(rootPath, "test", tt.genFilename),
				[]byte(tt.genContent),
				fs.ModePerm,
			))

			got, err := getExecutables(rootPath, config.Config{
				SkipPaths: []string{"test"},
			})
			require.NoError(t, err)

			assert.Len(t, got, len(tt.want))
			for k, v := range tt.want {
				assert.Equal(t, path.Join(rootPath, v), got[k])
			}
		})
	}
}

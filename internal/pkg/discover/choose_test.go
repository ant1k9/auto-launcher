package discover

import "testing"

func Test_prepareCommand(t *testing.T) {
	tests := []struct {
		name string
		ext  string
		path string
		want string
	}{
		{
			name: "go command",
			ext:  ".go",
			path: "main.go",
			want: "go run main.go",
		},
		{
			name: "rust command",
			ext:  ".rs",
			path: "main.rs",
			want: "cargo run",
		},
		{
			name: "c++ command",
			ext:  ".cpp",
			path: "main.cpp",
			want: "g++ -O2 -std=c++17 -o main main.cpp && ./main",
		},
		{
			name: "c command",
			ext:  ".c",
			path: "main.c",
			want: "gcc -O2 -o main main.c && ./main",
		},
		{
			name: "Makefile command",
			ext:  "Makefile",
			path: "Makefile",
			want: "make",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prepareCommand(tt.ext, tt.path); got != tt.want {
				t.Errorf("prepareCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

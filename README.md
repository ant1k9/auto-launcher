![CI](https://github.com/ant1k9/auto-launcher/workflows/test/badge.svg)
[![codecov](https://codecov.io/gh/ant1k9/auto-launcher/branch/main/graph/badge.svg)](https://codecov.io/gh/ant1k9/auto-launcher)

### Installation

```bash
$ go get github.com/ant1k9/auto-launcher/cmd/...
```

### Usage

🔎 The utility tries to find a file to (compile &) execute in the current folder. If it finds one, it saves a launch command to the _.run_ file. Then you can use the auto-launcher to execute it any time with the base params.

```bash
$ auto-launcher
$ auto-launcher help
$ auto-launcher edit  # edit .run file to edit command or params
$ auto-launcher rm    # rm .run file
$ AUTO_LAUNCHER_CONFIG_PATH=config.example.toml auto-launcher  # use custom config
```

### Supported formats

 - Go
 - Rust
 - C/C++
 - Python
 - Makefile
 - JavaScript
 - Bash, Fish

package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

const envConfigPath = "AUTO_LAUNCHER_CONFIG_PATH"

type (
	Config struct {
		// SkipPaths is a list of directories where auto-launcher will not
		// search for executables
		SkipPaths []string `toml:"skip_paths"`
	}
)

func GetConfig() Config {
	var cfg Config
	if path := os.Getenv(envConfigPath); path != "" {
		if _, err := toml.DecodeFile(path, &cfg); err != nil {
			log.Printf("error parsing config file: %s", err)
		}
		return cfg
	}

	return defaultConfig()
}

func defaultConfig() Config {
	return Config{
		SkipPaths: []string{
			".git",
			"test",
			"target",
			".ccls",
		},
	}
}

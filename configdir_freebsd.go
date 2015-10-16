package main

import (
	"os"
	"os/user"
	"path/filepath"
)

func configDir() string {
	confDir := os.Getenv("XDG_CONFIG_HOME")
	if confDir != "" {
		return confDir
	}

	usr, err := user.Current()
	if err != nil {
		return ""
	}

	return filepath.Join(usr.HomeDir, ".config")
}

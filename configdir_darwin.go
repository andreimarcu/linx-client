package main

import (
	"os/user"
	"path/filepath"
)

func configDir() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}

	return filepath.Join(usr.HomeDir, ".config")
}

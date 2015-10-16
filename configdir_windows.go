package main

import (
	"os"
	"os/user"
)

func configDir() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}

	return usr.HomeDir
}

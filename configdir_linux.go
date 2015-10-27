package main

import (
	"os"
	"path/filepath"
)

func getHomeDir() (homeDir string) {
	homeDir = os.Getenv("HOME")
	if homeDir == "" {
		homeDir = getInput("Path to home directory", false)
	}

	return
}

func getConfigDir() (configDir string) {
	configDir = filepath.Join(getHomeDir(), ".config")

	return
}

package main

import (
	"os"
)

func getHomeDir() (homeDir string) {
	homeDir = os.Getenv("HOMEPATH")
	if homeDir == "" {
		homeDir = getInput("Path to home directory", false)
	}

	return
}

func getConfigDir() (configDir string) {
	configDir = os.Getenv("APPDATA")
	if configDir == "" {
		configDir = getInput("Path to config directory", false)
	}
	return
}

package main

import "fmt"

// formatVersion returns the version string in the format:
// <semver> (<commit>, <date>)
// Cobra adds the "kfg version" prefix automatically.
func formatVersion() string {
	return fmt.Sprintf("%s (%s, %s)", version, commit, date)
}
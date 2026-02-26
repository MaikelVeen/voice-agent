package main

import (
	"github.com/MaikelVeen/voice-agent/cmd"
	"github.com/MaikelVeen/voice-agent/internal/version"
)

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func main() {
	version.Version = Version
	version.Commit = Commit
	version.Date = Date

	version.ValidateVersion()

	cmd.Execute()
}

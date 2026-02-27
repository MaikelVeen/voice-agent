package main

import (
	"github.com/MaikelVeen/voice-agent/cmd"
	"github.com/MaikelVeen/voice-agent/internal/version"
	"github.com/joho/godotenv"
)

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func main() {
	_ = godotenv.Load()

	version.Version = Version
	version.Commit = Commit
	version.Date = Date

	version.ValidateVersion()

	cmd.Execute()
}

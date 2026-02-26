package version

import (
	"fmt"
	"log"
	"runtime"

	goversion "github.com/hashicorp/go-version"
)

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
}

func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func (i Info) String() string {
	return fmt.Sprintf("%s\ncommit: %s, built: %s, go: %s, platform: %s",
		i.Version, i.Commit, i.Date, i.GoVersion, i.Platform)
}

func ValidateVersion() {
	if Version == "dev" {
		return
	}

	v, err := goversion.NewVersion(Version)
	if err != nil {
		log.Fatalf("Invalid semantic version '%s': %v", Version, err)
	}

	normalized := v.String()
	if Version != normalized && Version != "v"+normalized {
		log.Fatalf("Version '%s' is not in proper semantic version format. Expected: %s or %s",
			Version, normalized, "v"+normalized)
	}
}

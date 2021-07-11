package version

import (
	"encoding/json"
	"runtime"
)

// Version is loaded via  LDFLAGS:
// VERSION := `git fetch --tags && git tag | sort -V | tail -1`
// LDFLAGS=-ldflags "-X=github.com/groupe-edf/watchdog/internal/version.Version=$(VERSION)"
var (
	BuildDate = ""
	GitCommit = ""
	Commit    = ""
	GoVersion = ""
	Version   = "develop"
)

// BuildInfo build informations
type BuildInfo struct {
	BuildDate string `json:"build_date"`
	GitCommit string `json:"git_commit"`
	GoVersion string `json:"go_version"`
	Commit    string `json:"commit"`
	Platform  string `json:"platform"`
	Version   string `json:"version"`
}

// ToJSON return build info in JSON format
func (buildInfo *BuildInfo) ToJSON() []byte {
	version, _ := json.Marshal(buildInfo)
	return version
}

// GetBuildInfo get CLI information
func GetBuildInfo() *BuildInfo {
	return &BuildInfo{
		BuildDate: BuildDate,
		GitCommit: GitCommit,
		GoVersion: runtime.Version(),
		Commit:    Commit,
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
		Version:   Version,
	}
}

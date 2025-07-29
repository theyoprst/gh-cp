package version

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

const version = "0.1"

type versionInfo struct {
	Version       string
	GitCommit     string
	CommitTime    string
	ModuleVersion string
	ModulePath    string
	GoVersion     string
}

func getVersionInfo() versionInfo {
	info := versionInfo{
		Version: version,
	}

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return info
	}

	info.GoVersion = buildInfo.GoVersion
	info.ModulePath = buildInfo.Main.Path
	if buildInfo.Main.Version != "" && buildInfo.Main.Version != "(devel)" {
		info.ModuleVersion = buildInfo.Main.Version
	}

	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs.revision":
			info.GitCommit = setting.Value
		case "vcs.time":
			info.CommitTime = setting.Value
		}
	}

	return info
}

func PrintVersion() {
	info := getVersionInfo()

	fmt.Printf("gh-cp v%s\n", info.Version)

	if info.GitCommit != "" && info.CommitTime != "" {
		commitHashShort := info.GitCommit
		if len(info.GitCommit) >= 7 {
			commitHashShort = info.GitCommit[:7]
		}

		commitTime := info.CommitTime
		if t, err := time.Parse(time.RFC3339, info.CommitTime); err == nil {
			commitTime = t.Format("2006-01-02T15:04:05Z")
		}
		fmt.Printf("commit %s %s\n", commitHashShort, commitTime)
	}

	if info.ModuleVersion != "" && info.ModulePath != "" {
		fmt.Printf("module %s %s\n", info.ModulePath, info.ModuleVersion)
	}

	if info.GoVersion != "" {
		fmt.Printf("built with go %s\n", strings.TrimPrefix(info.GoVersion, "go"))
	}
}

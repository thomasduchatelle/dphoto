// Package meta expose information about the application like the version.
package meta

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

const defaultSemver = "devel"

var (
	SemVer   = ""     // SemVer is the target version, set during build with ld-flags
	Snapshot = "true" // Snapshot is overridden by ld-flags, set to true, it will show the GIT status from where the bin has been built
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if ok && info.Main.Version != "" {
		BuildVersion.Version = info.Main.Version
	}

	for _, kv := range info.Settings {
		if kv.Value == "" {
			continue
		}
		switch kv.Key {
		case "vcs.revision":
			BuildVersion.Revision = kv.Value
		case "vcs.time":
			BuildVersion.LastCommit, _ = time.Parse(time.RFC3339, kv.Value)
		case "vcs.modified":
			BuildVersion.DirtyBuild = kv.Value == "true"
		}
	}
}

type BuildStat struct {
	Version    string
	Revision   string
	LastCommit time.Time
	DirtyBuild bool
}

var (
	BuildVersion = BuildStat{}
)

// Version returns the version of the app, updated during build
func Version() string {
	version := VersionOrDefault(BuildVersion.Version, SemVer)

	if Snapshot == "true" {
		return fmt.Sprintf("%s%s", version, BuildVersion.GitInfoParts("-"))
	}
	return version
}

func VersionOrDefault(versions ...string) string {
	for _, version := range versions {
		if version != "" && version != "unknown" && version != "(devel)" {
			return version
		}
	}

	return defaultSemver
}

func (b BuildStat) String() string {
	return fmt.Sprintf("%s %s %s", b.Version, b.GitInfoParts("#"), b.LastCommit)
}

func (b BuildStat) GitInfoParts(prefix string) string {
	if b.Revision != "unknown" && b.Revision != "" {
		parts := make([]string, 0, 3)

		commit := b.Revision
		if len(commit) > 7 {
			commit = commit[:7]
		}
		parts = append(parts, prefix)
		parts = append(parts, commit)
		if b.DirtyBuild {
			parts = append(parts, "*")
		}

		return strings.Join(parts, "")

	} else {
		return ""
	}

}

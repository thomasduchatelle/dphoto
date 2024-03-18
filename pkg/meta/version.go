// Package meta expose information about the application like the version.
package meta

import (
	"runtime/debug"
	"strings"
	"time"
)

const (
	version = "1.5.2"
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
	BuildVersion = BuildStat{
		Version:  "",
		Revision: "",
	}
)

// Version returns the version of the app, updated by ci/pre-release.sh
func Version() string {
	return version
}

func (b BuildStat) String() string {
	parts := make([]string, 0, 3)
	if b.Version != "unknown" && b.Version != "(devel)" {
		parts = append(parts, b.Version)
	}
	if b.Revision != "unknown" && b.Revision != "" {
		parts = append(parts, "#")
		commit := b.Revision
		if len(commit) > 7 {
			commit = commit[:7]
		}
		parts = append(parts, commit)
		if b.DirtyBuild {
			parts = append(parts, "*")
		}
	}
	if len(parts) == 0 {
		return "devel"
	}
	return strings.Join(parts, "")
}

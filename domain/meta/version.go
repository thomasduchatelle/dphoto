package version

const (
	version = "1.4.0-delta"
)

// Version returns the version of the app, updated by ci/pre-release.sh
func Version() string {
	return version
}

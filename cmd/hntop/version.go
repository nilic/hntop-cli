package main

var (
	version string // set by build process
)

func getVersion() string {
	if version != "" {
		return "v" + version
	}
	return "v0.0.0"
}

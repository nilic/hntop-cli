package main

import "strings"

var (
	version      string // set by build process
	appName      = "hntop"
	appNameUpper = strings.ToUpper(appName)
)

func getVersion() string {
	if version != "" {
		return "v" + version
	}
	return "v0.0.0"
}

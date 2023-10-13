package main

import (
	"log"
	"os"
	"strings"
)

const (
	defaultResultCount = 20
	minResultCount     = 1
	maxResultCount     = 1000                        // maximum number of results returned by Algolia API
	defaultTags        = "story,poll,show_hn,ask_hn" // return all post types
)

var (
	appName          = "hntop"
	appNameUpper     = strings.ToUpper(appName)
	availableTags    = []string{"story", "poll", "show_hn", "ask_hn"}
	availableOutputs = []string{"list", "mail"}
)

func main() {
	app := newApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

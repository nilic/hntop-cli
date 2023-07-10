package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
)

const (
	defaultResultCount = 20
	minResultCount     = 1
	maxResultCount     = 1000                        // maximum number of results returned by Algolia API
	defaultTag         = "story,poll,show_hn,ask_hn" // return all post types
)

var (
	tags = []string{"story", "poll", "show_hn", "ask_hn"}
)

func main() {

	app := &cli.App{
		Name:  appName,
		Usage: "display top Hacker News posts in a given time range",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "last",
				Aliases: []string{"l"},
				EnvVars: []string{appNameUpper + "_LAST"},
				Usage:   fmt.Sprintf("interval since current time to show top HN posts from, eg. \"12h\" (last 12 hours), \"100d\" (last 100 days), \"6m\" (last 6 months)\nfollowing units are supported: %s", printUnits(intervals)),
				Action: func(cCtx *cli.Context, s string) error {
					if len(s) == 1 {
						return fmt.Errorf("interval too short, needs to be in format <number><unit>, eg. 12h for 12 hours or 6m for 6 months")
					}
					length := s[:len(s)-1]
					if _, err := strconv.Atoi(length); err != nil {
						return fmt.Errorf("interval length error, needs to be in format <number><unit>, eg. 12h for 12 hours or 6m for 6 months")
					}
					last := s[len(s)-1:]
					units := getIntervalUnits(intervals)
					if !slices.Contains(units, last) {
						return fmt.Errorf("interval unit error, must end in one of %v", units)
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "from",
				EnvVars: []string{appNameUpper + "_FROM"},
				Usage:   "start of the time range to show top HN posts from in RFC3339 format \"yyyy-MM-dd'T'HH:mm:ss'Z'\" (for UTC) or \"yyyy-MM-dd'T'HH:mm:ss±hh:mm\" (for a specific timezone, ±hh:mm is the offset to UTC)\nexamples: \"2006-01-02T15:04:05Z\" (UTC time) and \"2006-01-02T15:04:05+01:00\" (CET)",
				Action: func(cCtx *cli.Context, s string) error {
					_, err := time.Parse(time.RFC3339, s)
					if err != nil {
						cli.ShowAppHelp(cCtx)
						fmt.Println()
						return fmt.Errorf("invalid time format for start of the time range, please use RFC3339 (see help for more information)")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "to",
				EnvVars: []string{appNameUpper + "_TO"},
				Usage:   "end of the time range to show top HN posts from in RFC3339 format \"yyyy-MM-dd'T'HH:mm:ss'Z'\" (for UTC) or \"yyyy-MM-dd'T'HH:mm:ss±hh:mm\" (for a specific timezone, ±hh:mm is the offset to UTC); used in conjuction with --from; if omitted, current time will be used\nexamples: \"2006-01-02T15:04:05Z\" (UTC time) and \"2006-01-02T15:04:05+01:00\" (CET)",
				Action: func(cCtx *cli.Context, s string) error {
					_, err := time.Parse(time.RFC3339, s)
					if err != nil {
						cli.ShowAppHelp(cCtx)
						fmt.Println()
						return fmt.Errorf("invalid time format for end of the time range, please use RFC3339 (see help for more information)")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "tags",
				Aliases: []string{"t"},
				EnvVars: []string{appNameUpper + "_TAG"},
				Value:   defaultTag,
				Usage:   fmt.Sprintf("filter results by item tag; available tags: %v; multiple tags can be combined with a comma, eg. show_hn,poll", tags),
				Action: func(cCtx *cli.Context, s string) error {
					// TODO
					return nil
				},
			},
			&cli.IntFlag{ // TODO
				Name:    "count",
				Aliases: []string{"c"},
				EnvVars: []string{appNameUpper + "_COUNT"},
				Value:   defaultResultCount,
				Usage:   fmt.Sprintf("number of results to retrieve, must be between %d and %d", minResultCount, maxResultCount),
				Action: func(cCtx *cli.Context, i int) error {
					if i < minResultCount || i > maxResultCount {
						return fmt.Errorf("count should be between %d and %d", minResultCount, maxResultCount)
					}
					return nil
				},
			},
			&cli.BoolFlag{ // TODO
				Name:    "front-page",
				Aliases: []string{"f"},
				EnvVars: []string{appNameUpper + "_FRONT_PAGE"},
				Usage: "display current front page posts; " +
					"have in mind that the results will be sorted differently than they appear on the front page, by points, then number of comments; " +
					"if selected, all other flags are ignored",
			},
		},
		CommandNotFound: func(cCtx *cli.Context, command string) { // TODO
			fmt.Printf("No matching command '%s'", command)
			cli.ShowAppHelp(cCtx)
		},
		Action: func(cCtx *cli.Context) error {
			err := Execute(cCtx)
			if err != nil {
				return err
			}
			return nil
		},
		Version: getVersion(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

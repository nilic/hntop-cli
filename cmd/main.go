package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
)

const (
	defaultResultCount = 20
	minResultCount     = 1
	maxResultCount     = 1000                        // maximum number of results returned by Algolia API
	defaultTags        = "story,poll,show_hn,ask_hn" // return all post types
)

var (
	appName       = "hntop"
	appNameUpper  = strings.ToUpper(appName)
	availableTags = []string{"story", "poll", "show_hn", "ask_hn"}
)

func main() {

	app := &cli.App{
		Name:  appName,
		Usage: "display top Hacker News posts",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "last",
				Aliases: []string{"l"},
				EnvVars: []string{appNameUpper + "_LAST"},
				Usage:   "Interval since current time to show top HN posts from, eg. \"12h\" (last 12 hours), \"6m\" (last 6 months).",
				Action: func(cCtx *cli.Context, s string) error {
					if len(s) == 1 {
						return fmt.Errorf("interval too short, needs to be in format <number><unit>, eg. 12h for 12 hours or 6m for 6 months")
					}
					length := s[:len(s)-1]
					if _, err := strconv.Atoi(length); err != nil {
						return fmt.Errorf("invalid interval length, needs to be in format <number><unit>, eg. 12h for 12 hours or 6m for 6 months")
					}
					last := s[len(s)-1:]
					units := getIntervalUnits(intervals)
					if !slices.Contains(units, last) {
						return fmt.Errorf("invalid interval unit, must end in one of %v", units)
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "from",
				EnvVars: []string{appNameUpper + "_FROM"},
				Usage:   "Start of the time range to show top HN posts from in RFC3339 format.",
				Action: func(cCtx *cli.Context, s string) error {
					_, err := time.Parse(time.RFC3339, s)
					if err != nil {
						return fmt.Errorf("invalid time format for start of the time range, please use RFC3339")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "to",
				EnvVars: []string{appNameUpper + "_TO"},
				Usage:   "End of the time range to show top HN posts from in RFC3339 format. Used in conjuction with --from. If omitted, current time will be used.",
				Action: func(cCtx *cli.Context, s string) error {
					if cCtx.String("from") == "" {
						return fmt.Errorf("start of the time range missing, please use --from <value in RFC3339> to specify")
					}
					to, err := time.Parse(time.RFC3339, s)
					if err != nil {
						return fmt.Errorf("invalid time format for end of the time range, please use RFC3339")
					}
					from, _ := time.Parse(time.RFC3339, cCtx.String("from"))
					if !to.After(from) {
						return fmt.Errorf("end of the time range should be later than start of the time range")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "tags",
				Aliases: []string{"t"},
				EnvVars: []string{appNameUpper + "_TAGS"},
				Value:   defaultTags,
				Usage:   fmt.Sprintf("Filter results by post tag. Available tags: %v.", availableTags),
				Action: func(cCtx *cli.Context, s string) error {
					tags := strings.Split(s, ",")
					for _, t := range tags {
						if !slices.Contains(availableTags, t) {
							return fmt.Errorf("invalid tag value \"%s\", available tags: %v", t, availableTags)
						}
					}
					return nil
				},
			},
			&cli.IntFlag{
				Name:    "count",
				Aliases: []string{"c"},
				EnvVars: []string{appNameUpper + "_COUNT"},
				Value:   defaultResultCount,
				Usage:   fmt.Sprintf("Number of results to retrieve, must be between %d and %d.", minResultCount, maxResultCount),
				Action: func(cCtx *cli.Context, i int) error {
					if i < minResultCount || i > maxResultCount {
						return fmt.Errorf("count must be between %d and %d", minResultCount, maxResultCount)
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "front-page",
				Aliases: []string{"f"},
				EnvVars: []string{appNameUpper + "_FRONT_PAGE"},
				Usage:   "Display current front page posts. If selected, all other flags are ignored.",
			},
			&cli.StringFlag{
				Name:     "mail-from",
				EnvVars:  []string{appNameUpper + "_MAIL_FROM"},
				Usage:    "Mail From address.",
				Category: "Mail options:",
			},
			&cli.StringFlag{
				Name:     "mail-to",
				EnvVars:  []string{appNameUpper + "_MAIL_TO"},
				Usage:    "Mail To address.",
				Category: "Mail options:",
			},
			&cli.StringFlag{
				Name:     "mail-server",
				EnvVars:  []string{appNameUpper + "_MAIL_SERVER"},
				Usage:    "Mail server.",
				Category: "Mail options:",
			},
			&cli.IntFlag{
				Name:     "mail-port",
				EnvVars:  []string{appNameUpper + "_MAIL_PORT"},
				Usage:    "Mail server port.",
				Category: "Mail options:",
			},
			&cli.StringFlag{
				Name:     "mail-username",
				EnvVars:  []string{appNameUpper + "_MAIL_USERNAME"},
				Usage:    "Mail server username.",
				Category: "Mail options:",
			},
			&cli.StringFlag{
				Name:     "mail-password",
				EnvVars:  []string{appNameUpper + "_MAIL_PASSWORD"},
				Usage:    "Mail server password.",
				Category: "Mail options:",
			},
			&cli.StringFlag{
				Name:     "mail-auth",
				EnvVars:  []string{appNameUpper + "_MAIL_AUTH"},
				Usage:    fmt.Sprintf("Mail server authentication mechanism, one of: %v.", availableMailAuthMechanisms),
				Category: "Mail options:",
				Action: func(cCtx *cli.Context, s string) error {
					if !slices.Contains(availableMailAuthMechanisms, s) {
						return fmt.Errorf("invalid mail server authentication mechanism, must be one of %v", availableMailAuthMechanisms)
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "mail-tls",
				EnvVars:  []string{appNameUpper + "_MAIL_TLS"},
				Usage:    fmt.Sprintf("Mail TLS policy, one of: %v.", availableMailTLSPolicies),
				Category: "Mail options:",
				Action: func(cCtx *cli.Context, s string) error {
					if !slices.Contains(availableMailTLSPolicies, s) {
						return fmt.Errorf("invalid mail TLS policy, must be one of %v", availableMailTLSPolicies)
					}
					return nil
				},
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

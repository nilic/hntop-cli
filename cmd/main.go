package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
)

func main() {

	app := &cli.App{
		Name:  appName,
		Usage: "display top Hacker News stories",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "interval",
				Aliases:     []string{"i"},
				EnvVars:     []string{appNameUpper + "_INTERVAL"},
				Usage:       fmt.Sprintf("interval since current time to show top HN stories from, eg. \"week\" (last week) or \"year\" (last year)\none of %v", printKeys(intervals)),
				DefaultText: "week",
				Action: func(cCtx *cli.Context, s string) error {
					if _, ok := intervals[s]; !ok {
						return fmt.Errorf("invalid interval value, should be one of %v", printKeys(intervals))
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "custom-interval",
				Aliases: []string{"c"},
				EnvVars: []string{appNameUpper + "_CUSTOM_INTERVAL"},
				Usage:   fmt.Sprintf("custom interval since current time to show top HN stories from, eg. \"12h\" (last 12 hours), \"100d\" (last 100 days) or \"6m\" (last 6 months)\nfollowing units are supported: %v", customIntervalSuffixes),
				Action: func(cCtx *cli.Context, s string) error {
					if len(s) == 1 {
						return fmt.Errorf("custom interval too short, needs to be in format <length><interval>, eg. 12h for 12 hours or 6m for 6 months")
					}
					length := s[:len(s)-1]
					if _, err := strconv.Atoi(length); err != nil {
						return fmt.Errorf("custom interval length error, needs to be in format <length><interval>, eg. 12h for 12 hours or 6m for 6 months")
					}
					last := s[len(s)-1:]
					if !slices.Contains(customIntervalSuffixes, last) {
						return fmt.Errorf("custom interval unit error, must end in one of %v", customIntervalSuffixes)
					}
					return nil
				},
			},
			// TODO: timerange
			// TODO: send e-mail
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

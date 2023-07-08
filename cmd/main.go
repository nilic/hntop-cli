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
				Name:    "last",
				Aliases: []string{"l"},
				EnvVars: []string{appNameUpper + "_LAST"},
				Usage:   fmt.Sprintf("interval since current time to show top HN stories from, eg. \"12h\" (last 12 hours), \"100d\" (last 100 days), \"6m\" (last 6 months)\nfollowing units are supported: %s", printUnits(intervals)),
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

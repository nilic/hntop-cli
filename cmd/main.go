package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  appName,
		Usage: "display top Hacker News stories",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "interval",
				Aliases: []string{"i"},
				EnvVars: []string{appNameUpper + "_INTERVAL"},
				Value:   "week",
				Usage:   "interval to show top HN stories from",
				Action: func(ctx *cli.Context, s string) error {
					if _, ok := intervals[s]; !ok {
						return fmt.Errorf("invalid interval value, should be one of %v", printKeys(intervals))
					}
					return nil
				},
			},
			// TODO: custom interval
			// TODO: timerange
			// TODO: send e-mail
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

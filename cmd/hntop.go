package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	intervals = []string{"day", "week", "month", "year"}
)

func doAction(cCtx *cli.Context) error {
	fmt.Println("Interval value is", cCtx.String("interval"))
	return nil
}

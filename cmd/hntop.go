package main

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

func Execute(cCtx *cli.Context) error {
	var startTime, endTime int64
	if cCtx.String("last") != "" {
		endTime = time.Now().Unix()
		interval := intervaltoSecs(cCtx.String("last"))
		startTime = endTime - interval
	} else if cCtx.String("from") != "" {
		if cCtx.String("to") != "" {
			e, _ := time.Parse(time.RFC3339, cCtx.String("to"))
			endTime = e.Unix()
		} else {
			endTime = time.Now().Unix()
		}
		s, _ := time.Parse(time.RFC3339, cCtx.String("from"))
		startTime = s.Unix()
	} else { //should print usage information
		endTime = time.Now().Unix()
		interval := intervaltoSecs("1w")
		startTime = endTime - interval
	}

	if startTime < 0 {
		startTime = 0
	}

	hnclient := NewClient()

	query := fmt.Sprintf("search?numericFilters=created_at_i>%d,created_at_i<%d", startTime, endTime)
	req, err := hnclient.NewRequest(query)
	if err != nil {
		return err
	}

	var h Hits
	err = hnclient.Do(req, &h)
	if err != nil {
		return err
	}

	h.PrintConsole(startTime, endTime)

	return nil
}

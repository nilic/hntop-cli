package main

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

const (
	itemBaseURL = "https://news.ycombinator.com/item?id="
	userBaseURL = "https://news.ycombinator.com/user?id="
)

var (
	intervals = map[string]int64{
		"day":   60 * 60 * 24,
		"week":  60 * 60 * 24 * 7,
		"month": 60 * 60 * 24 * 30,
		"year":  60 * 60 * 24 * 365,
	}
)

func Execute(cCtx *cli.Context) error {
	var h Hits
	endTime := time.Now().Unix()
	startTime := endTime - intervals[cCtx.String("interval")]

	hnclient := NewClient()
	req, err := hnclient.NewRequest(fmt.Sprintf("search?numericFilters=created_at_i>%d,created_at_i<%d", startTime, endTime))
	if err != nil {
		return err
	}
	_, err = hnclient.Do(req, &h)
	fmt.Printf("%v", h)
	return nil
}

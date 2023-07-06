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
		"hour":  60 * 60,
		"day":   60 * 60 * 24,
		"week":  60 * 60 * 24 * 7,
		"month": 60 * 60 * 24 * 30,
		"year":  60 * 60 * 24 * 365,
	}
)

func Execute(cCtx *cli.Context) error {
	endTime := time.Now().Unix()
	startTime := endTime - intervals[cCtx.String("interval")]

	hnclient := NewClient()

	query := fmt.Sprintf("search?numericFilters=created_at_i>%d,created_at_i<%d", startTime, endTime)
	req, err := hnclient.NewRequest(query)
	if err != nil {
		return err
	}

	var h Hits
	_, err = hnclient.Do(req, &h) //TODO
	if err != nil {
		return err
	}
	fmt.Printf("%v", h)
	return nil
}

func printKeys(m map[string]int64) []string {
	keys := make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

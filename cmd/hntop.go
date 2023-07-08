package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	intervals = map[string]int64{
		"hour":  60 * 60,
		"day":   60 * 60 * 24,
		"week":  60 * 60 * 24 * 7,
		"month": 60 * 60 * 24 * 30,
		"year":  60 * 60 * 24 * 365,
	}
	customIntervalSuffixes = []string{"h", "d", "w", "m", "y"}
)

func Execute(cCtx *cli.Context) error {
	endTime := time.Now().Unix()
	interval := intervaltoTime(cCtx)
	startTime := endTime - interval
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

func intervaltoTime(cCtx *cli.Context) int64 {
	if cCtx.String("interval") != "" {
		return intervals[cCtx.String("interval")]
	}

	if cCtx.String("custom-interval") != "" {
		s := cCtx.String("custom-interval")
		l, _ := strconv.Atoi(s[:len(s)-1])
		length := int64(l)
		switch unit := s[len(s)-1:]; unit {
		case "h":
			return length * intervals["hour"]
		case "d":
			return length * intervals["day"]
		case "w":
			return length * intervals["week"]
		case "m":
			return length * intervals["month"]
		case "y":
			return length * intervals["year"]
		default:
			return intervals["week"]
		}
	}

	return intervals["week"]
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

package main

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

const (
	queryPrefix        = "search?"
	defaultInterval    = "1w"
	frontPagePostCount = 30
)

type Query struct {
	ResultCount int
	StartTime   int64
	EndTime     int64
	FrontPage   bool
	Tags        string
	Query       string
}

func (q *Query) buildQuery(cCtx *cli.Context) {
	if cCtx.Bool("front-page") {
		q.FrontPage = true
		q.ResultCount = frontPagePostCount
		q.Query = queryPrefix +
			"tags=front_page" +
			fmt.Sprintf("&hitsPerPage=%d", q.ResultCount)
		return
	}

	if cCtx.String("last") != "" {
		q.EndTime = time.Now().Unix()
		interval := intervaltoSecs(cCtx.String("last"))
		q.StartTime = q.EndTime - interval
	} else if cCtx.String("from") != "" {
		if cCtx.String("to") != "" {
			e, _ := time.Parse(time.RFC3339, cCtx.String("to"))
			q.EndTime = e.Unix()
		} else {
			q.EndTime = time.Now().Unix()
		}
		s, _ := time.Parse(time.RFC3339, cCtx.String("from"))
		q.StartTime = s.Unix()
	} else {
		q.EndTime = time.Now().Unix()
		interval := intervaltoSecs(defaultInterval)
		q.StartTime = q.EndTime - interval
	}

	if q.StartTime < 0 {
		q.StartTime = 0
	}

	q.Tags = cCtx.String("tags")
	q.ResultCount = cCtx.Int("count")

	q.Query = queryPrefix +
		fmt.Sprintf("numericFilters=created_at_i>%d,created_at_i<%d", q.StartTime, q.EndTime) +
		fmt.Sprintf("&hitsPerPage=%d", q.ResultCount) +
		fmt.Sprintf("&tags=(%s)", q.Tags)
}

package hntopclient

import (
	"fmt"
	"time"
)

const (
	queryPrefix        = "search?"
	defaultInterval    = "1w"
	frontPagePostCount = 30
)

type QueryParams struct {
	FrontPage bool
	Last      string
	From      string
	To        string
	Tags      string
	Count     int
}

type Query struct {
	ResultCount int
	StartTime   int64
	EndTime     int64
	FrontPage   bool
	Tags        string
	Query       string
}

func NewQuery(qp QueryParams) *Query {
	var q Query
	if qp.FrontPage {
		q.FrontPage = true
		q.ResultCount = frontPagePostCount
		q.Query = queryPrefix +
			"tags=front_page" +
			fmt.Sprintf("&hitsPerPage=%d", q.ResultCount)
		return &q
	}

	if qp.Last != "" {
		q.EndTime = time.Now().Unix()
		interval := intervaltoSecs(qp.Last)
		q.StartTime = q.EndTime - interval
	} else if qp.From != "" {
		if qp.To != "" {
			e, _ := time.Parse(time.RFC3339, qp.To)
			q.EndTime = e.Unix()
		} else {
			q.EndTime = time.Now().Unix()
		}
		s, _ := time.Parse(time.RFC3339, qp.From)
		q.StartTime = s.Unix()
	} else {
		q.EndTime = time.Now().Unix()
		interval := intervaltoSecs(defaultInterval)
		q.StartTime = q.EndTime - interval
	}

	if q.StartTime < 0 {
		q.StartTime = 0
	}

	q.Tags = qp.Tags
	q.ResultCount = qp.Count

	q.Query = queryPrefix +
		fmt.Sprintf("numericFilters=created_at_i>%d,created_at_i<%d", q.StartTime, q.EndTime) +
		fmt.Sprintf("&hitsPerPage=%d", q.ResultCount) +
		fmt.Sprintf("&tags=(%s)", q.Tags)

	return &q
}

package main

import (
	"fmt"

	"github.com/xeonx/timeago"
)

const (
	itemBaseURL = "https://news.ycombinator.com/item?id="
	userBaseURL = "https://news.ycombinator.com/user?id="
)

func (h *Hits) Print() {
	for i, s := range h.Hits {
		fmt.Printf("%d. %s\n%s\n%d points by %s, %s | %d comments\n\n", i+1, s.Title, s.getURL(), s.Points, s.Author, timeago.English.Format(s.CreatedAt), s.NumComments)
	}
}

func (h *Hit) getURL() string {
	if h.URL != "" {
		return h.URL
	}
	return itemBaseURL + h.ObjectID
}

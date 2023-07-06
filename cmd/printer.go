package main

import (
	"fmt"

	"github.com/sersh88/timeago"
)

const (
	itemBaseURL = "https://news.ycombinator.com/item?id="
	userBaseURL = "https://news.ycombinator.com/user?id="
)

func (h *Hits) PrintConsole() {
	for i, s := range h.Hits {
		fmt.Printf("%d. %s\n%s\n%d points by %s, %s | %d comments\n\n", i+1, s.Title, s.getItemURL(), s.Points, s.Author, timeago.New(s.CreatedAt).Format(), s.NumComments)
	}
}

func (h *Hits) PrintHTML() {
	for i, s := range h.Hits {
		fmt.Printf("%d. <a href=\"%s\">%s</a>\n%d points by <a href=\"%s\">%s</a>, %s | %d comments\n\n", i+1, s.getItemURL(), s.Title, s.Points, s.getUserURL(), s.Author, timeago.New(s.CreatedAt).Format(), s.NumComments)
	}
}

func (h *Hit) getItemURL() string {
	if h.URL != "" {
		return h.URL
	}
	return itemBaseURL + h.ObjectID
}

func (h *Hit) getUserURL() string {
	return userBaseURL + h.Author
}

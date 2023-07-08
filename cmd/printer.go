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
		fmt.Printf("%d. %s\n%s\n%s\n%d points by %s, %s | %d comments\n\n", i+1, s.Title, s.getExternalURL(), s.getItemURL(), s.Points, s.Author, timeago.New(s.CreatedAt).Format(), s.NumComments)
	}
}

func (h *Hits) PrintHTML() {
	for i, s := range h.Hits {
		fmt.Printf("%d. <a href=\"%s\">%s</a>\n%d points by <a href=\"%s\">%s</a>, %s | <a href=\"%s\">%d comments</a>\n\n", i+1, s.getExternalURL(), s.Title, s.Points, s.getUserURL(), s.Author, timeago.New(s.CreatedAt).Format(), s.getItemURL(), s.NumComments)
	}
}

func (h *Hit) getItemURL() string {
	return itemBaseURL + h.ObjectID
}

func (h *Hit) getExternalURL() string {
	if h.URL != "" {
		return h.URL
	}

	return h.getItemURL()
}

func (h *Hit) getUserURL() string {
	return userBaseURL + h.Author
}

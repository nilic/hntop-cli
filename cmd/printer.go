package main

import (
	"fmt"
	"time"

	"github.com/sersh88/timeago"
)

const (
	itemBaseURL = "https://news.ycombinator.com/item?id="
	userBaseURL = "https://news.ycombinator.com/user?id="
)

func (h *Hits) PrintConsole(q *Query) {
	var heading string
	if q.FrontPage {
		heading = "Displaying HN posts currently on the front page\n"
	} else {
		heading = fmt.Sprintf("Displaying %d top HN posts from %s to %s\n", q.ResultCount, (time.Unix(q.StartTime, 0)).Format(time.RFC822), (time.Unix(q.EndTime, 0)).Format(time.RFC822))
	}
	fmt.Printf("\n" + heading + "\n")

	for i, s := range h.Hits {
		fmt.Printf("%d. %s\n", i+1, s.Title)
		fmt.Println(s.getExternalURL())
		if s.getItemURL() != s.getExternalURL() {
			fmt.Println(s.getItemURL())
		}
		fmt.Printf("%d points by %s %s | %d comments\n\n", s.Points, s.Author, timeago.New(s.CreatedAt).Format(), s.NumComments)
	}
}

func (h *Hits) PrintHTML() {
	for i, s := range h.Hits {
		fmt.Printf("%d. <a href=\"%s\">%s</a>\n", i+1, s.getExternalURL(), s.Title)
		fmt.Printf("%d points by <a href=\"%s\">%s</a> %s | <a href=\"%s\">%d comments</a>\n\n", s.Points, s.getUserURL(), s.Author, timeago.New(s.CreatedAt).Format(), s.getItemURL(), s.NumComments)
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

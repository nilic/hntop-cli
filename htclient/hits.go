package htclient

import (
	"net/url"
	"time"
)

const (
	fromBaseURL = "https://news.ycombinator.com/from?site="
	itemBaseURL = "https://news.ycombinator.com/item?id="
	userBaseURL = "https://news.ycombinator.com/user?id="
)

type Hit struct {
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Author      string    `json:"author"`
	Points      int       `json:"points"`
	NumComments int       `json:"num_comments"`
	ObjectID    string    `json:"objectID"`
}

type Hits struct {
	Hits        []Hit `json:"hits"`
	NbHits      int   `json:"nbHits"`
	Page        int   `json:"page"`
	NbPages     int   `json:"nbPages"`
	HitsPerPage int   `json:"hitsPerPage"`
}

func (h Hit) ItemURL() string {
	return itemBaseURL + h.ObjectID
}

func (h Hit) ExternalURL() string {
	if h.URL != "" {
		return h.URL
	}

	return h.ItemURL()
}

func (h Hit) UserURL() string {
	return userBaseURL + h.Author
}

func (h Hit) BaseExternalURL() string {
	if h.URL == "" {
		return ""
	}

	u, err := url.Parse(h.URL)
	if err != nil {
		return ""
	}

	return u.Hostname()
}

func (h Hit) FromURL() string {
	b := h.BaseExternalURL()

	if b == "" {
		return ""
	}

	return fromBaseURL + b
}

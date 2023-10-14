package htclient

import (
	"fmt"
	"time"

	"github.com/nilic/hntop-cli/internal/client"
)

const (
	apiBaseURL = "https://hn.algolia.com/api/v1/"
)

type HTClient struct {
	URL       string
	UserAgent string
}

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

func NewClient(query, userAgent string) *HTClient {
	return &HTClient{
		URL:       apiBaseURL + query,
		UserAgent: userAgent,
	}
}

func (c *HTClient) Do() (*Hits, error) {
	var h Hits
	h, err := client.MakeHTTPRequest("GET", c.URL, c.UserAgent, nil, nil, h)
	if err != nil {
		return nil, fmt.Errorf("calling HN API: %w", err)
	}
	return &h, err
}

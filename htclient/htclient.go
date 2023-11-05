package htclient

import (
	"fmt"

	"github.com/nilic/hntop-cli/internal/client"
)

const (
	apiBaseURL = "https://hn.algolia.com/api/v1/"
)

type HTClient struct {
	URL       string
	UserAgent string
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

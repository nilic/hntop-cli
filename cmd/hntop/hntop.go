package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nilic/hntop-cli/internal/client"
	"github.com/urfave/cli/v2"
)

const (
	apiBaseURL = "https://hn.algolia.com/api/v1/"
)

var (
	userAgent = appName + "/" + getVersion()
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

func Execute(cCtx *cli.Context) error {
	q := buildQuery(cCtx)

	fullURL := apiBaseURL + q.Query

	apiClient, err := client.NewClient(fullURL, userAgent)
	if err != nil {
		return fmt.Errorf("creating API client: %w", err)
	}

	req, err := apiClient.NewRequest(http.MethodGet, nil, nil)
	if err != nil {
		return fmt.Errorf("creating API query: %w", err)
	}

	var h Hits
	err = apiClient.Do(req, &h)
	if err != nil {
		return fmt.Errorf("calling HN API: %w", err)
	}

	err = h.Output(cCtx, q)
	if err != nil {
		return fmt.Errorf("creating output: %w", err)
	}

	return nil
}

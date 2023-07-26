package main

import (
	"fmt"
	"net/http"

	"github.com/nilic/hntop-cli/internal/client"
	"github.com/urfave/cli/v2"
)

const (
	apiBaseURL = "https://hn.algolia.com/api/v1/"
)

var (
	userAgent = appName + "/" + getVersion()
)

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

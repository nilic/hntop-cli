package main

import (
	"fmt"

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

	hnclient, err := client.NewClient(apiBaseURL, userAgent)
	if err != nil {
		return fmt.Errorf("creating API client: %w", err)
	}

	req, err := hnclient.NewRequest(q.Query)
	if err != nil {
		return fmt.Errorf("creating API query: %w", err)
	}

	var h Hits
	err = hnclient.Do(req, &h)
	if err != nil {
		return fmt.Errorf("calling HN API: %w", err)
	}

	err = h.Output(cCtx, q)
	if err != nil {
		return fmt.Errorf("creating output: %w", err)
	}

	return nil
}

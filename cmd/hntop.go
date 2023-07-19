package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Execute(cCtx *cli.Context) error {
	q := buildQuery(cCtx)

	hnclient, err := NewClient()
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

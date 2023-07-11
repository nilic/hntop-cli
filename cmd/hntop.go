package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Execute(cCtx *cli.Context) error {
	q := &Query{}
	q.buildQuery(cCtx)

	hnclient := NewClient()
	req, err := hnclient.NewRequest(q.Query)
	if err != nil {
		return fmt.Errorf("creating API query: %w", err)
	}

	var h Hits
	err = hnclient.Do(req, &h)
	if err != nil {
		return fmt.Errorf("calling HN API: %w", err)
	}

	h.PrintConsole(q)

	return nil
}

package main

import (
	"github.com/urfave/cli/v2"
)

func Execute(cCtx *cli.Context) error {
	q := &Query{}
	q.buildQuery(cCtx)

	hnclient := NewClient()
	req, err := hnclient.NewRequest(q.Query)
	if err != nil {
		return err
	}

	var h Hits
	err = hnclient.Do(req, &h)
	if err != nil {
		return err
	}

	h.PrintConsole(q)

	return nil
}

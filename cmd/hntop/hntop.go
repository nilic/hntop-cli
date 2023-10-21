package main

import (
	"fmt"

	"github.com/nilic/hntop-cli/pkg/htclient"
	"github.com/urfave/cli/v2"
)

var (
	userAgent = appName + "/" + getVersion()
)

func execute(cCtx *cli.Context) error {
	qp := htclient.QueryParams{
		FrontPage: cCtx.Bool("front-page"),
		Last:      cCtx.String("last"),
		From:      cCtx.String("from"),
		To:        cCtx.String("to"),
		Tags:      cCtx.String("tags"),
		Count:     cCtx.Int("count"),
	}

	q := htclient.NewQuery(qp)

	c := htclient.NewClient(q.Query, userAgent)
	h, err := c.Do()
	if err != nil {
		return fmt.Errorf("invoking HN API: %w", err)
	}

	err = output(cCtx, q, h)
	if err != nil {
		return fmt.Errorf("creating output: %w", err)
	}

	return nil
}

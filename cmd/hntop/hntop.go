package main

import (
	"fmt"

	"github.com/nilic/hntop-cli/internal/mailer"
	"github.com/nilic/hntop-cli/pkg/hntopclient"
	"github.com/urfave/cli/v2"
)

const (
	mailSubject = "[hntop] Top HN posts"
)

var (
	userAgent = appName + "/" + getVersion()
)

func Execute(cCtx *cli.Context) error {
	q := buildQuery(cCtx)

	c := hntopclient.NewClient(q.Query, userAgent)
	h, err := c.Do()
	if err != nil {
		return fmt.Errorf("invoking HN API: %w", err)
	}

	output, err := h.Output(cCtx.String("output"), q)
	if err != nil {
		return fmt.Errorf("creating output: %w", err)
	}

	if cCtx.String("output") == "list" {
		fmt.Print(output)
	} else {
		mc, err := mailer.NewMailConfig(cCtx.String("mail-from"),
			cCtx.String("mail-to"),
			mailSubject,
			"html",
			output,
			cCtx.String("mail-server"),
			cCtx.Int("mail-port"),
			cCtx.String("mail-username"),
			cCtx.String("mail-password"),
			cCtx.String("mail-auth"),
			cCtx.String("mail-tls"))
		if err != nil {
			return fmt.Errorf("configuring mail options: %w", err)
		}

		err = mc.SendMail()
		if err != nil {
			return fmt.Errorf("output to mail error: %w", err)
		}
	}

	return nil
}

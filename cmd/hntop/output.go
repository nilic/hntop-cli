package main

import (
	"bytes"
	"fmt"
	tt "text/template"
	"time"

	"github.com/nilic/hntop-cli/internal/mailer"
	"github.com/nilic/hntop-cli/pkg/htclient"
	"github.com/sersh88/timeago"
	"github.com/urfave/cli/v2"
)

const (
	listTemplate = `{{if .FrontPage}}Displaying HN posts currently on the front page{{else}}Displaying top {{.ResultCount}} HN posts from {{.StartTime}} to {{.EndTime}}{{end -}}
{{if .Hits}}
{{range $i, $e := .Hits}}
{{increment $i}}. {{.Title}}{{if ne .GetBaseExternalURL ""}} ({{.GetBaseExternalURL}}){{end}}
{{.GetExternalURL}}
{{- if ne .GetItemURL .GetExternalURL}}
{{.GetItemURL}}
{{- end}}
{{.Points}} points by {{.Author}} {{timeAgo .CreatedAt}} | {{.NumComments}} comments
{{end}}
{{end}}`
)

var templateFuncs = tt.FuncMap{
	"increment": increment,
	"timeAgo":   timeAgo,
}

type templateData struct {
	FrontPage   bool
	ResultCount int
	StartTime   string
	EndTime     string
	Hits        []htclient.Hit
}

func output(cCtx *cli.Context, q *htclient.Query, h *htclient.Hits) error {
	var td = templateData{
		FrontPage:   q.FrontPage,
		ResultCount: q.ResultCount,
		StartTime:   (time.Unix(q.StartTime, 0)).Format(time.RFC822),
		EndTime:     (time.Unix(q.EndTime, 0)).Format(time.RFC822),
		Hits:        h.Hits,
	}

	switch cCtx.String("output") {
	case "list":
		output, err := outputList(td)
		if err != nil {
			return fmt.Errorf("creating list output: %w", err)
		}

		fmt.Print(output)

		return nil
	case "mail":
		mc, err := mailer.New(cCtx.String("mail-from"),
			cCtx.String("mail-server"),
			cCtx.Int("mail-port"),
			cCtx.String("mail-username"),
			cCtx.String("mail-password"),
			cCtx.String("mail-auth"),
			cCtx.String("mail-tls"))
		if err != nil {
			return fmt.Errorf("configuring mail client: %w", err)
		}

		err = mc.SendTemplate(cCtx.String("mail-to"), "hntop.tmpl", templateFuncs, td)
		if err != nil {
			return fmt.Errorf("sending mail: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("unknown output type: %s", cCtx.String("output"))
	}
}

func outputList(td templateData) (string, error) {
	t, err := tt.New("list").Funcs(templateFuncs).Parse(listTemplate)
	if err != nil {
		return "", fmt.Errorf("creating template: %w", err)
	}

	buf := new(bytes.Buffer)
	err = t.ExecuteTemplate(buf, "list", td)
	if err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

func increment(i int) int {
	return i + 1
}

func timeAgo(t time.Time) string {
	return timeago.New(t).Format()
}

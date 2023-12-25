package main

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/nilic/hntop-cli/htclient"
	"github.com/nilic/hntop-cli/internal/mailer"
	"github.com/sersh88/timeago"
	"github.com/urfave/cli/v2"
)

const (
	listTemplate = `{{if .FrontPage}}Displaying HN posts currently on the front page{{else}}Displaying top {{.ResultCount}} HN posts from {{.StartTime}} to {{.EndTime}}{{end -}}
{{if .Hits}}
{{range $i, $e := .Hits}}
{{increment $i}}. {{.Title}}{{if ne .BaseExternalURL ""}} ({{.BaseExternalURL}}){{end}}
{{.ExternalURL}}
{{- if ne .ItemURL .ExternalURL}}
{{.ItemURL}}
{{- end}}
{{.Points}} points by {{.Author}} {{timeAgo .CreatedAt}} | {{.NumComments}} comments
{{end}}
{{end}}`
)

var templateFuncs = template.FuncMap{
	"increment": increment,
	"mod":       mod,
	"timeAgo":   timeAgo,
}

type templateData struct {
	FrontPage   bool
	ResultCount int
	StartTime   string
	EndTime     string
	Christmas   bool
	Hits        []htclient.Hit
}

func output(cCtx *cli.Context, q *htclient.Query, h *htclient.Hits) error {
	t := time.Now()
	var td = templateData{
		FrontPage:   q.FrontPage,
		ResultCount: q.ResultCount,
		StartTime:   (time.Unix(q.StartTime, 0)).Format(time.RFC822),
		EndTime:     (time.Unix(q.EndTime, 0)).Format(time.RFC822),
		Christmas:   t.Month().String() == "December" && t.Day() == 25,
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

		if err := mc.SendTemplate(cCtx.String("mail-to"), "hntop.tmpl", templateFuncs, td); err != nil {
			return fmt.Errorf("sending mail: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("unknown output type: %q", cCtx.String("output"))
	}
}

func outputList(td templateData) (string, error) {
	t, err := template.New("list").Funcs(templateFuncs).Parse(listTemplate)
	if err != nil {
		return "", fmt.Errorf("creating template: %w", err)
	}

	buf := new(bytes.Buffer)
	if err := t.ExecuteTemplate(buf, "list", td); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

func increment(i int) int {
	return i + 1
}

func mod(i, j int) bool {
	return i%j == 0
}

func timeAgo(t time.Time) string {
	return timeago.New(t).Format()
}

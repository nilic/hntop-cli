package main

import (
	"bytes"
	"fmt"
	ht "html/template"
	"net/url"
	tt "text/template"
	"time"

	"github.com/nilic/hntop-cli/internal/mailer"
	"github.com/sersh88/timeago"
	"github.com/urfave/cli/v2"
)

const (
	fromBaseURL = "https://news.ycombinator.com/from?site="
	itemBaseURL = "https://news.ycombinator.com/item?id="
	userBaseURL = "https://news.ycombinator.com/user?id="
	mailSubject = "[hntop] Top HN posts"

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

	htmlBodyTemplate = `{{if .FrontPage}}HN posts currently on the front page
{{else}}Top {{.ResultCount}} HN posts from {{.StartTime}} to {{.EndTime}}{{end}}<br><br>
{{if .Hits}}
	{{range $i, $e := .Hits}}
	{{increment $i}}. <a href="{{.GetExternalURL}}">{{.Title}}</a>{{if ne .GetBaseExternalURL ""}} <a href="{{.GetFromURL}}">({{.GetBaseExternalURL}})</a>{{end}}<br>
	{{.Points}} points by <a href="{{.GetUserURL}}">{{.Author}}</a> {{timeAgo .CreatedAt}} | <a href="{{.GetItemURL}}">{{.NumComments}} comments</a><br><br>
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
	Hits        []Hit
}

func (h *Hits) Output(cCtx *cli.Context, q *Query) error {
	var data = templateData{
		FrontPage:   q.FrontPage,
		ResultCount: q.ResultCount,
		StartTime:   (time.Unix(q.StartTime, 0)).Format(time.RFC822),
		EndTime:     (time.Unix(q.EndTime, 0)).Format(time.RFC822),
		Hits:        h.Hits,
	}

	switch cCtx.String("output") {
	case "list":
		list, err := outputList(data)
		if err != nil {
			return fmt.Errorf("creating list output: %w", err)
		}
		fmt.Print(list)
	case "mail":
		body, err := outputHTML(data)
		if err != nil {
			return fmt.Errorf("creating mail body: %w", err)
		}

		mc, err := mailer.NewMailConfig(cCtx.String("mail-from"),
			cCtx.String("mail-to"),
			mailSubject,
			"html",
			body,
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
	default:
		return fmt.Errorf("unknown output type: %s", cCtx.String("output"))
	}

	return nil
}

func outputList(data templateData) (string, error) {
	t, err := tt.New("list").Funcs(templateFuncs).Parse(listTemplate)
	if err != nil {
		return "", fmt.Errorf("creating template: %w", err)
	}

	buf := new(bytes.Buffer)
	err = t.ExecuteTemplate(buf, "list", data)
	if err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

func outputHTML(data templateData) (string, error) {
	t, err := ht.New("htmlBody").Funcs(templateFuncs).Parse(htmlBodyTemplate)
	if err != nil {
		return "", fmt.Errorf("creating template: %w", err)
	}

	buf := new(bytes.Buffer)
	err = t.ExecuteTemplate(buf, "htmlBody", data)
	if err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

func (h Hit) GetItemURL() string {
	return itemBaseURL + h.ObjectID
}

func (h Hit) GetExternalURL() string {
	if h.URL != "" {
		return h.URL
	}

	return h.GetItemURL()
}

func (h Hit) GetUserURL() string {
	return userBaseURL + h.Author
}

func (h Hit) GetBaseExternalURL() string {
	if h.URL == "" {
		return ""
	}

	u, err := url.Parse(h.URL)
	if err != nil {
		return ""
	}

	return u.Hostname()
}

func (h Hit) GetFromURL() string {
	b := h.GetBaseExternalURL()

	if b == "" {
		return ""
	}

	return fromBaseURL + b
}

func increment(i int) int {
	return i + 1
}

func timeAgo(t time.Time) string {
	return timeago.New(t).Format()
}

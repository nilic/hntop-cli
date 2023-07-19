package main

import (
	"bytes"
	"fmt"
	ht "html/template"
	tt "text/template"
	"time"

	"github.com/nilic/hntop-cli/pkg/mailer"
	"github.com/sersh88/timeago"
	"github.com/urfave/cli/v2"
)

const (
	itemBaseURL  = "https://news.ycombinator.com/item?id="
	userBaseURL  = "https://news.ycombinator.com/user?id="
	mailSubject  = "[hntop] Top HN posts"
	listTemplate = `{{if .FrontPage}}Displaying HN posts currently on the front page{{else}}Displaying top {{.ResultCount}} HN posts from {{.StartTime}} to {{.EndTime}}{{end -}}
{{if .Hits}}
{{range $i, $e := .Hits}}
{{increment $i}}. {{.Title}}
{{.GetExternalURL}}
{{if ne .GetItemURL .GetExternalURL}}{{.GetItemURL}}{{end -}}
{{.Points}} points by {{.Author}} {{timeAgo .CreatedAt}} | {{.NumComments}}
{{end}}
{{end}}`
	htmlBodyTemplate = `{{if .FrontPage}}HN posts currently on the front page
{{else}}Top {{.ResultCount}} HN posts from {{.StartTime}} to {{.EndTime}}{{end}}<br><br>
{{if .Hits}}
	{{range $i, $e := .Hits}}
	{{increment $i}}. <a href="{{.GetExternalURL}}">{{.Title}}</a><br>
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
	switch cCtx.String("output") {
	case "list":
		list, err := h.ToList(q)
		if err != nil {
			return fmt.Errorf("creating mail body: %w", err)
		}
		fmt.Print(list)
	case "mail":
		body, err := h.ToHTML(q)
		if err != nil {
			return fmt.Errorf("creating mail body: %w", err)
		}
		err = h.ToMail(cCtx, body)
		if err != nil {
			return fmt.Errorf("output to mail error: %w", err)
		}
	default:
		return fmt.Errorf("unknown output type: %s", cCtx.String("output"))
	}

	return nil
}

func (h *Hits) ToMail(cCtx *cli.Context, body string) error {
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

	m, err := mailer.NewMailer(mc)
	if err != nil {
		return fmt.Errorf("creating mail message: %w", err)
	}

	err = m.Send()
	if err != nil {
		return fmt.Errorf("sending mail: %w", err)
	}

	return nil
}

func (h *Hits) ToList(q *Query) (string, error) {
	var data = templateData{
		FrontPage:   q.FrontPage,
		ResultCount: q.ResultCount,
		StartTime:   (time.Unix(q.StartTime, 0)).Format(time.RFC822),
		EndTime:     (time.Unix(q.EndTime, 0)).Format(time.RFC822),
		Hits:        h.Hits,
	}
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

func (h *Hits) ToHTML(q *Query) (string, error) {
	var data = templateData{
		FrontPage:   q.FrontPage,
		ResultCount: q.ResultCount,
		StartTime:   (time.Unix(q.StartTime, 0)).Format(time.RFC822),
		EndTime:     (time.Unix(q.EndTime, 0)).Format(time.RFC822),
		Hits:        h.Hits,
	}

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

func increment(i int) int {
	return i + 1
}

func timeAgo(t time.Time) string {
	return timeago.New(t).Format()
}

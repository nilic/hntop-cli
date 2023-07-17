package main

import (
	"fmt"
	"time"

	"github.com/nilic/hntop-cli/pkg/mailer"
	"github.com/sersh88/timeago"
	"github.com/urfave/cli/v2"
)

const (
	itemBaseURL = "https://news.ycombinator.com/item?id="
	userBaseURL = "https://news.ycombinator.com/user?id="
	mailSubject = "[hntop] Top HN posts"
)

func (h *Hits) Output(cCtx *cli.Context, q *Query) error {
	switch cCtx.String("output") {
	case "list":
		h.ToList(q)
	case "mail":
		body := h.ToHTML(q)
		h.ToMail(cCtx, body)
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

func (h *Hits) ToList(q *Query) {
	var heading string
	if q.FrontPage {
		heading = "Displaying HN posts currently on the front page\n"
	} else {
		heading = fmt.Sprintf("Displaying top %d HN posts from %s to %s\n", q.ResultCount, (time.Unix(q.StartTime, 0)).Format(time.RFC822), (time.Unix(q.EndTime, 0)).Format(time.RFC822))
	}
	fmt.Printf("\n" + heading + "\n")

	for i, s := range h.Hits {
		fmt.Printf("%d. %s\n", i+1, s.Title)
		fmt.Println(s.getExternalURL())
		if s.getItemURL() != s.getExternalURL() {
			fmt.Println(s.getItemURL())
		}
		fmt.Printf("%d points by %s %s | %d comments\n\n", s.Points, s.Author, timeago.New(s.CreatedAt).Format(), s.NumComments)
	}
}

func (h *Hits) ToHTML(q *Query) string {
	var out string
	if q.FrontPage {
		out = "HN posts currently on the front page\n\n"
	} else {
		out = fmt.Sprintf("Top %d HN posts from %s to %s\n\n", q.ResultCount, (time.Unix(q.StartTime, 0)).Format(time.RFC822), (time.Unix(q.EndTime, 0)).Format(time.RFC822))
	}
	for i, s := range h.Hits {
		out += fmt.Sprintf("%d. <a href=\"%s\">%s</a>\n", i+1, s.getExternalURL(), s.Title)
		out += fmt.Sprintf("%d points by <a href=\"%s\">%s</a> %s | <a href=\"%s\">%d comments</a>\n\n", s.Points, s.getUserURL(), s.Author, timeago.New(s.CreatedAt).Format(), s.getItemURL(), s.NumComments)
	}

	return out
}

func (h *Hit) getItemURL() string {
	return itemBaseURL + h.ObjectID
}

func (h *Hit) getExternalURL() string {
	if h.URL != "" {
		return h.URL
	}

	return h.getItemURL()
}

func (h *Hit) getUserURL() string {
	return userBaseURL + h.Author
}

package main

import (
	"fmt"

	"github.com/wneessen/go-mail"
)

var (
	availableMailAuthMechanisms = []string{"login", "plain", "crammd5", "xoauth2"}
	availableMailTLSPolicies    = []string{"mandatory", "opportunistic", "notls"}
)

type Mailer struct {
	msg    *mail.Msg
	config *MailConfig
}

type MailConfig struct {
	from        string
	to          string
	subject     string
	contentType mail.ContentType
	server      string
	port        int
	username    string
	password    string
	auth        mail.SMTPAuthType
	tls         mail.TLSPolicy
}

func (mc *MailConfig) NewMailConfig(from, to, subject, contentType, server string, port int, username, password, auth, tls string) error {
	mc.from = from
	mc.to = to
	mc.subject = subject

	switch contentType {
	case "html":
		mc.contentType = mail.TypeTextHTML
	case "text":
		mc.contentType = mail.TypeTextPlain
	default:
		return fmt.Errorf("unknown mail content type: %s", contentType)
	}

	mc.server = server
	mc.port = port
	mc.username = username
	mc.password = password

	switch auth {
	case "login":
		mc.auth = mail.SMTPAuthLogin
	case "plain":
		mc.auth = mail.SMTPAuthPlain
	case "crammd5":
		mc.auth = mail.SMTPAuthCramMD5
	case "xoauth2":
		mc.auth = mail.SMTPAuthXOAUTH2
	default:
		return fmt.Errorf("unknown mail authentication mechanism: %s", auth)
	}

	switch tls {
	case "mandatory":
		mc.tls = mail.TLSMandatory
	case "opportunistic":
		mc.tls = mail.TLSOpportunistic
	case "notls":
		mc.tls = mail.NoTLS
	default:
		return fmt.Errorf("unknown mail TLS policy: %s", tls)
	}

	return nil
}

func (m *Mailer) NewMailer(h *Hits, mc *MailConfig) error {
	m.config = mc

	m.msg = mail.NewMsg()
	if err := m.msg.From(m.config.from); err != nil {
		return fmt.Errorf("failed to set mail From address: %s", err)
	}
	if err := m.msg.To(m.config.to); err != nil {
		return fmt.Errorf("failed to set mail To address: %s", err)
	}

	m.msg.Subject(m.config.subject)
	if m.config.contentType == mail.TypeTextHTML {
		m.msg.SetBodyString(m.config.contentType, "Do you like this mail? I certainly do!") // TODO HTML
	} else {
		m.msg.SetBodyString(m.config.contentType, "Do you like this mail? I certainly do!") // TODO text
	}

	return nil
}

func (m *Mailer) Send() error {
	c, err := mail.NewClient(m.config.server, mail.WithPort(m.config.port), mail.WithSMTPAuth(m.config.auth),
		mail.WithTLSPolicy(m.config.tls), mail.WithUsername(m.config.username), mail.WithPassword(m.config.password))
	if err != nil {
		return fmt.Errorf("failed to create mail client: %s", err)
	}

	if err := c.DialAndSend(m.msg); err != nil {
		return fmt.Errorf("failed to send mail: %s", err)
	}

	return nil
}

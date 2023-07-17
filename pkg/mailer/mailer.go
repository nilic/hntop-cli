package mailer

import (
	"fmt"

	"github.com/wneessen/go-mail"
)

const (
	DefaultAuthMechanism = "login"
	DefaultTLSPolicy     = "mandatory"
)

var (
	AvailableAuthMechanisms = []string{"login", "plain", "crammd5", "xoauth2"}
	AvailableTLSPolicies    = []string{"mandatory", "opportunistic", "notls"}
)

type Mailer struct {
	Msg    *mail.Msg
	Config *MailConfig
}

type MailConfig struct {
	From        string
	To          string
	Subject     string
	ContentType mail.ContentType
	Body        string
	Server      string
	Port        int
	Username    string
	Password    string
	Auth        mail.SMTPAuthType //default?
	Tls         mail.TLSPolicy    // default
}

func NewMailConfig(from, to, subject, contentType, body, server string, port int, username, password, auth, tls string) (*MailConfig, error) {
	var mc *MailConfig
	mc.From = from
	mc.To = to
	mc.Subject = subject

	switch contentType {
	case "html":
		mc.ContentType = mail.TypeTextHTML
	case "text":
		mc.ContentType = mail.TypeTextPlain
	default:
		return nil, fmt.Errorf("unknown mail content type: %s", contentType)
	}

	mc.Body = body
	mc.Server = server
	mc.Port = port
	mc.Username = username
	mc.Password = password

	switch auth {
	case "login":
		mc.Auth = mail.SMTPAuthLogin
	case "plain":
		mc.Auth = mail.SMTPAuthPlain
	case "crammd5":
		mc.Auth = mail.SMTPAuthCramMD5
	case "xoauth2":
		mc.Auth = mail.SMTPAuthXOAUTH2
	default:
		return nil, fmt.Errorf("unknown mail authentication mechanism: %s", auth)
	}

	switch tls {
	case "mandatory":
		mc.Tls = mail.TLSMandatory
	case "opportunistic":
		mc.Tls = mail.TLSOpportunistic
	case "notls":
		mc.Tls = mail.NoTLS
	default:
		return nil, fmt.Errorf("unknown mail TLS policy: %s", tls)
	}

	return mc, nil
}

func NewMailer(mc *MailConfig) (*Mailer, error) {
	var m *Mailer
	m.Config = mc

	m.Msg = mail.NewMsg()
	if err := m.Msg.From(m.Config.From); err != nil {
		return nil, fmt.Errorf("failed to set mail From address: %s", err)
	}
	if err := m.Msg.To(m.Config.To); err != nil {
		return nil, fmt.Errorf("failed to set mail To address: %s", err)
	}

	m.Msg.Subject(m.Config.Subject)
	m.Msg.SetBodyString(m.Config.ContentType, m.Config.Body)

	return m, nil
}

func (m *Mailer) Send() error {
	c, err := mail.NewClient(m.Config.Server, mail.WithPort(m.Config.Port), mail.WithSMTPAuth(m.Config.Auth),
		mail.WithTLSPolicy(m.Config.Tls), mail.WithUsername(m.Config.Username), mail.WithPassword(m.Config.Password))
	if err != nil {
		return fmt.Errorf("failed to create mail client: %s", err)
	}

	if err := c.DialAndSend(m.Msg); err != nil {
		return fmt.Errorf("failed to send mail: %s", err)
	}

	return nil
}

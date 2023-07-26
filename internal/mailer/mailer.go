package mailer

import (
	"fmt"

	"github.com/wneessen/go-mail"
)

const (
	DefaultAuthMechanism = mail.SMTPAuthLogin
	DefaultTLSPolicy     = mail.TLSMandatory
	DefaultPort          = 587
)

var (
	AvailableAuthMechanisms = []string{"login", "plain", "crammd5", "xoauth2", "none"}
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
	Auth        mail.SMTPAuthType
	Tls         mail.TLSPolicy
}

func NewMailConfig(from, to, subject, contentType, body, server string, port int, username, password, auth, tls string) (*MailConfig, error) {
	var mc MailConfig

	if from != "" {
		mc.From = from
	} else {
		return nil, fmt.Errorf("mail From address is required")
	}

	if to != "" {
		mc.To = to
	} else {
		return nil, fmt.Errorf("mail To address is required")
	}

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

	if server != "" {
		mc.Server = server
	} else {
		return nil, fmt.Errorf("mail server DNS or IP is required")
	}

	if port != 0 {
		mc.Port = port
	} else {
		mc.Port = DefaultPort
	}

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
	case "none":
		mc.Auth = ""
	default:
		mc.Auth = DefaultAuthMechanism
	}

	switch tls {
	case "mandatory":
		mc.Tls = mail.TLSMandatory
	case "opportunistic":
		mc.Tls = mail.TLSOpportunistic
	case "notls":
		mc.Tls = mail.NoTLS
	default:
		mc.Tls = DefaultTLSPolicy
	}

	return &mc, nil
}

func (mc *MailConfig) SendMail() error {
	m, err := NewMailer(mc)
	if err != nil {
		return fmt.Errorf("creating mail message: %w", err)
	}

	err = m.Send()
	if err != nil {
		return fmt.Errorf("sending mail: %w", err)
	}

	return nil
}

func NewMailer(mc *MailConfig) (*Mailer, error) {
	var m Mailer
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

	return &m, nil
}

func (m *Mailer) Send() error {
	c, err := mail.NewClient(m.Config.Server, mail.WithPort(m.Config.Port),
		mail.WithTLSPolicy(m.Config.Tls), mail.WithUsername(m.Config.Username), mail.WithPassword(m.Config.Password))
	if err != nil {
		return fmt.Errorf("failed to create mail client: %s", err)
	}

	if m.Config.Auth != "" {
		c.SetSMTPAuth(m.Config.Auth)
	}

	fmt.Printf("Sending mail to %s.. ", m.Config.To)

	if err := c.DialAndSend(m.Msg); err != nil {
		return fmt.Errorf("failed to send mail: %s", err)
	}

	fmt.Print("Done.\n")
	return nil
}

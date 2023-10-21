package mailer

import (
	"fmt"
	"time"

	"github.com/wneessen/go-mail"
)

const (
	DefaultAuthMechanism = mail.SMTPAuthLogin
	DefaultTLSPolicy     = mail.TLSMandatory
	DefaultPort          = 587
	sendRetryCount       = 3
	sendRetryWait        = 10 * time.Second
)

var (
	AvailableAuthMechanisms = []string{"login", "plain", "crammd5", "xoauth2", "none"}
	AvailableTLSPolicies    = []string{"mandatory", "opportunistic", "notls"}
)

type Mailer struct {
	From   string
	Client *mail.Client
}

func New(from, server string, port int, username, password, auth, tls string) (*Mailer, error) {
	if from == "" {
		return nil, fmt.Errorf("mail From address is required")
	}

	if server == "" {
		return nil, fmt.Errorf("mail server DNS or IP is required")
	}

	if port == 0 {
		port = DefaultPort
	}

	var authentication mail.SMTPAuthType
	switch auth {
	case "login":
		authentication = mail.SMTPAuthLogin
	case "plain":
		authentication = mail.SMTPAuthPlain
	case "crammd5":
		authentication = mail.SMTPAuthCramMD5
	case "xoauth2":
		authentication = mail.SMTPAuthXOAUTH2
	case "none":
		authentication = ""
	default:
		authentication = DefaultAuthMechanism
	}

	var tlsPolicy mail.TLSPolicy
	switch tls {
	case "mandatory":
		tlsPolicy = mail.TLSMandatory
	case "opportunistic":
		tlsPolicy = mail.TLSOpportunistic
	case "notls":
		tlsPolicy = mail.NoTLS
	default:
		tlsPolicy = DefaultTLSPolicy
	}

	client, err := mail.NewClient(server, mail.WithPort(port),
		mail.WithTLSPolicy(tlsPolicy), mail.WithUsername(username), mail.WithPassword(password))
	if err != nil {
		return nil, fmt.Errorf("failed to create mail client: %w", err)
	}

	if authentication != "" {
		client.SetSMTPAuth(authentication)
	}

	return &Mailer{
		From:   from,
		Client: client,
	}, nil
}

func (m *Mailer) SendString(to, subject, plainBody, htmlBody string) error {
	var err error

	msg := mail.NewMsg()

	if err = msg.From(m.From); err != nil {
		return fmt.Errorf("failed to set mail From address: %w", err)
	}

	if err = msg.To(to); err != nil {
		return fmt.Errorf("failed to set mail To address: %w", err)
	}

	msg.Subject(subject)
	msg.SetBodyString("text/plain", plainBody)
	msg.AddAlternativeString("text/html", htmlBody)

	fmt.Printf("Sending mail to %s.. ", to)
	m.sendMsg(msg)
	if err != nil {
		return fmt.Errorf("sending mail from string: %w", err)
	}

	fmt.Print("Done.\n")
	return nil
}

func (m *Mailer) sendMsg(msg *mail.Msg) error {
	var err error

	for i := 1; i <= sendRetryCount; i++ {
		err := m.Client.DialAndSend(msg)
		if nil == err {
			return nil
		}

		time.Sleep(sendRetryWait)
	}

	return fmt.Errorf("failed to send mail: %w", err)
}

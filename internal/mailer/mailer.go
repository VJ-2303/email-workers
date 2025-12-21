package mailer

import (
	"time"

	"github.com/go-mail/mail/v2"
)

type Mailer struct {
	dailer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) *Mailer {
	d := mail.NewDialer(host, port, username, password)
	d.Timeout = 5 * time.Second

	return &Mailer{
		dailer: d,
		sender: sender,
	}
}

func (m *Mailer) Send(recipient, replyTo, subject, body string) error {
	msg := mail.NewMessage()

	msg.SetHeader("From", m.sender)

	msg.SetHeader("To", recipient)

	if replyTo != "" {
		msg.SetAddressHeader("Reply-To", replyTo, "User via Our App")
	}

	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	err := m.dailer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}

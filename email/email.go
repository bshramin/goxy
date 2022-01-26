package email

import (
	"github.com/go-mail/mail"
)

func GetDialer(host string, port int, username, password string) *mail.Dialer {
	d := mail.NewDialer(host, port, username, password)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	return d
}

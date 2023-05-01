package mailer

import (
	"github.com/mail-server/env"
)

var Mailer *mailer

type mailer struct {
	SMTP SMTP
}

func init() {
	Mailer = &mailer{
		SMTP: SMTP{
			Server:   "smtp.gmail.com",
			Port:     587,
			Email:    env.GetEnv().SMTPEmail,
			Password: env.GetEnv().SMTPPassword,
		},
	}
}

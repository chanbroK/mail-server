package controller

import (
	"github.com/mail-server/dto"
	"github.com/mail-server/env"
	"github.com/mail-server/pkg/mailer"
	"github.com/mail-server/routes"
	"net/mail"
)

type MailController struct{}

func (*MailController) SendMail(c *routes.Context) {
	var body dto.ReqSendMail
	c.BindBody(&body)

	// default sender address
	senderAddress := mail.Address{
		Name:    "test",
		Address: env.GetEnv().SMTPEmail,
	}
	if body.Sender != nil {
		senderAddress = *body.Sender
	}

	if len(body.Receivers) == 0 {
		panic("invalid body")
	}
	err := mailer.Mailer.SMTP.Send(&mailer.Mail{
		From:    &senderAddress,
		To:      body.Receivers,
		Subject: body.Subject,
		Body:    body.Body,
		// TODO 파일 첨부 지원 X, 추후 cloud storage 주소 받아서 첨부 지원
		Attachment: nil,
	})

	if err != nil {
		panic(err)
	}

	c.JSON(map[string]string{
		"success": "true",
	})
}

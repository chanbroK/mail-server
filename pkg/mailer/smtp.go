package mailer

import (
	"crypto/tls"
	"io"
	"net/smtp"
	"strconv"
)

type SMTP struct {
	Server   string
	Port     int
	Email    string
	Password string
	client   *smtp.Client
}

// Send mail by smtp.
// 메일 보낼때마다 smtp tls 연결 & 연결 해제
func (m *SMTP) Send(msg *Mail) (err error) {
	// smtp tls 연결
	err = m.connect()

	err = m.client.Mail(msg.From.Address)
	if err != nil {
		return err
	}
	for i := range msg.To {
		err = m.client.Rcpt(msg.To[i].Address)
		if err != nil {
			return err
		}
	}

	var in io.WriteCloser
	in, err = m.client.Data()
	if err != nil {
		return err
	}
	// 메일 내용 buffer로 복사
	r, err := msg.GetReader()
	if err != nil {
		return err
	}
	_, err = io.Copy(in, r)
	if err != nil {
		return err
	}

	err = in.Close()
	if err != nil {
		return err
	}

	// smtp 연결 종료
	err = m.disconnect()

	return
}

func (m *SMTP) connect() (err error) {
	m.client, err = smtp.Dial(m.Server + ":" + strconv.Itoa(m.Port))
	if err != nil {
		return
	}

	err = m.client.StartTLS(&tls.Config{ServerName: m.Server, InsecureSkipVerify: false})
	if err != nil {
		return
	}

	auth := smtp.PlainAuth("", m.Email, m.Password, m.Server)
	err = m.client.Auth(auth)
	if err != nil {
		return
	}

	return
}

func (m *SMTP) disconnect() (err error) {
	err = m.client.Quit()
	return
}

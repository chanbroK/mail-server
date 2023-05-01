package dto

import "net/mail"

type ReqSendMail struct {
	Receivers []*mail.Address
	Subject   string
	Body      string
	Sender    *mail.Address // nullable
}

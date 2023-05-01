package router

import (
	"github.com/mail-server/controller"
	"github.com/mail-server/routes"
	"net/http"
)

var mailController = new(controller.MailController)

func NewMailRouter() http.Handler {
	router := routes.NewRouter()

	router.Post("/", mailController.SendMail)

	return router
}

package main

import (
	"github.com/mail-server/router"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/mail", router.NewMailRouter())

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

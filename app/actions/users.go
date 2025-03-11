// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package actions

import (
	"github.com/playbymail/empyr/app/domains"
	"github.com/playbymail/empyr/app/responders"
	"log"
	"net/http"
)

type CreateUserAction struct {
	Service   *domains.UserService
	Responder *responders.CreateUserResponder
}

func (a *CreateUserAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var rsp responders.CreateUserResponse

	log.Printf("%s %s: %s: htmx %q\n", r.RemoteAddr, r.Method, r.URL.Path, r.Header.Get("HX-Request"))
	username := r.FormValue("username")
	email := r.FormValue("email")

	user, err := a.Service.CreateUser(username, email)
	if err != nil {
		a.Responder.Respond(w, rsp, err)
	}
	rsp.User = user
	a.Responder.Respond(w, rsp, nil)
}

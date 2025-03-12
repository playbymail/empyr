// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package responders

import (
	"github.com/playbymail/empyr/internal/domains"
	"log"
	"net/http"
)

type LogoutUserResponder struct {
	View *ResponderTemplate
}

func NewLogoutUserResponder(view *ResponderTemplate) *LogoutUserResponder {
	return &LogoutUserResponder{View: view}
}

func (r *LogoutUserResponder) Respond(w http.ResponseWriter, isHTMX bool, err error) {
	// render the logout page
	r.View.Render(w, "logout.gohtml", nil)
}

type ShowLoginResponder struct {
	View *ResponderTemplate
}

func NewShowLoginResponder(view *ResponderTemplate) *ShowLoginResponder {
	return &ShowLoginResponder{View: view}
}

func (r *ShowLoginResponder) Respond(w http.ResponseWriter, isHTMX bool, user domains.User, err error) {
	log.Printf("responders: ShowLoginResponder: entered\n")
	// render the login page
	r.View.Render(w, "login.gohtml", struct {
		IsHTMX bool
		User   domains.User
		Error  error
	}{
		IsHTMX: isHTMX,
		User:   user,
		Error:  err,
	})
}

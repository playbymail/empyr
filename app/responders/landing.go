// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package responders

import (
	"github.com/playbymail/empyr/internal/domains"
	"net/http"
)

type ShowLandingResponder struct {
	View *ResponderTemplate
}

func NewShowLandingResponder(view *ResponderTemplate) *ShowLandingResponder {
	return &ShowLandingResponder{View: view}
}

func (r *ShowLandingResponder) Respond(w http.ResponseWriter, isHTMX bool, user domains.User, err error) {
	// render the landing page
	r.View.Render(w, "landing.gohtml", struct {
		IsHTMX bool
		User   domains.User
		Error  error
	}{
		IsHTMX: isHTMX,
		User:   user,
		Error:  err,
	})
}

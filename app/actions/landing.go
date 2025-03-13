// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package actions

import (
	"github.com/playbymail/empyr/app/responders"
	"github.com/playbymail/empyr/internal/services/sessions"
	"log"
	"net/http"
)

type ShowLandingAction struct {
	Sessions  sessions.Service
	Responder *responders.ShowLandingResponder
}

func (a *ShowLandingAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// action depends on the route (we're the catch-all for not found routes)
	// and if there is a user session.
	if r.URL.Path != "/" {
		// The / path is the garbage bin for the stdlib's router. We have to explicitly
		// handle those not found paths here. Some people would implement the static
		// file server here, but that's not the job of this action. Meh.
		log.Printf("%s %s: not found\n", r.Method, r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	// set the response parameters
	isHTMX := r.Header.Get("HX-Request") != ""
	user := a.Sessions.GetUser(r)
	// respond by showing the landing page
	a.Responder.Respond(w, isHTMX, user, nil)
}

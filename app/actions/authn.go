// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package actions

import (
	"github.com/playbymail/empyr/app/responders"
	"log"
	"net/http"
)

type ShowLoginAction struct {
	Responder *responders.ShowLoginResponder
}

func (a *ShowLoginAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var rsp responders.ShowLoginResponse
	if key := r.Header.Get("HX-Request"); key != "" {
		rsp.IsHTMX = true
	}

	log.Printf("%s %s: %s: htmx %v (%q)\n", r.RemoteAddr, r.Method, r.URL.Path, rsp.IsHTMX, r.Header.Get("HX-Request"))

	a.Responder.Respond(w, rsp, nil)
}

type LogoutAction struct {
	Responder *responders.LogoutResponder
}

func (a *LogoutAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var rsp responders.LogoutResponse
	if key := r.Header.Get("HX-Request"); key != "" {
		rsp.IsHTMX = true
	}
	a.Responder.Respond(w, rsp, nil)
}

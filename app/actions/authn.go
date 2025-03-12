// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package actions

import (
	"github.com/playbymail/empyr/app/responders"
	"github.com/playbymail/empyr/internal/services/auth"
	"github.com/playbymail/empyr/internal/services/sessions"
	"log"
	"net/http"
)

type LoginUserAction struct {
	Authentication auth.Service
	Sessions       sessions.Service
	Responder      *responders.ShowLoginResponder
}

func (a *LoginUserAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// action is to fetch the key from the route, use it to authenticate the user
	// and then create a session.
	// as a precaution, we delete any existing sessions
	_ = a.Sessions.DeleteSession(w, 0)
	// get the key from the route
	key := r.PathValue("magicKey")
	log.Printf("action: loginUser: key: %q\n", key)
	// use the key to authenticate the user
	user, err := a.Authentication.AuthenticateMagicKey(key)
	log.Printf("action: loginUser: user %+v\n", user)
	log.Printf("action: loginUser: err %v\n", err)
	if err == nil {
		// create a new session
		_, _ = a.Sessions.CreateSession(w, user.ID)
	}
	// set the response parameters
	isHTMX := r.Header.Get("HX-Request") != ""
	// respond
	a.Responder.Respond(w, isHTMX, user, err)
}

type LogoutUserAction struct {
	Sessions  sessions.Service
	Responder *responders.LogoutUserResponder
}

func (a *LogoutUserAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// action is to delete the session
	_ = a.Sessions.DeleteSession(w, 0)
	// set the response parameters
	isHTMX := r.Header.Get("HX-Request") != ""
	// respond
	a.Responder.Respond(w, isHTMX, nil)
}

type ShowLoginAction struct {
	Sessions  sessions.Service
	Responder *responders.ShowLoginResponder
}

func (a *ShowLoginAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// action is to do nothing
	// set the response parameters
	isHTMX := r.Header.Get("HX-Request") != ""
	user := a.Sessions.GetUser(r)
	// respond
	a.Responder.Respond(w, isHTMX, user, nil)
}

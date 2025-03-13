// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package actions

import (
	"fmt"
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
	// get the handle and key from the route
	handle, key := r.PathValue("handle"), r.PathValue("magicKey")
	log.Printf("action: loginUser: %q %q\n", handle, key)
	// use the key to authenticate the user
	user, err := a.Authentication.AuthenticateMagicKey(handle, key)
	if err != nil {
		// respond by showing the login page with an error
		a.Responder.Respond(w, false, user, err)
		return
	} else if !user.IsUser {
		// respond by showing the login page with an error
		a.Responder.Respond(w, false, user, fmt.Errorf("invalid credentials"))
		return
	}
	log.Printf("action: loginUser: user %+v\n", user)
	log.Printf("action: loginUser: err %v\n", err)
	if err == nil {
		// create a new session
		_, _ = a.Sessions.CreateSession(w, user.ID)
	}
	// respond by redirecting user to the login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
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
	log.Printf("action: showLogin: user %+v\n", user)
	// respond
	a.Responder.Respond(w, isHTMX, user, nil)
}

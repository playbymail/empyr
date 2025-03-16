// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package actions

import (
	"github.com/playbymail/empyr/app/responders"
	"github.com/playbymail/empyr/internal/domains"
	"github.com/playbymail/empyr/internal/services/games"
	"github.com/playbymail/empyr/internal/services/sessions"
	"log"
	"net/http"
	"strconv"
)

type ShowGamesAction struct {
	Sessions  sessions.Service
	Games     games.Service
	Responder *responders.ShowGamesResponder
}

func (a *ShowGamesAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// action depends on if there is a user session.
	user := a.Sessions.GetUser(r)
	if !user.IsUser {
		// respond by showing the login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var payload responders.ShowGamesData
	payload.IsHTMX = r.Header.Get("HX-Request") != ""
	payload.User = user

	// get the games
	var err error
	payload.Games, err = a.Games.GetListOfActiveGamesForUser(user.ID)
	if err != nil {
		log.Printf("%s %s: %v\n", r.Method, r.URL.Path, err)
	}

	// respond by showing the page
	a.Responder.Respond(w, payload, err)
}

type ShowGameAction struct {
	Sessions  sessions.Service
	Games     games.Service
	Responder *responders.ShowGameResponder
}

func (a *ShowGameAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// action depends on if there is a user session.
	user := a.Sessions.GetUser(r)
	if !user.IsUser {
		// respond by showing the login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var payload responders.ShowGameData
	payload.User = user
	payload.IsHTMX = r.Header.Get("HX-Request") != ""

	gameCode := r.PathValue("gameCode")
	if gameCode == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	empireID, err := strconv.Atoi(r.PathValue("empireID"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// get the game
	game, err := a.Games.GetEmpireGameSummary(gameCode, user.ID, int64(empireID))
	if err != nil {
		log.Printf("%s %s: %v\n", r.Method, r.URL.Path, err)
	} else {
		payload.Game = domains.Game{
			ID:          game.ID,
			Code:        game.Code,
			Name:        game.Name,
			DisplayName: game.DisplayName,
			IsActive:    game.EmpireIsActive,
			CurrentTurn: game.CurrentTurn,
		}
	}

	// respond by showing the page
	a.Responder.Respond(w, payload, err)
}

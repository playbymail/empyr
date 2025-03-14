// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package actions

import (
	"github.com/playbymail/empyr/app/responders"
	"github.com/playbymail/empyr/internal/services/games"
	"github.com/playbymail/empyr/internal/services/sessions"
	"log"
	"net/http"
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
	if games, err := a.Games.GetAllGameInfo(); err != nil {
		log.Printf("%s %s: %v\n", r.Method, r.URL.Path, err)
		a.Responder.Respond(w, payload, err)
		return
	} else {
		for _, game := range games {
			payload.Games = append(payload.Games, responders.GameInfo{
				Code:        game.Code,
				Name:        game.Name,
				DisplayName: game.DisplayName,
				IsActive:    game.IsActive,
				CurrentTurn: game.CurrentTurn,
			})
		}
	}

	// respond by showing the landing page
	a.Responder.Respond(w, payload, nil)
}

// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package responders

import (
	"github.com/playbymail/empyr/internal/domains"
	"net/http"
)

type ShowGamesResponder struct {
	View *ResponderTemplate
}

func NewShowGamesResponder(view *ResponderTemplate) *ShowGamesResponder {
	return &ShowGamesResponder{View: view}
}

type ShowGamesData struct {
	IsHTMX  bool
	User    domains.User
	IsAdmin bool
	Games   []domains.GameListing
}

func (r *ShowGamesResponder) Respond(w http.ResponseWriter, data ShowGamesData, err error) {
	// render the games page
	r.View.Render(w, "games.gohtml", data)
}

type ShowGameResponder struct {
	View *ResponderTemplate
}

func NewShowGameResponder(view *ResponderTemplate) *ShowGameResponder {
	return &ShowGameResponder{View: view}
}

type ShowGameData struct {
	IsHTMX bool
	User   domains.User
	Game   domains.Game
}

func (r *ShowGameResponder) Respond(w http.ResponseWriter, data ShowGameData, err error) {
	// render the game page
	r.View.Render(w, "game.gohtml", data)
}

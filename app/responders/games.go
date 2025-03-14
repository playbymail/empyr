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
	IsHTMX bool
	User   domains.User
	Games  []GameInfo
}
type GameInfo struct {
	Code        string
	Name        string
	DisplayName string
	IsActive    bool
	CurrentTurn int64
	EmpireCount int64
	PlayerCount int64
}

func (r *ShowGamesResponder) Respond(w http.ResponseWriter, data ShowGamesData, err error) {
	// render the games page
	r.View.Render(w, "games.gohtml", data)
}

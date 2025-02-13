// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package ec implements the logic for Empyrean Challenge
package ec

import (
	"github.com/playbymail/empyr/models/games"
	"github.com/playbymail/empyr/models/player"
)

// Engine holds the state of a single game
type Engine struct {
	// Game holds information about the current game
	Game games.Game

	// Players holds every player that has ever been in this game.
	Players map[string]player.Player

	// Orders holds every player's set of orders for the current turn.
	Orders []*Orders
}

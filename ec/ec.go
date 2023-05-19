// empyr - a game engine for Empyrean Challenge
// Copyright (C) 2023 Michael D Henderson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

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

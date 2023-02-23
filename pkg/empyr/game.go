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

package empyr

import (
	"encoding/json"
	"fmt"
	"os"
)

// Game implements the entire game state.
type Game struct {
	Id   int
	Name string
	Turn int
}

// ReadGame loads a game's data file.
func ReadGame(filename string) (Game, error) {
	var g Game
	data, err := os.ReadFile(filename)
	if err != nil {
		return g, fmt.Errorf("game: open: %w", nil)
	}
	err = json.Unmarshal(data, &g)
	if err != nil {
		return g, fmt.Errorf("game: parse: %w", nil)
	}
	return g, nil
}

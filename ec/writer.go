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

package ec

import (
	"path/filepath"
)

func (e *Engine) SaveGame(path string) error {
	path = filepath.Join(path, "out")
	game := GameJS{
		Id:   e.Game.Id,
		Name: e.Game.Name,
		Turn: e.Game.Turn,
	}
	for _, player := range e.Players {
		game.Players = append(game.Players, player.Id)
	}
	if err := tojson(path, "game", game); err != nil {
		return err
	}

	players := make(map[string]PlayerJS)
	for k, player := range e.Players {
		players[k] = PlayerJS{
			Handle: player.Handle,
			Secret: player.Secret,
			Nation: player.Nation,
		}
	}
	if err := tojson(path, "players", players); err != nil {
		return err
	}

	return nil
}

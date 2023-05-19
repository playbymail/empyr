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
	"github.com/playbymail/empyr/models/player"
)

type GameJS struct {
	Id      string   `json:"id,omitempty"`
	Name    string   `json:"name,omitempty"`
	Turn    int      `json:"turn,omitempty"`
	Players []string `json:"players,omitempty"`
}

type PlayerJS struct {
	Handle string `json:"handle,omitempty"`
	Secret string `json:"secret,omitempty"`
	Nation string `json:"nation,omitempty"`
}

func LoadGame(path string) (*Engine, error) {
	var e Engine
	e.Players = make(map[string]player.Player)

	var game GameJS
	if err := fromjson(path, "game", &game); err != nil {
		return nil, err
	}
	e.Game.Id = game.Id
	e.Game.Name = game.Name
	e.Game.Turn = game.Turn

	players := make(map[string]PlayerJS)
	if err := fromjson(path, "players", &players); err != nil {
		return nil, err
	}
	for k, p := range players {
		e.Players[k] = player.Player{
			Id:     k,
			Handle: p.Handle,
			Secret: p.Secret,
			Nation: p.Nation,
		}
	}

	return &e, nil
}

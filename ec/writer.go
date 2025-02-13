// Copyright (c) 2025 Michael D Henderson. All rights reserved.

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

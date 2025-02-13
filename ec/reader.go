// Copyright (c) 2025 Michael D Henderson. All rights reserved.

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

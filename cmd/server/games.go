// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/playbymail/empyr/cmd/server/pkg/games"
	"github.com/playbymail/empyr/cmd/server/pkg/prng"
	"github.com/playbymail/empyr/cmd/server/pkg/users"
	"time"
)

type GameMeta struct {
	ID   string
	Name string // public name of the game (not guaranteed to be unique)
	Game *games.Game
	// this keeps confusing me. the player name maps to a specific user.
	// in the game, the player name maps to a specific tribe.
	Players   map[string]*users.User
	CreatedAt time.Time
}

// addPlayer adds a new player to an existing game.
// If the name is already in user or if the user is already a player
// in the game, an error is returned.
func (meta *GameMeta) addPlayer(user *users.User, name string) error {
	if _, ok := meta.Players[name]; ok {
		return ErrDuplicatePlayer
	}
	for _, u := range meta.Players {
		if u.ID == user.ID {
			return ErrDuplicateUserName
		}
	}
	meta.Players[name] = user
	return nil
}

// createGame creates a new game and registers it with the engine.
func (s *server) createGame(id, name, seedString string) (*GameMeta, error) {
	game, err := games.DefaultGenerator()(prng.New(s.seed(seedString)))
	if err != nil {
		return nil, err
	}
	if id == "" {
		id = uuid.New().String()
	}
	if name == "" {
		name = fmt.Sprintf("GAME-%06X", len(s.games)+1)
	}
	meta := &GameMeta{
		ID:        id,
		Name:      name,
		Game:      game,
		Players:   make(map[string]*users.User),
		CreatedAt: time.Now(),
	}
	s.games[meta.ID] = meta
	return meta, nil
}

func (meta *GameMeta) MarshalJSON() ([]byte, error) {
	data := struct {
		ID        string            `json:"game_id"`
		Name      string            `json:"game_name"`
		CreatedAt time.Time         `json:"created_at"`
		Players   map[string]string `json:"players,omitempty"`
		Game      *games.Game       `json:"data,omitempty"`
	}{
		ID:        meta.ID,
		Name:      meta.Name,
		CreatedAt: meta.CreatedAt,
		Players:   make(map[string]string),
		Game:      meta.Game,
	}
	for name, user := range meta.Players {
		data.Players[name] = user.ID
	}
	return json.Marshal(&data)
}

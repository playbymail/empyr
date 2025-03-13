// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package store

import "github.com/playbymail/empyr/internal/domains"

// this file implements functions for the Games service.

func (s *Store) GetAllGames() ([]domains.Game, error) {
	// get all games
	rows, err := s.Queries.ReadAllGames(s.Context)
	if err != nil {
		return nil, err
	}
	// convert rows to games
	games := make([]domains.Game, 0, len(rows))
	for _, row := range rows {
		games = append(games, domains.Game{
			ID:          domains.GameID(row.ID),
			Code:        row.Code,
			Name:        row.Name,
			DisplayName: row.DisplayName,
			IsActive:    row.IsActive == 1,
			CurrentTurn: int(row.CurrentTurn),
		})
	}
	return games, nil
}

func (s *Store) GetUsersGames(userID domains.UserID) ([]domains.Game, error) {
	// get all games
	rows, err := s.Queries.ReadUsersGames(s.Context, int64(userID))
	if err != nil {
		return nil, err
	}
	// convert rows to games
	games := make([]domains.Game, 0, len(rows))
	for _, row := range rows {
		games = append(games, domains.Game{
			ID:          domains.GameID(row.ID),
			Code:        row.Code,
			Name:        row.Name,
			DisplayName: row.DisplayName,
			IsActive:    row.IsActive == 1,
			CurrentTurn: int(row.CurrentTurn),
		})
	}
	return games, nil
}

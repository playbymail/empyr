// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package store

import "github.com/playbymail/empyr/internal/domains"

// this file implements functions for the Games service.

func (s *Store) GetAllGameInfo() ([]domains.GameInfo, error) {
	// get all games
	rows, err := s.Queries.ReadAllGameInfo(s.Context)
	if err != nil {
		return nil, err
	}
	// convert rows to games
	games := make([]domains.GameInfo, 0, len(rows))
	for _, row := range rows {
		games = append(games, domains.GameInfo{
			ID:          domains.GameID(row.ID),
			Code:        row.Code,
			Name:        row.Name,
			DisplayName: row.DisplayName,
			IsActive:    row.IsActive == 1,
			CurrentTurn: row.CurrentTurn,
			EmpireCount: row.EmpireCount,
			PlayerCount: row.PlayerCount,
		})
	}
	return games, nil
}

func (s *Store) GetUserGames(userID domains.UserID) ([]domains.UserGame, error) {
	// get all games
	rows, err := s.Queries.ReadUsersGames(s.Context, int64(userID))
	if err != nil {
		return nil, err
	}
	// convert rows to games
	games := make([]domains.UserGame, 0, len(rows))
	for _, row := range rows {
		games = append(games, domains.UserGame{
			ID:          domains.GameID(row.ID),
			Code:        row.Code,
			Name:        row.Name,
			DisplayName: row.DisplayName,
			IsActive:    row.IsActive == 1,
			EmpireID:    domains.EmpireID(row.EmpireID.Int64),
			EmpireNo:    row.EmpireNo.Int64,
			CurrentTurn: row.CurrentTurn,
		})
	}
	return games, nil
}

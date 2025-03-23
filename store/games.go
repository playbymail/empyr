// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package store

import (
	"github.com/playbymail/empyr/internal/domains"
)

// GetAllGameInfo returns a list of all active games.
// I am not sure if this should be called by all users or just admins.
func (s *Store) GetAllGameInfo() ([]domains.GameInfo, error) {
	row, err := s.Queries.ReadAllGameInfo(s.Context)
	if err != nil {
		return nil, err
	}
	var games []domains.GameInfo
	empireCount, _ := s.Queries.ReadActiveEmpireCount(s.Context)
	games = append(games, domains.GameInfo{
		Code:        row.Code,
		Name:        row.Name,
		DisplayName: row.DisplayName,
		CurrentTurn: row.CurrentTurn,
		EmpireCount: empireCount,
		PlayerCount: empireCount,
	})
	return games, nil
}

// GetEmpireGameSummary returns a game summary for the given empire.
func (s *Store) GetEmpireGameSummary(empireID int64) (domains.UserGameSummary, error) {
	// get all games
	row, err := s.Queries.ReadEmpireByID(s.Context, empireID)
	if err != nil {
		return domains.UserGameSummary{}, err
	}
	// convert row to game
	return domains.UserGameSummary{
		Code:           row.GameCode,
		Name:           row.GameName,
		DisplayName:    row.GameDisplayName,
		EmpireID:       row.EmpireID,
		EmpireIsActive: true,
		CurrentTurn:    row.GameCurrentTurn,
	}, nil
}

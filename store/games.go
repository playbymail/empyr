// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package store

import (
	"github.com/playbymail/empyr/internal/domains"
	"github.com/playbymail/empyr/store/sqlc"
	"log"
)

// this file implements functions for the Games service.

// GetListOfActiveGamesForUser returns a list of active games that the user has an empire in.
func (s *Store) GetListOfActiveGamesForUser(userID domains.UserID) ([]domains.GameListing, error) {
	var games []domains.GameListing
	// ensure that the user is active
	user, err := s.Queries.ReadUserByID(s.Context, int64(userID))
	if err != nil {
		log.Printf("store: getListOfGames: %v\n", err)
		return nil, err
	} else if user.IsActive != 1 {
		// don't return data for inactive users
		return games, nil
	}
	rows, err := s.Queries.ReadAllGamesByUser(s.Context, int64(userID))
	if err != nil {
		log.Printf("store: getListOfGames: readAllGamesByUser: %v\n", err)
		return games, err
	}
	for _, row := range rows {
		if row.IsActive != 1 {
			// don't return data for inactive games
			continue
		}
		games = append(games, domains.GameListing{
			ID:          domains.GameID(row.ID),
			Code:        row.Code,
			Name:        row.Name,
			DisplayName: row.DisplayName,
			CurrentTurn: row.CurrentTurn,
			EmpireID:    domains.EmpireID(row.EmpireID),
			EmpireNo:    row.EmpireNo,
		})
	}
	return games, nil
}

// GetAllGameInfo returns a list of all active games.
// I am not sure if this should be called by all users or just admins.
func (s *Store) GetAllGameInfo() ([]domains.GameInfo, error) {
	// get all games
	rows, err := s.Queries.ReadAllGameInfo(s.Context)
	if err != nil {
		return nil, err
	}
	// convert rows to games
	games := make([]domains.GameInfo, 0, len(rows))
	for _, row := range rows {
		empireCount, _ := s.Queries.ReadEmpireCountByGameID(s.Context, row.ID)
		games = append(games, domains.GameInfo{
			ID:           domains.GameID(row.ID),
			Code:         row.Code,
			Name:         row.Name,
			DisplayName:  row.DisplayName,
			IsActive:     row.IsActive == 1,
			CurrentTurn:  row.CurrentTurn,
			LastEmpireNo: row.LastEmpireNo,
			EmpireCount:  empireCount,
			PlayerCount:  empireCount,
		})
	}
	return games, nil
}

// GetUserGames returns a list of games that the user is in.
func (s *Store) GetUserGames(userID domains.UserID) ([]domains.UserGame, error) {
	// get all games
	rows, err := s.Queries.ReadAllGamesByUser(s.Context, int64(userID))
	if err != nil {
		return nil, err
	}
	// convert rows to list of games
	games := make([]domains.UserGame, 0, len(rows))
	for _, row := range rows {
		games = append(games, domains.UserGame{
			ID:          domains.GameID(row.ID),
			Code:        row.Code,
			Name:        row.Name,
			DisplayName: row.DisplayName,
			IsActive:    row.IsActive == 1,
			EmpireID:    domains.EmpireID(row.EmpireID),
			EmpireNo:    row.EmpireNo,
			CurrentTurn: row.CurrentTurn,
		})
	}
	return games, nil
}

// GetEmpireGameSummary returns a game summary for the given empire.
func (s *Store) GetEmpireGameSummary(gameCode string, userID domains.UserID, empireNo int64) (domains.UserGameSummary, error) {
	// get all games
	row, err := s.Queries.ReadEmpireByGameCodeUserEmpireNo(s.Context, sqlc.ReadEmpireByGameCodeUserEmpireNoParams{
		GameCode: gameCode,
		UserID:   int64(userID),
		EmpireNo: empireNo,
	})
	if err != nil {
		return domains.UserGameSummary{}, err
	}
	// convert row to game
	return domains.UserGameSummary{
		ID:             domains.GameID(row.GameID),
		Code:           gameCode,
		Name:           row.GameName,
		DisplayName:    row.GameDisplayName,
		EmpireID:       domains.EmpireID(row.EmpireID),
		EmpireNo:       row.EmpireNo,
		EmpireIsActive: row.EmpireIsActive == 1,
		CurrentTurn:    row.GameCurrentTurn,
	}, nil
}

// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package store

import (
	"github.com/playbymail/empyr/internal/domains"
	"github.com/playbymail/empyr/store/sqlc"
	"log"
)

// this file implements functions for the Games service.

// GetListOfActiveGamesForUser returns a list of active games that the user is a player in.
func (s *Store) GetListOfActiveGamesForUser(userID domains.UserID) ([]domains.GameListing, error) {
	var games []domains.GameListing
	user, err := s.Queries.ReadUserByID(s.Context, int64(userID))
	if err != nil {
		log.Printf("store: getListOfGames: %v\n", err)
		return nil, err
	} else if user.IsActive != 1 {
		// don't return data for inactive users
		return games, nil
	}
	rows, err := s.Queries.GetListOfActiveGamesForUser(s.Context, int64(userID))
	if err != nil {
		log.Printf("store: getListOfGames: %v\n", err)
		return games, err
	}
	for _, row := range rows {
		games = append(games, domains.GameListing{
			ID:          domains.GameID(row.ID),
			Code:        row.Code,
			DisplayName: row.DisplayName,
			CurrentTurn: row.CurrentTurn,
			EmpireID:    domains.EmpireID(row.EmpireID),
			EmpireNo:    row.EmpireNo,
		})
	}
	return games, nil
}

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

func (s *Store) GetEmpireGameSummary(gameCode string, userID domains.UserID, empireNo int64) (domains.UserGameSummary, error) {
	// get all games
	row, err := s.Queries.ReadEmpireGameSummary(s.Context, sqlc.ReadEmpireGameSummaryParams{
		GameCode: gameCode,
		UserID:   int64(userID),
		EmpireNo: empireNo,
	})
	if err != nil {
		return domains.UserGameSummary{}, err
	}
	// convert row to game
	return domains.UserGameSummary{
		ID:          domains.GameID(row.ID),
		Code:        row.Code,
		Name:        row.Name,
		DisplayName: row.DisplayName,
		IsActive:    row.IsActive == 1,
		EmpireID:    domains.EmpireID(row.EmpireID.Int64),
		EmpireNo:    row.EmpireNo.Int64,
		CurrentTurn: row.CurrentTurn,
	}, nil
}

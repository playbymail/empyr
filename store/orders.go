// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package store

import (
	"github.com/playbymail/empyr/store/sqlc"
	"log"
)

// ResetTurnResults resets the results for the current turn of the game.
func (s *Store) ResetTurnResults(gameCode string) error {
	// start a transaction
	q, tx, err := s.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// get the game id and turn number
	var gameID, turnNo int64
	if row, err := q.ReadGameInfoByCode(s.Context, gameCode); err != nil {
		return err
	} else {
		gameID, turnNo = row.ID, row.CurrentTurn
	}
	log.Printf("game %q: id %d: turn: %d\n", gameCode, gameID, turnNo)
	// reset turn results so that we can re-run the turn from scratch
	// 1. delete all reports
	err = q.DeleteAllTurnReports(s.Context, sqlc.DeleteAllTurnReportsParams{GameID: gameID, TurnNo: turnNo})
	if err != nil {
		log.Printf("game %q: id %d: turn: %d: err %v\n", gameCode, gameID, turnNo, err)
		return err
	}
	log.Printf("game %q: id %d: turn: %d: purged reports\n", gameCode, gameID, turnNo)
	// 2. reset probe order status
	err = q.ResetProbeOrdersStatus(s.Context, sqlc.ResetProbeOrdersStatusParams{GameID: gameID, TurnNo: turnNo})
	if err != nil {
		log.Printf("game %q: id %d: turn: %d: probes: err %v\n", gameCode, gameID, turnNo, err)
		return err
	}
	log.Printf("game %q: id %d: turn: %d: reset probe status\n", gameCode, gameID, turnNo)
	// 3. reset survey order status
	err = q.ResetSurveyOrdersStatus(s.Context, sqlc.ResetSurveyOrdersStatusParams{GameID: gameID, TurnNo: turnNo})
	if err != nil {
		log.Printf("game %q: id %d: turn: %d: surveys: err %v\n", gameCode, gameID, turnNo, err)
		return err
	}
	log.Printf("game %q: id %d: turn: %d: reset survey status\n", gameCode, gameID, turnNo)
	// commit the transaction
	return tx.Commit()
}

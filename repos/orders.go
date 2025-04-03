// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package repos

import (
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
	// get the game turn number
	turnNo, err := q.ReadCurrentTurn(s.Context)
	if err != nil {
		return err
	}
	log.Printf("game %q: turn: %d\n", gameCode, turnNo)
	// reset turn results so that we can re-run the turn from scratch
	// 1. delete all reports
	log.Printf("game %q: turn: %d: purged reports\n", gameCode, turnNo)
	// 2. reset probe results
	err = q.DeleteSCProbeStarResultsByTurn(s.Context, turnNo)
	if err != nil {
		log.Printf("game %q: turn: %d: probes: err %v\n", gameCode, turnNo, err)
		return err
	}
	log.Printf("game %q: turn: %d: reset probe results\n", gameCode, turnNo)
	// 3. reset survey results
	err = q.DeleteSCSurveyOrbitResultsByTurn(s.Context, turnNo)
	if err != nil {
		log.Printf("game %q: turn: %d: surveys: err %v\n", gameCode, turnNo, err)
		return err
	}
	log.Printf("game %q: turn: %d: reset survey results\n", gameCode, turnNo)
	// commit the transaction
	return tx.Commit()
}

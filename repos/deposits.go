// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package repos

import (
	"github.com/playbymail/empyr/internal/domains"
	"github.com/playbymail/empyr/repos/sqlite"
	"log"
)

// CreateDeposit creates a new deposit.
func (s *Store) CreateDeposit(orbitID, depositNo int64, kind string, qty, yieldPct int64) (int64, error) {
	// start a transaction
	q, tx, err := s.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	ps := sqlite.CreateDepositParams{
		OrbitID:   orbitID,
		DepositNo: depositNo,
		Kind:      kind,
		YieldPct:  yieldPct,
	}
	depositID, err := q.CreateDeposit(s.Context, ps)
	if err != nil {
		log.Printf("error: createDeposit %v: %v\n", ps, err)
		return 0, err
	}

	hs := sqlite.CreateDepositHistoryParams{
		DepositIt: depositID,
		Effdt:     0,
		Enddt:     domains.MaxGameTurnNo,
		Qty:       qty,
	}
	err = q.CreateDepositHistory(s.Context, hs)
	if err != nil {
		log.Printf("error: createDepositHistory %v: %v\n", hs, err)
		return 0, err
	}

	// commit the transaction
	return depositID, tx.Commit()
}

// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package empires

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/playbymail/empyr/internal/cerr"
	"github.com/playbymail/empyr/internal/domains"
	"github.com/playbymail/empyr/repos"
	"github.com/playbymail/empyr/repos/sqlite"
	"log"
	"strings"
)

const (
	ErrEmpireNotAvailable = cerr.Error("empire not available")
	ErrGameInProgress     = cerr.Error("game in progress")
	ErrInvalidEmpireID    = cerr.Error("invalid empire id")
	ErrInvalidEffDt       = cerr.Error("invalid effective date")
)

type Repo struct {
	store *repos.Store
}

func NewRepo(store *repos.Store) *Repo {
	return &Repo{store: store}
}

// CreateEmpire creates a new empire in the repository. It's important
// to remind myself again that this should have no business logic. It
// should insert records. The game engine must be responsible for the
// business logic.
func (r *Repo) CreateEmpire() (int64, error) {
	db, ctx := r.store, r.store.Context
	q, tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// make sure we're not in the middle of a game
	gameRow, err := q.ReadAllGameInfo(ctx)
	if err != nil {
		return 0, err
	}
	if gameRow.CurrentTurn != 0 {
		return 0, ErrGameInProgress
	}

	// if id is not zero, then we'll use it to create the empire.
	// otherwise, we'll get the next available id from the database.
	empireID, err := q.CreateEmpire(ctx, sqlite.CreateEmpireParams{
		HomeSystemID: gameRow.HomeSystemID,
		HomeStarID:   gameRow.HomeStarID,
		HomeOrbitID:  gameRow.HomeOrbitID,
		IsActive:     1,
	})
	if err != nil {
		return 0, err
	}

	// create a default empire name. the caller must update it later.
	err = q.CreateEmpireName(ctx, sqlite.CreateEmpireNameParams{
		EmpireID: empireID,
		Effdt:    0,
		Enddt:    domains.MaxGameTurnNo,
		Name:     fmt.Sprintf("E%03d", empireID),
	})
	if err != nil {
		return 0, err
	}

	// create a default empire player. the caller must update it later.
	err = q.CreateEmpirePlayer(ctx, sqlite.CreateEmpirePlayerParams{
		EmpireID: empireID,
		Effdt:    0,
		Enddt:    domains.MaxGameTurnNo,
		Username: fmt.Sprintf("p%03d", empireID),
		Email:    fmt.Sprintf("p%03d@%s.epimethean.dev", empireID, strings.ToLower(gameRow.Code)),
	})

	// create a default system name. the caller must update it later.
	err = q.CreateEmpireSystemName(ctx, sqlite.CreateEmpireSystemNameParams{
		EmpireID: empireID,
		Effdt:    0,
		Enddt:    domains.MaxGameTurnNo,
		Name:     "Not Named",
	})

	return empireID, nil
}

// CreateEmpireWithID creates a new empire in the repository.
// It's important to remind myself again that this should have no business logic.
// It should insert records. The game engine must be responsible for the
// business logic.
func (r *Repo) CreateEmpireWithID(id int64) (int64, error) {
	if id < 1 || id > 255 {
		return 0, ErrInvalidEmpireID
	}

	db, ctx := r.store, r.store.Context
	q, tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// make sure we're not in the middle of a game
	gameRow, err := q.ReadAllGameInfo(ctx)
	if err != nil {
		return 0, err
	}
	if gameRow.CurrentTurn != 0 {
		return 0, ErrGameInProgress
	}

	empireID := id
	err = q.CreateEmpireWithID(ctx, sqlite.CreateEmpireWithIDParams{
		ID:           empireID,
		HomeSystemID: gameRow.HomeSystemID,
		HomeStarID:   gameRow.HomeStarID,
		HomeOrbitID:  gameRow.HomeOrbitID,
		IsActive:     1,
	})
	if err != nil {
		return 0, err
	}

	// create a default empire name. the caller must update it later.
	err = q.CreateEmpireName(ctx, sqlite.CreateEmpireNameParams{
		EmpireID: empireID,
		Effdt:    0,
		Enddt:    domains.MaxGameTurnNo,
		Name:     fmt.Sprintf("E%03d", empireID),
	})
	if err != nil {
		return 0, err
	}

	// create a default empire player. the caller must update it later.
	err = q.CreateEmpirePlayer(ctx, sqlite.CreateEmpirePlayerParams{
		EmpireID: empireID,
		Effdt:    0,
		Enddt:    domains.MaxGameTurnNo,
		Username: fmt.Sprintf("p%03d", empireID),
		Email:    fmt.Sprintf("p%03d@%s.epimethean.dev", empireID, strings.ToLower(gameRow.Code)),
	})

	// create a default system name. the caller must update it later.
	err = q.CreateEmpireSystemName(ctx, sqlite.CreateEmpireSystemNameParams{
		EmpireID: empireID,
		Effdt:    0,
		Enddt:    domains.MaxGameTurnNo,
		Name:     "Not Named",
	})

	return empireID, nil
}

// UpdateEmpireName updates the empire's name in the repository.
func (r *Repo) UpdateEmpireName(id, effdt int64, name string) error {
	if effdt < 0 || effdt > domains.MaxGameTurnNo {
		return ErrInvalidEffDt
	}

	db, ctx := r.store, r.store.Context
	q, tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// is there already a record for this date?
	row, err := q.ReadEmpireName(ctx, sqlite.ReadEmpireNameParams{EmpireID: id, AsOfDt: effdt})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("readEmpireName: %v", err)
			return err
		}
		// there isn't. this should never happen; the empire creation routines are
		// supposed to insert default values. it did, though, so create a new record.
		err = q.CreateEmpireName(ctx, sqlite.CreateEmpireNameParams{
			EmpireID: id,
			Effdt:    effdt,
			Enddt:    domains.MaxGameTurnNo,
			Name:     name,
		})
		if err != nil {
			return err
		}
		return tx.Commit()
	}

	// there is a record for that date.
	// if the name is the same, then we're done.
	if row.Name == name {
		return nil
	}

	// if the effdt is the same, we must be making a correction.
	if row.Effdt == effdt {
		err = q.CorrectEmpireName(ctx, sqlite.CorrectEmpireNameParams{
			EmpireID: id,
			Effdt:    row.Effdt,
			Enddt:    row.Enddt,
			Name:     name,
		})
		if err != nil {
			return err
		}
		return tx.Commit()
	}

	// otherwise, we're inserting a new record and terminating the old record.
	// first, end the old record.
	err = q.ExpireEmpireName(ctx, sqlite.ExpireEmpireNameParams{
		EmpireID: id,
		Effdt:    row.Enddt,
		Enddt:    effdt,
	})
	if err != nil {
		return err
	}
	// then create the new record
	err = q.CreateEmpireName(ctx, sqlite.CreateEmpireNameParams{
		EmpireID: id,
		Effdt:    effdt,
		Enddt:    row.Enddt,
		Name:     name,
	})
	if err != nil {
		return err
	}

	// commit the transaction
	return tx.Commit()
}

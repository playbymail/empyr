// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package sqlc

import (
	"context"
)

const createGame = `-- name: CreateGame :one

INSERT INTO games (code, name, display_name)
VALUES (?1, ?2, ?3)
RETURNING id
`

type CreateGameParams struct {
	Code        string
	Name        string
	DisplayName string
}

//	Copyright (c) 2025 Michael D Henderson. All rights reserved.
//
// CreateGame creates a new game.
func (q *Queries) CreateGame(ctx context.Context, arg CreateGameParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createGame, arg.Code, arg.Name, arg.DisplayName)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createOrbit = `-- name: CreateOrbit :one
INSERT INTO orbits (star_id, orbit_no, kind)
VALUES (?1, ?2, ?3)
RETURNING id
`

type CreateOrbitParams struct {
	StarID  int64
	OrbitNo int64
	Kind    string
}

// CreateOrbit creates a new orbit.
func (q *Queries) CreateOrbit(ctx context.Context, arg CreateOrbitParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createOrbit, arg.StarID, arg.OrbitNo, arg.Kind)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createStar = `-- name: CreateStar :one
INSERT INTO stars (system_id, sequence, scarcity)
VALUES (?1, ?2, ?3)
RETURNING id
`

type CreateStarParams struct {
	SystemID int64
	Sequence string
	Scarcity int64
}

// CreateStar creates a new star in an existing system.
func (q *Queries) CreateStar(ctx context.Context, arg CreateStarParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createStar, arg.SystemID, arg.Sequence, arg.Scarcity)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createSystem = `-- name: CreateSystem :one
INSERT INTO systems (game_id, x, y, z, scarcity)
VALUES (?1, ?2, ?3, ?4, ?5)
RETURNING id
`

type CreateSystemParams struct {
	GameID   int64
	X        int64
	Y        int64
	Z        int64
	Scarcity int64
}

// CreateSystem creates a new system.
func (q *Queries) CreateSystem(ctx context.Context, arg CreateSystemParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createSystem,
		arg.GameID,
		arg.X,
		arg.Y,
		arg.Z,
		arg.Scarcity,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createSystemDistance = `-- name: CreateSystemDistance :exec
INSERT INTO system_distances (from_system_id, to_system_id, distance)
VALUES (?1, ?2, ?3)
`

type CreateSystemDistanceParams struct {
	FromSystemID int64
	ToSystemID   int64
	Distance     int64
}

// CreateSystemDistance inserts the distance between two systems.
func (q *Queries) CreateSystemDistance(ctx context.Context, arg CreateSystemDistanceParams) error {
	_, err := q.db.ExecContext(ctx, createSystemDistance, arg.FromSystemID, arg.ToSystemID, arg.Distance)
	return err
}

const deleteEmptyOrbits = `-- name: DeleteEmptyOrbits :exec
DELETE FROM orbits
WHERE kind = 'empty'
`

// DeleteEmptyOrbits deletes all orbits with no planets.
func (q *Queries) DeleteEmptyOrbits(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteEmptyOrbits)
	return err
}

const deleteGame = `-- name: DeleteGame :exec
DELETE
FROM games
WHERE code = ?1
`

// DeleteGame deletes an existing game
func (q *Queries) DeleteGame(ctx context.Context, code string) error {
	_, err := q.db.ExecContext(ctx, deleteGame, code)
	return err
}

const getCurrentGameTurn = `-- name: GetCurrentGameTurn :one
SELECT current_turn
FROM games
WHERE id = ?1
`

// GetCurrentGameTurn gets the current game turn.
func (q *Queries) GetCurrentGameTurn(ctx context.Context, gameID int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getCurrentGameTurn, gameID)
	var current_turn int64
	err := row.Scan(&current_turn)
	return current_turn, err
}

const updateGameTurn = `-- name: UpdateGameTurn :exec
UPDATE games
SET current_turn = ?1
WHERE id = ?2
`

type UpdateGameTurnParams struct {
	TurnNumber int64
	GameID     int64
}

// UpdateGameTurn increments the game turn number.
func (q *Queries) UpdateGameTurn(ctx context.Context, arg UpdateGameTurnParams) error {
	_, err := q.db.ExecContext(ctx, updateGameTurn, arg.TurnNumber, arg.GameID)
	return err
}

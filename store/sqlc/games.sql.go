// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: games.sql

package sqlc

import (
	"context"
)

const readAllGames = `-- name: ReadAllGames :many

SELECT id,
       code,
       name,
       display_name,
       is_active,
       current_turn,
       last_empire_no,
       home_system_id,
       home_star_id,
       home_orbit_id,
       home_planet_id
FROM games
ORDER BY code
`

type ReadAllGamesRow struct {
	ID           int64
	Code         string
	Name         string
	DisplayName  string
	IsActive     int64
	CurrentTurn  int64
	LastEmpireNo int64
	HomeSystemID int64
	HomeStarID   int64
	HomeOrbitID  int64
	HomePlanetID int64
}

//	Copyright (c) 2025 Michael D Henderson. All rights reserved.
//
// ReadAllGames returns all games in the database, even the inactive ones.
func (q *Queries) ReadAllGames(ctx context.Context) ([]ReadAllGamesRow, error) {
	rows, err := q.db.QueryContext(ctx, readAllGames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadAllGamesRow
	for rows.Next() {
		var i ReadAllGamesRow
		if err := rows.Scan(
			&i.ID,
			&i.Code,
			&i.Name,
			&i.DisplayName,
			&i.IsActive,
			&i.CurrentTurn,
			&i.LastEmpireNo,
			&i.HomeSystemID,
			&i.HomeStarID,
			&i.HomeOrbitID,
			&i.HomePlanetID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const readUsersGames = `-- name: ReadUsersGames :many
SELECT games.id,
       games.code,
       games.name,
       games.display_name,
       games.is_active,
       games.current_turn,
       games.last_empire_no,
       games.home_system_id,
       games.home_star_id,
       games.home_orbit_id,
       games.home_planet_id
FROM games
         LEFT JOIN empires ON games.id = empires.game_id AND empires.user_id = ?1
WHERE is_active = 1
ORDER BY code
`

type ReadUsersGamesRow struct {
	ID           int64
	Code         string
	Name         string
	DisplayName  string
	IsActive     int64
	CurrentTurn  int64
	LastEmpireNo int64
	HomeSystemID int64
	HomeStarID   int64
	HomeOrbitID  int64
	HomePlanetID int64
}

// ReadUsersGames returns all active games that the user has a player in.
func (q *Queries) ReadUsersGames(ctx context.Context, userID int64) ([]ReadUsersGamesRow, error) {
	rows, err := q.db.QueryContext(ctx, readUsersGames, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadUsersGamesRow
	for rows.Next() {
		var i ReadUsersGamesRow
		if err := rows.Scan(
			&i.ID,
			&i.Code,
			&i.Name,
			&i.DisplayName,
			&i.IsActive,
			&i.CurrentTurn,
			&i.LastEmpireNo,
			&i.HomeSystemID,
			&i.HomeStarID,
			&i.HomeOrbitID,
			&i.HomePlanetID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

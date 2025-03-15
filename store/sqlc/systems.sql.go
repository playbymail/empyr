// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: systems.sql

package sqlc

import (
	"context"
	"database/sql"
)

const readAllSystems = `-- name: ReadAllSystems :many

SELECT systems.id      AS id,
       systems.x       as x,
       systems.y       as y,
       systems.z       as z,
       count(stars.id) AS number_of_stars
FROM games
         LEFT JOIN systems on games.id = systems.game_id
         LEFT JOIN stars on systems.id = stars.system_id
WHERE games.id = ?1
GROUP BY systems.id, systems.x, systems.y, systems.z
ORDER BY systems.id
`

type ReadAllSystemsRow struct {
	ID            sql.NullInt64
	X             sql.NullInt64
	Y             sql.NullInt64
	Z             sql.NullInt64
	NumberOfStars int64
}

//	Copyright (c) 2025 Michael D Henderson. All rights reserved.
//
// ReadAllSystems reads the system data for a game.
func (q *Queries) ReadAllSystems(ctx context.Context, gameID int64) ([]ReadAllSystemsRow, error) {
	rows, err := q.db.QueryContext(ctx, readAllSystems, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadAllSystemsRow
	for rows.Next() {
		var i ReadAllSystemsRow
		if err := rows.Scan(
			&i.ID,
			&i.X,
			&i.Y,
			&i.Z,
			&i.NumberOfStars,
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

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: clusters.sql

package sqlc

import (
	"context"
)

const readClusterMap = `-- name: ReadClusterMap :many
select id      as system_id,
       x       as x,
       y       as y,
       z       as z,
       nbr_of_stars
from systems
order by systems.id
`

type ReadClusterMapRow struct {
	SystemID   int64
	X          int64
	Y          int64
	Z          int64
	NbrOfStars int64
}

// ReadClusterMap reads the cluster map for a game
func (q *Queries) ReadClusterMap(ctx context.Context) ([]ReadClusterMapRow, error) {
	rows, err := q.db.QueryContext(ctx, readClusterMap)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadClusterMapRow
	for rows.Next() {
		var i ReadClusterMapRow
		if err := rows.Scan(
			&i.SystemID,
			&i.X,
			&i.Y,
			&i.Z,
			&i.NbrOfStars,
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

const readClusterMeta = `-- name: ReadClusterMeta :one
select home_system_id,
       home_star_id,
       home_orbit_id
from games
`

type ReadClusterMetaRow struct {
	HomeSystemID int64
	HomeStarID   int64
	HomeOrbitID  int64
}

// ReadClusterMeta reads the cluster metadata for a game
func (q *Queries) ReadClusterMeta(ctx context.Context) (ReadClusterMetaRow, error) {
	row := q.db.QueryRowContext(ctx, readClusterMeta)
	var i ReadClusterMetaRow
	err := row.Scan(&i.HomeSystemID, &i.HomeStarID, &i.HomeOrbitID)
	return i, err
}

const updateEmpireMetadata = `-- name: UpdateEmpireMetadata :exec
update games
set home_system_id = ?1,
    home_star_id   = ?2,
    home_orbit_id  = ?3
`

type UpdateEmpireMetadataParams struct {
	HomeSystemID int64
	HomeStarID   int64
	HomeOrbitID  int64
}

// UpdateEmpireMetadata updates the empire metadata in the games table.
func (q *Queries) UpdateEmpireMetadata(ctx context.Context, arg UpdateEmpireMetadataParams) error {
	_, err := q.db.ExecContext(ctx, updateEmpireMetadata, arg.HomeSystemID, arg.HomeStarID, arg.HomeOrbitID)
	return err
}

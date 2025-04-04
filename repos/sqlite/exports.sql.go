// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: exports.sql

package sqlite

import (
	"context"
)

const exportCoverTabByID = `-- name: ExportCoverTabByID :one

select games.code          as game_code,
       games.name          as game_name,
       games.display_name  as game_display_name,
       games.current_turn  as game_current_turn,
       empire.id           as empire_id,
       empire_name.name    as empire_name,
       empire_player.username,
       empire_player.email,
       empire.home_system_id,
       systems.system_name as home_system_name,
       empire.home_star_id,
       stars.star_name     as home_star_name,
       empire.home_orbit_id,
       orbits.orbit_no
from empire,
     empire_name,
     empire_player,
     systems,
     stars,
     orbits,
     games
where empire.id = ?1
  and empire.is_active = 1
  and empire_name.empire_id = empire.id
  and (empire_name.effdt <= ?2 and ?2 <= empire_name.enddt)
  and empire_player.empire_id = empire.id
  and (empire_player.effdt <= ?2 and ?2 <= empire_player.enddt)
  and systems.id = empire.home_system_id
  and stars.id = empire.home_star_id
  and orbits.id = empire.home_orbit_id
`

type ExportCoverTabByIDParams struct {
	EmpireID int64
	AsOfDt   int64
}

type ExportCoverTabByIDRow struct {
	GameCode        string
	GameName        string
	GameDisplayName string
	GameCurrentTurn int64
	EmpireID        int64
	EmpireName      string
	Username        string
	Email           string
	HomeSystemID    int64
	HomeSystemName  string
	HomeStarID      int64
	HomeStarName    string
	HomeOrbitID     int64
	OrbitNo         int64
}

//	Copyright (c) 2025 Michael D Henderson. All rights reserved.
//
// ExportCoverTabByID returns the exportable data for a single empire's cover tab.
func (q *Queries) ExportCoverTabByID(ctx context.Context, arg ExportCoverTabByIDParams) (ExportCoverTabByIDRow, error) {
	row := q.db.QueryRowContext(ctx, exportCoverTabByID, arg.EmpireID, arg.AsOfDt)
	var i ExportCoverTabByIDRow
	err := row.Scan(
		&i.GameCode,
		&i.GameName,
		&i.GameDisplayName,
		&i.GameCurrentTurn,
		&i.EmpireID,
		&i.EmpireName,
		&i.Username,
		&i.Email,
		&i.HomeSystemID,
		&i.HomeSystemName,
		&i.HomeStarID,
		&i.HomeStarName,
		&i.HomeOrbitID,
		&i.OrbitNo,
	)
	return i, err
}

const exportStarProbes = `-- name: ExportStarProbes :many
select scs.id,
       sc_probe_star_result.probe_id,
       sc_probe_star_result.effdt as as_of_dt,
       sc_probe_star_result.location,
       sc_probe_star_result.nbr_of_orbits
from scs,
     sc_probe_order,
     sc_probe_star_result
where scs.empire_id = ?1
  and sc_probe_order.sc_id = scs.id
  and sc_probe_star_result.probe_id = sc_probe_order.id
  and sc_probe_star_result.effdt = ?2
order by scs.id,
         sc_probe_star_result.location,
         sc_probe_order.id
`

type ExportStarProbesParams struct {
	EmpireID int64
	AsOfDt   int64
}

type ExportStarProbesRow struct {
	ID          int64
	ProbeID     int64
	AsOfDt      int64
	Location    string
	NbrOfOrbits int64
}

// ExportStarProbes returns the exportable data for all star probes by an empire.
func (q *Queries) ExportStarProbes(ctx context.Context, arg ExportStarProbesParams) ([]ExportStarProbesRow, error) {
	rows, err := q.db.QueryContext(ctx, exportStarProbes, arg.EmpireID, arg.AsOfDt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExportStarProbesRow
	for rows.Next() {
		var i ExportStarProbesRow
		if err := rows.Scan(
			&i.ID,
			&i.ProbeID,
			&i.AsOfDt,
			&i.Location,
			&i.NbrOfOrbits,
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

const exportSystems = `-- name: ExportSystems :many
select system_name, x, y, z, nbr_of_stars
from systems
order by system_name
`

type ExportSystemsRow struct {
	SystemName string
	X          int64
	Y          int64
	Z          int64
	NbrOfStars int64
}

// ExportSystems returns the exportable data for all systems in a game.
func (q *Queries) ExportSystems(ctx context.Context) ([]ExportSystemsRow, error) {
	rows, err := q.db.QueryContext(ctx, exportSystems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExportSystemsRow
	for rows.Next() {
		var i ExportSystemsRow
		if err := rows.Scan(
			&i.SystemName,
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

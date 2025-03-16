// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: sorcs.sql

package sqlc

import (
	"context"
)

const createSorC = `-- name: CreateSorC :one
INSERT INTO sorcs (empire_id, sorc_cd, tech_level, name, orbit_id, is_on_surface, rations, sol, birth_rate, death_rate)
VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10)
RETURNING id
`

type CreateSorCParams struct {
	EmpireID    int64
	SorcCd      string
	TechLevel   int64
	Name        string
	OrbitID     int64
	IsOnSurface int64
	Rations     float64
	Sol         float64
	BirthRate   float64
	DeathRate   float64
}

// CreateSorC creates a new ship or colony.
func (q *Queries) CreateSorC(ctx context.Context, arg CreateSorCParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createSorC,
		arg.EmpireID,
		arg.SorcCd,
		arg.TechLevel,
		arg.Name,
		arg.OrbitID,
		arg.IsOnSurface,
		arg.Rations,
		arg.Sol,
		arg.BirthRate,
		arg.DeathRate,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createSorCInventory = `-- name: CreateSorCInventory :exec
INSERT INTO inventory (sorc_id, unit_cd, qty)
VALUES (?1, ?2, ?3)
`

type CreateSorCInventoryParams struct {
	SorcID int64
	UnitCd string
	Qty    int64
}

// CreateSorCInventory creates a new ship or colony inventory entry.
func (q *Queries) CreateSorCInventory(ctx context.Context, arg CreateSorCInventoryParams) error {
	_, err := q.db.ExecContext(ctx, createSorCInventory, arg.SorcID, arg.UnitCd, arg.Qty)
	return err
}

const createSorCPopulation = `-- name: CreateSorCPopulation :exec
INSERT INTO population (sorc_id, population_cd, qty, pay_rate, rebel_qty)
VALUES (?1, ?2, ?3, ?4, ?5)
`

type CreateSorCPopulationParams struct {
	SorcID       int64
	PopulationCd string
	Qty          int64
	PayRate      float64
	RebelQty     int64
}

// CreateSorCPopulation creates a new ship or colony population entry.
func (q *Queries) CreateSorCPopulation(ctx context.Context, arg CreateSorCPopulationParams) error {
	_, err := q.db.ExecContext(ctx, createSorCPopulation,
		arg.SorcID,
		arg.PopulationCd,
		arg.Qty,
		arg.PayRate,
		arg.RebelQty,
	)
	return err
}

const readAllColoniesByEmpire = `-- name: ReadAllColoniesByEmpire :many
SELECT sorcs.id        AS sorcs_id,
       systems.x,
       systems.y,
       systems.z,
       stars.sequence  as suffix,
       orbits.orbit_no,
       sorc_codes.name AS sorc_kind,
       sorcs.tech_level,
       sorcs.name,
       sorcs.rations,
       sorcs.birth_rate,
       sorcs.death_rate,
       sorcs.sol
FROM sorcs,
     sorc_codes,
     orbits,
     stars,
     systems
WHERE sorcs.empire_id = ?1
  AND sorcs.sorc_cd IN ('COPN', 'CENC', 'CORB')
  AND sorcs.sorc_cd = sorc_codes.code
  AND orbits.star_id = stars.id
  AND systems.id = stars.system_id
ORDER BY sorcs.id
`

type ReadAllColoniesByEmpireRow struct {
	SorcsID   int64
	X         int64
	Y         int64
	Z         int64
	Suffix    string
	OrbitNo   int64
	SorcKind  string
	TechLevel int64
	Name      string
	Rations   float64
	BirthRate float64
	DeathRate float64
	Sol       float64
}

// ReadAllColoniesByEmpire reads the colonies for a given empire in a game.
func (q *Queries) ReadAllColoniesByEmpire(ctx context.Context, empireID int64) ([]ReadAllColoniesByEmpireRow, error) {
	rows, err := q.db.QueryContext(ctx, readAllColoniesByEmpire, empireID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadAllColoniesByEmpireRow
	for rows.Next() {
		var i ReadAllColoniesByEmpireRow
		if err := rows.Scan(
			&i.SorcsID,
			&i.X,
			&i.Y,
			&i.Z,
			&i.Suffix,
			&i.OrbitNo,
			&i.SorcKind,
			&i.TechLevel,
			&i.Name,
			&i.Rations,
			&i.BirthRate,
			&i.DeathRate,
			&i.Sol,
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

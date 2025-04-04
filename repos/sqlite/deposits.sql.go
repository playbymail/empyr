// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: deposits.sql

package sqlite

import (
	"context"
)

const createDeposit = `-- name: CreateDeposit :one
insert into deposits (orbit_id, deposit_no, kind, yield_pct)
values (?1, ?2, ?3, ?4)
returning id
`

type CreateDepositParams struct {
	OrbitID   int64
	DepositNo int64
	Kind      string
	YieldPct  int64
}

// CreateDeposit creates a new deposit on an existing orbit.
func (q *Queries) CreateDeposit(ctx context.Context, arg CreateDepositParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createDeposit,
		arg.OrbitID,
		arg.DepositNo,
		arg.Kind,
		arg.YieldPct,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createDepositHistory = `-- name: CreateDepositHistory :exec
insert into deposit_history (deposit_id, effdt, enddt, qty)
values (?1, ?2, ?3, ?4)
`

type CreateDepositHistoryParams struct {
	DepositIt int64
	Effdt     int64
	Enddt     int64
	Qty       int64
}

// CreateDepositHistory creates a new deposit history record.
func (q *Queries) CreateDepositHistory(ctx context.Context, arg CreateDepositHistoryParams) error {
	_, err := q.db.ExecContext(ctx, createDepositHistory,
		arg.DepositIt,
		arg.Effdt,
		arg.Enddt,
		arg.Qty,
	)
	return err
}

const createDepositSummary = `-- name: CreateDepositSummary :exec
insert into deposits_summary (orbit_id, effdt, enddt,
                              fuel_qty, fuel_est_qty,
                              gold_qty, gold_est_qty,
                              mets_qty, mets_est_qty,
                              nmts_qty, nmts_est_qty)
VALUES (?1, ?2, ?3,
        ?4, ?5,
        ?6, ?7,
        ?8, ?9,
        ?10, ?11)
`

type CreateDepositSummaryParams struct {
	OrbitID    int64
	Effdt      int64
	Enddt      int64
	FuelQty    int64
	FuelEstQty int64
	GoldQty    int64
	GoldEstQty int64
	MetsQty    int64
	MetsEstQty int64
	NmtsQty    int64
	NmtsEstQty int64
}

// CreateDepositSummary creates a new deposit summary record for a planet.
func (q *Queries) CreateDepositSummary(ctx context.Context, arg CreateDepositSummaryParams) error {
	_, err := q.db.ExecContext(ctx, createDepositSummary,
		arg.OrbitID,
		arg.Effdt,
		arg.Enddt,
		arg.FuelQty,
		arg.FuelEstQty,
		arg.GoldQty,
		arg.GoldEstQty,
		arg.MetsQty,
		arg.MetsEstQty,
		arg.NmtsQty,
		arg.NmtsEstQty,
	)
	return err
}

const createDepositsSummaryPivot = `-- name: CreateDepositsSummaryPivot :exec
insert into deposits_summary_pivot (deposit_id, effdt, enddt,
                                    fuel_qty, fuel_est_qty,
                                    gold_qty, gold_est_qty,
                                    mets_qty, mets_est_qty,
                                    nmts_qty, nmts_est_qty)
VALUES (?1, ?2, ?3,
        ?4, ?5,
        ?6, ?7,
        ?8, ?9,
        ?10, ?11)
`

type CreateDepositsSummaryPivotParams struct {
	DepositID  int64
	Effdt      int64
	Enddt      int64
	FuelQty    int64
	FuelEstQty int64
	GoldQty    int64
	GoldEstQty int64
	MetsQty    int64
	MetsEstQty int64
	NmtsQty    int64
	NmtsEstQty int64
}

// CreateDepositsSummaryPivot creates a new deposit summary pivot record.
func (q *Queries) CreateDepositsSummaryPivot(ctx context.Context, arg CreateDepositsSummaryPivotParams) error {
	_, err := q.db.ExecContext(ctx, createDepositsSummaryPivot,
		arg.DepositID,
		arg.Effdt,
		arg.Enddt,
		arg.FuelQty,
		arg.FuelEstQty,
		arg.GoldQty,
		arg.GoldEstQty,
		arg.MetsQty,
		arg.MetsEstQty,
		arg.NmtsQty,
		arg.NmtsEstQty,
	)
	return err
}

const readDepositByOrbitDepositNo = `-- name: ReadDepositByOrbitDepositNo :one
select deposits.id as deposit_id,
       deposit_no,
       kind,
       qty,
       yield_pct
from deposits,
     deposit_history
where deposits.orbit_id = ?1
  and deposits.deposit_no = ?2
  and deposit_history.deposit_id = deposits.id
  and deposit_history.effdt <= ?3
  and ?3 < deposit_history.enddt
`

type ReadDepositByOrbitDepositNoParams struct {
	OrbitID   int64
	DepositNo int64
	TurnNo    int64
}

type ReadDepositByOrbitDepositNoRow struct {
	DepositID int64
	DepositNo int64
	Kind      string
	Qty       int64
	YieldPct  int64
}

// ReadDepositByOrbitDepositNo reads a deposit by its deposit number and orbit ID.
func (q *Queries) ReadDepositByOrbitDepositNo(ctx context.Context, arg ReadDepositByOrbitDepositNoParams) (ReadDepositByOrbitDepositNoRow, error) {
	row := q.db.QueryRowContext(ctx, readDepositByOrbitDepositNo, arg.OrbitID, arg.DepositNo, arg.TurnNo)
	var i ReadDepositByOrbitDepositNoRow
	err := row.Scan(
		&i.DepositID,
		&i.DepositNo,
		&i.Kind,
		&i.Qty,
		&i.YieldPct,
	)
	return i, err
}

const readDepositSummaryByOrbitId = `-- name: ReadDepositSummaryByOrbitId :many
select deposits.orbit_id,
       deposits.id                                                          as deposit_id,
       case when deposits.kind = 'FUEL' then deposit_history.qty else 0 end as fuel_qty,
       case when deposits.kind = 'GOLD' then deposit_history.qty else 0 end as gold_qty,
       case when deposits.kind = 'METS' then deposit_history.qty else 0 end as mets_qty,
       case when deposits.kind = 'NMTS' then deposit_history.qty else 0 end as nmts_qty
from deposits,
     deposit_history
where orbit_id = ?1
  and deposit_history.deposit_id = deposits.id
  and deposit_history.effdt <= ?2
  and ?2 < deposit_history.enddt
`

type ReadDepositSummaryByOrbitIdParams struct {
	OrbitID int64
	TurnNo  int64
}

type ReadDepositSummaryByOrbitIdRow struct {
	OrbitID   int64
	DepositID int64
	FuelQty   int64
	GoldQty   int64
	MetsQty   int64
	NmtsQty   int64
}

// ReadDepositSummaryByOrbitId reads a summary of deposits on a planet.
func (q *Queries) ReadDepositSummaryByOrbitId(ctx context.Context, arg ReadDepositSummaryByOrbitIdParams) ([]ReadDepositSummaryByOrbitIdRow, error) {
	rows, err := q.db.QueryContext(ctx, readDepositSummaryByOrbitId, arg.OrbitID, arg.TurnNo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadDepositSummaryByOrbitIdRow
	for rows.Next() {
		var i ReadDepositSummaryByOrbitIdRow
		if err := rows.Scan(
			&i.OrbitID,
			&i.DepositID,
			&i.FuelQty,
			&i.GoldQty,
			&i.MetsQty,
			&i.NmtsQty,
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

const readDepositsByOrbit = `-- name: ReadDepositsByOrbit :many
select id, deposit_no, kind, qty, yield_pct
from deposits,
     deposit_history
where deposits.orbit_id = ?1
  and deposit_history.deposit_id = deposits.id
  and deposit_history.effdt <= ?2
  and ?2 < deposit_history.enddt
order by deposit_no
`

type ReadDepositsByOrbitParams struct {
	OrbitID int64
	TurnNo  int64
}

type ReadDepositsByOrbitRow struct {
	ID        int64
	DepositNo int64
	Kind      string
	Qty       int64
	YieldPct  int64
}

// ReadDepositsByOrbit returns a list of all deposits on a orbit.
func (q *Queries) ReadDepositsByOrbit(ctx context.Context, arg ReadDepositsByOrbitParams) ([]ReadDepositsByOrbitRow, error) {
	rows, err := q.db.QueryContext(ctx, readDepositsByOrbit, arg.OrbitID, arg.TurnNo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadDepositsByOrbitRow
	for rows.Next() {
		var i ReadDepositsByOrbitRow
		if err := rows.Scan(
			&i.ID,
			&i.DepositNo,
			&i.Kind,
			&i.Qty,
			&i.YieldPct,
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

const updateDepositsSummaryEndDt = `-- name: UpdateDepositsSummaryEndDt :exec
update deposits_summary
set enddt = ?1
where orbit_id = ?2
  and effdt = ?3
`

type UpdateDepositsSummaryEndDtParams struct {
	Enddt   int64
	OrbitID int64
	Effdt   int64
}

// UpdateDepositsSummaryEndDt updates the end turn for a deposit summary record.
func (q *Queries) UpdateDepositsSummaryEndDt(ctx context.Context, arg UpdateDepositsSummaryEndDtParams) error {
	_, err := q.db.ExecContext(ctx, updateDepositsSummaryEndDt, arg.Enddt, arg.OrbitID, arg.Effdt)
	return err
}

const updateDepositsSummaryPivotEndDt = `-- name: UpdateDepositsSummaryPivotEndDt :exec
update deposits_summary_pivot
set enddt = ?1
where deposit_id = ?2
  and effdt = ?3
`

type UpdateDepositsSummaryPivotEndDtParams struct {
	Enddt     int64
	DepositID int64
	Effdt     int64
}

// UpdateDepositsSummaryPivotEndDt updates the end turn for a deposit summary record.
func (q *Queries) UpdateDepositsSummaryPivotEndDt(ctx context.Context, arg UpdateDepositsSummaryPivotEndDtParams) error {
	_, err := q.db.ExecContext(ctx, updateDepositsSummaryPivotEndDt, arg.Enddt, arg.DepositID, arg.Effdt)
	return err
}

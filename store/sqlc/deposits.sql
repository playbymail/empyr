-- CreateDeposit creates a new deposit on an existing orbit.
--
-- name: CreateDeposit :one
insert into deposits (system_id, star_id, orbit_id, deposit_no, kind, qty, yield_pct)
values (:system_id, :star_id, :orbit_id, :deposit_no, :kind, :qty, :yield_pct)
returning id;

-- CreateDepositSummary creates a new deposit summary record for a planet.
--
-- name: CreateDepositSummary :exec
insert into deposits_summary (orbit_id, eff_turn, end_turn,
                              fuel_qty, fuel_est_qty,
                              gold_qty, gold_est_qty,
                              mets_qty, mets_est_qty,
                              nmts_qty, nmts_est_qty)
VALUES (:orbit_id, :eff_turn, :end_turn,
        :fuel_qty, :fuel_est_qty,
        :gold_qty, :gold_est_qty,
        :mets_qty, :mets_est_qty,
        :nmts_qty, :nmts_est_qty);

-- CreateDepositsSummaryPivot creates a new deposit summary pivot record.
--
-- name: CreateDepositsSummaryPivot :exec
insert into deposits_summary_pivot (deposit_id, eff_turn, end_turn,
                                   fuel_qty, fuel_est_qty,
                                   gold_qty, gold_est_qty,
                                   mets_qty, mets_est_qty,
                                   nmts_qty, nmts_est_qty)
VALUES (:deposit_id, :eff_turn, :end_turn,
        :fuel_qty, :fuel_est_qty,
        :gold_qty, :gold_est_qty,
        :mets_qty, :mets_est_qty,
        :nmts_qty, :nmts_est_qty);

-- ReadDepositByOrbitDepositNo reads a deposit by its deposit number and orbit ID.
--
-- name: ReadDepositByOrbitDepositNo :one
select id, deposit_no, kind, qty, yield_pct
from deposits
where orbit_id = :orbit_id
  and deposit_no = :deposit_no;

-- ReadDepositSummaryByOrbitId reads a summary of deposits on a planet.
--
-- name: ReadDepositSummaryByOrbitId :many
select deposits.orbit_id,
       deposits.id                                                   as deposit_id,
       case when deposits.kind = 'FUEL' then deposits.qty else 0 end as fuel_qty,
       case when deposits.kind = 'GOLD' then deposits.qty else 0 end as gold_qty,
       case when deposits.kind = 'METS' then deposits.qty else 0 end as mets_qty,
       case when deposits.kind = 'NMTS' then deposits.qty else 0 end as nmts_qty
from deposits
where orbit_id = :orbit_id;

-- ReadDepositsByOrbit returns a list of all deposits on a orbit.
--
-- name: ReadDepositsByOrbit :many
select id, deposit_no, kind, qty, yield_pct
from deposits
where orbit_id = :orbit_id
order by deposit_no;

-- UpdateDepositsSummaryEndDt updates the end turn for a deposit summary record.
--
-- name: UpdateDepositsSummaryEndDt :exec
update deposits_summary
set end_turn = :end_turn
where orbit_id = :orbit_id
  and eff_turn = :eff_turn;

-- UpdateDepositsSummaryPivotEndDt updates the end turn for a deposit summary record.
--
-- name: UpdateDepositsSummaryPivotEndDt :exec
update deposits_summary_pivot
set end_turn = :end_turn
where deposit_id = :deposit_id
  and eff_turn = :eff_turn;


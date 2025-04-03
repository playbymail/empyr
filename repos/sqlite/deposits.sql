-- CreateDeposit creates a new deposit on an existing orbit.
--
-- name: CreateDeposit :one
insert into deposits (orbit_id, deposit_no, kind, yield_pct)
values (:orbit_id, :deposit_no, :kind, :yield_pct)
returning id;

-- CreateDepositHistory creates a new deposit history record.
--
-- name: CreateDepositHistory :exec
insert into deposit_history (deposit_id, effdt, enddt, qty)
values (:deposit_it, :effdt, :enddt, :qty);

-- CreateDepositSummary creates a new deposit summary record for a planet.
--
-- name: CreateDepositSummary :exec
insert into deposits_summary (orbit_id, effdt, enddt,
                              fuel_qty, fuel_est_qty,
                              gold_qty, gold_est_qty,
                              mets_qty, mets_est_qty,
                              nmts_qty, nmts_est_qty)
VALUES (:orbit_id, :effdt, :enddt,
        :fuel_qty, :fuel_est_qty,
        :gold_qty, :gold_est_qty,
        :mets_qty, :mets_est_qty,
        :nmts_qty, :nmts_est_qty);

-- CreateDepositsSummaryPivot creates a new deposit summary pivot record.
--
-- name: CreateDepositsSummaryPivot :exec
insert into deposits_summary_pivot (deposit_id, effdt, enddt,
                                    fuel_qty, fuel_est_qty,
                                    gold_qty, gold_est_qty,
                                    mets_qty, mets_est_qty,
                                    nmts_qty, nmts_est_qty)
VALUES (:deposit_id, :effdt, :enddt,
        :fuel_qty, :fuel_est_qty,
        :gold_qty, :gold_est_qty,
        :mets_qty, :mets_est_qty,
        :nmts_qty, :nmts_est_qty);

-- ReadDepositByOrbitDepositNo reads a deposit by its deposit number and orbit ID.
--
-- name: ReadDepositByOrbitDepositNo :one
select deposits.id as deposit_id,
       deposit_no,
       kind,
       qty,
       yield_pct
from deposits,
     deposit_history
where deposits.orbit_id = :orbit_id
  and deposits.deposit_no = :deposit_no
  and deposit_history.deposit_id = deposits.id
  and deposit_history.effdt <= :turn_no
  and :turn_no < deposit_history.enddt;

-- ReadDepositSummaryByOrbitId reads a summary of deposits on a planet.
--
-- name: ReadDepositSummaryByOrbitId :many
select deposits.orbit_id,
       deposits.id                                                          as deposit_id,
       case when deposits.kind = 'FUEL' then deposit_history.qty else 0 end as fuel_qty,
       case when deposits.kind = 'GOLD' then deposit_history.qty else 0 end as gold_qty,
       case when deposits.kind = 'METS' then deposit_history.qty else 0 end as mets_qty,
       case when deposits.kind = 'NMTS' then deposit_history.qty else 0 end as nmts_qty
from deposits,
     deposit_history
where orbit_id = :orbit_id
  and deposit_history.deposit_id = deposits.id
  and deposit_history.effdt <= :turn_no
  and :turn_no < deposit_history.enddt;

-- ReadDepositsByOrbit returns a list of all deposits on a orbit.
--
-- name: ReadDepositsByOrbit :many
select id, deposit_no, kind, qty, yield_pct
from deposits,
     deposit_history
where deposits.orbit_id = :orbit_id
  and deposit_history.deposit_id = deposits.id
  and deposit_history.effdt <= :turn_no
  and :turn_no < deposit_history.enddt
order by deposit_no;

-- UpdateDepositsSummaryEndDt updates the end turn for a deposit summary record.
--
-- name: UpdateDepositsSummaryEndDt :exec
update deposits_summary
set enddt = :enddt
where orbit_id = :orbit_id
  and effdt = :effdt;

-- UpdateDepositsSummaryPivotEndDt updates the end turn for a deposit summary record.
--
-- name: UpdateDepositsSummaryPivotEndDt :exec
update deposits_summary_pivot
set enddt = :enddt
where deposit_id = :deposit_id
  and effdt = :effdt;


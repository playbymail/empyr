-- CreateDeposit creates a new deposit on an existing planet.
--
-- name: CreateDeposit :one
INSERT INTO deposits (planet_id, deposit_no, kind, remaining_qty, yield_pct)
VALUES (:planet_id, :deposit_no, :kind, :remaining_qty, :yield_pct)
RETURNING id;

-- DeleteEmptyDeposits deletes all empty deposits.
--
-- name: DeleteEmptyDeposits :exec
DELETE
FROM deposits
WHERE kind = 'NONE';



-- CreateDeposit creates a new deposit on an existing orbit.
--
-- name: CreateDeposit :one
INSERT INTO deposits (system_id, star_id, orbit_id, deposit_no, kind, qty, yield_pct)
VALUES (:system_id, :star_id, :orbit_id, :deposit_no, :kind, :qty, :yield_pct)
RETURNING id;

-- ReadDepositByOrbitDepositNo reads a deposit by its deposit number and orbit ID.
--
-- name: ReadDepositByOrbitDepositNo :one
SELECT id, deposit_no, kind, qty, yield_pct
FROM deposits
WHERE orbit_id = :orbit_id
  AND deposit_no = :deposit_no;

-- ReadDepositsByOrbit returns a list of all deposits on a orbit.
--
-- name: ReadDepositsByOrbit :many
SELECT id, deposit_no, kind, qty, yield_pct
FROM deposits
WHERE orbit_id = :orbit_id
ORDER BY deposit_no;
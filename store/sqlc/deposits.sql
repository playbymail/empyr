-- CreateDeposit creates a new deposit on an existing planet.
--
-- name: CreateDeposit :one
INSERT INTO deposits (planet_id, deposit_no, kind, remaining_qty, yield_pct)
VALUES (:planet_id, :deposit_no, :kind, :remaining_qty, :yield_pct)
RETURNING id;

-- ReadDepositByNo reads a deposit by its deposit number and planet ID.
--
-- name: ReadDepositByNo :one
SELECT deposits.id, deposits.deposit_no, deposits.kind, deposits.remaining_qty, deposits.yield_pct
FROM deposits
WHERE deposits.planet_id = :planet_id
  AND deposits.deposit_no = :deposit_no;

-- ReadDepositsByPlanet returns a list of all deposits on a planet.
--
-- name: ReadDepositsByPlanet :many
SELECT deposits.id, deposits.deposit_no, deposits.kind, deposits.remaining_qty, deposits.yield_pct
FROM deposits
WHERE deposits.planet_id = :planet_id
ORDER BY deposits.deposit_no;
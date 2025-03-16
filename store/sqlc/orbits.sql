-- CreateOrbit creates a new orbit.
--
-- name: CreateOrbit :one
INSERT INTO orbits (star_id, orbit_no, kind)
VALUES (:star_id, :orbit_no, :kind)
RETURNING id;

-- DeleteEmptyOrbits deletes all orbits with no planets.
--
-- name: DeleteEmptyOrbits :exec
DELETE
FROM orbits
WHERE kind = 'EMPTY';


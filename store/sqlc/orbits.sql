-- CreateOrbit creates a new orbit.
--
-- name: CreateOrbit :one
INSERT INTO orbits (star_id, orbit_no, kind)
VALUES (:star_id, :orbit_no, :kind)
RETURNING id;

-- ReadOrbitPlanet returns the planet for a given orbit.
--
-- name: ReadOrbitPlanet :one
SELECT planets.id
FROM planets
WHERE planets.orbit_id = :orbit_id;

-- ReadOrbitSurvey reads the orbit survey data for a game.
--
-- name: ReadOrbitSurvey :many
SELECT systems.id             AS system_id,
       stars.id               AS star_id,
       orbits.id              AS orbit_id,
       orbits.orbit_no        AS orbit_no,
       planets.id             AS planet_id,
       planet_codes.name      AS planet_kind,
       deposits.id            AS deposit_id,
       deposits.deposit_no    AS deposit_no,
       unit_codes.code        AS deposit_kind,
       deposits.remaining_qty AS deposit_qty,
       deposits.yield_pct     AS yield_pct
FROM orbits,
     planets,
     planet_codes,
     deposits,
     unit_codes,
     stars,
     systems
WHERE orbits.id = :orbit_id
  AND planets.orbit_id = orbits.id
  AND planet_codes.code = planets.kind
  AND deposits.planet_id = planets.id
  AND unit_codes.code = deposits.kind
  AND orbits.id = planets.orbit_id
  AND stars.id = orbits.star_id
  AND systems.id = stars.system_id
ORDER BY deposits.deposit_no;

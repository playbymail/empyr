--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreatePlanet creates a new planet.
--
-- name: CreatePlanet :one
INSERT INTO planets (orbit_id, kind, habitability)
VALUES (:orbit_id, :kind, :habitability)
RETURNING id;

-- ReadPlanetSurvey reads the planet survey data for a game.
--
-- name: ReadPlanetSurvey :many
SELECT orbits.orbit_no        as orbit_no,
       planet_codes.name      as planet_kind,
       deposits.deposit_no    as deposit_no,
       unit_codes.name        as deposit_kind,
       deposits.remaining_qty as deposit_qty,
       deposits.yield_pct     as yield_pct
FROM orbits,
     planets,
     deposits,
     planet_codes,
     unit_codes
WHERE orbits.id = planets.orbit_id
  AND planets.id = :planet_id
  AND deposits.planet_id = planets.id
  AND planet_codes.code = planets.kind
  AND unit_codes.code = deposits.kind
ORDER BY deposits.deposit_no;

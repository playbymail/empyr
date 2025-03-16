--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreatePlanet creates a new planet.
--
-- name: CreatePlanet :one
INSERT INTO planets (orbit_id, kind, habitability)
VALUES (:orbit_id, :kind, :habitability)
RETURNING id;

-- -- ReadPlanetSurvey reads the planet survey data for a game.
-- --
-- -- name: ReadPlanetSurvey :many
-- SELECT systems.x              as x,
--        systems.y              as y,
--        systems.z              as z,
--        stars.sequence         as sequence,
--        orbits.orbit_no        as orbit_no,
--        planets.kind           as planet_kind,
--        deposits.deposit_no    as deposit_no,
--        deposits.kind          as deposit_kind,
--        deposits.remaining_qty as quantity
-- FROM systems,
--      stars,
--      orbits,
--      planets,
--      deposits
-- WHERE planets.id = :planet_id
--   AND stars.system_id = systems.id
--   AND orbits.star_id = stars.id
--   AND planets.orbit_id = orbits.id
--   AND deposits.planet_id = planets.id
-- ORDER BY deposits.deposit_no;

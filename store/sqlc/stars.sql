--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreateStar creates a new star in an existing system.
--
-- name: CreateStar :one
INSERT INTO stars (system_id, sequence)
VALUES (:system_id, :sequence)
RETURNING id;

-- ReadStarByID returns a star by its ID.
--
-- name: ReadStarByID :one
SELECT stars.id, stars.sequence, systems.x, systems.y, systems.z
FROM stars,
     systems
WHERE stars.id = :star_id
  AND systems.id = stars.system_id;

-- ReadAllStarsInCluster returns a list of all the stars in a cluster.
--
-- name: ReadAllStarsInCluster :many
SELECT systems.id     AS system_id,
       stars.id       AS star_id,
       stars.sequence AS sequence,
       systems.x      AS x,
       systems.y      AS y,
       systems.z      AS z
FROM clusters,
     systems,
     stars
WHERE clusters.id = :cluster_id
  AND systems.cluster_id = clusters.id
  AND stars.system_id = systems.id
ORDER BY systems.id, stars.sequence;

-- ReadAllStarsInSystem returns a list of stars in a system.
--
-- name: ReadAllStarsInSystem :many
SELECT stars.id, systems.x, systems.y, systems.z, stars.sequence
FROM systems,
     stars
WHERE systems.id = :system_id
  AND stars.system_id = systems.id
ORDER BY stars.id;

-- ReadStarSurvey reads the star survey data for star in a game.
--
-- name: ReadStarSurvey :many
SELECT orbits.id                   AS orbit_id,
       planets.id                  AS planet_id,
       orbits.kind                 AS orbit_kind,
       orbits.orbit_no             AS orbit_no,
       planets.kind                AS planet_kind,
       deposits.kind               AS deposit_kind,
       sum(deposits.remaining_qty) AS quantity
FROM stars,
     orbits,
     planets,
     deposits
WHERE stars.id = :star_id
  AND orbits.star_id = stars.id
  AND planets.orbit_id = orbits.id
  AND deposits.planet_id = planets.id
GROUP BY orbits.id, orbits.orbit_no, orbits.kind, planets.id, planets.kind, deposits.kind
ORDER BY orbits.orbit_no, deposits.kind;


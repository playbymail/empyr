--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- ReadAllStars reads the stars data for a game.
--
-- name: ReadAllStars :many
SELECT systems.id      AS system_id,
       stars.id        as star_id,
       stars.sequence  as sequence,
       systems.x       as x,
       systems.y       as y,
       systems.z       as z,
       count(stars.id) AS number_of_stars
FROM games
         LEFT JOIN systems on games.id = systems.game_id
         LEFT JOIN stars on systems.id = stars.system_id
WHERE games.id = :game_id
ORDER BY systems.id;

-- ReadStarSurvey reads the star survey data for a game.
--
-- name: ReadStarSurvey :many
SELECT systems.id                  AS system_id,
       stars.id                    as star_id,
       orbits.id                   as orbit_id,
       planets.id                  as planet_id,
       systems.x                   as x,
       systems.y                   as y,
       systems.z                   as z,
       stars.sequence              as sequence,
       orbits.kind                 as orbit_kind,
       orbits.orbit_no             as orbit_no,
       planets.kind                as planet_kind,
       deposits.kind               as deposit_kind,
       sum(deposits.remaining_qty) as quantity
FROM games,
     systems,
     stars,
     orbits,
     planets,
     deposits
WHERE stars.id = :star_id
  AND stars.system_id = systems.id
  AND systems.game_id = games.id
  AND orbits.star_id = stars.id
  AND planets.orbit_id = orbits.id
  AND deposits.planet_id = planets.id
GROUP BY systems.id, stars.id, orbits.id, planets.id, deposits.kind
ORDER BY systems.id, stars.id, orbits.id, planets.id, deposits.kind;

-- ReadStarsInSystem returns a list of stars in a system.
--
-- name: ReadStarsInSystem :many
SELECT stars.id
FROM systems,
     stars
WHERE systems.id = :system_id
  AND stars.system_id = systems.id
ORDER BY systems.id, stars.sequence;
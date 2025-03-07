--  Copyright (c) 2025 Michael D Henderson. All rights reserved.
--

-- CreateGame creates a new game.
--
-- name: CreateGame :one
INSERT INTO games (code, name, display_name)
VALUES (:code, :name, :display_name)
RETURNING id;

-- DeleteGame deletes an existing game
--
-- name: DeleteGame :exec
DELETE
FROM games
WHERE code = :code;

-- UpdateGameTurn increments the game turn number.
--
-- name: UpdateGameTurn :exec
UPDATE games
SET current_turn = :turn_number
WHERE id = :game_id;

-- GetCurrentGameTurn gets the current game turn.
--
-- name: GetCurrentGameTurn :one
SELECT current_turn
FROM games
WHERE id = :game_id;

-- CreateSystem creates a new system.
--
-- name: CreateSystem :one
INSERT INTO systems (game_id, x, y, z, scarcity)
VALUES (:game_id, :x, :y, :z, :scarcity)
RETURNING id;

-- CreateStar creates a new star in an existing system.
--
-- name: CreateStar :one
INSERT INTO stars (system_id, sequence, scarcity)
VALUES (:system_id, :sequence, :scarcity)
RETURNING id;

-- CreateSystemDistance inserts the distance between two systems.
--
-- name: CreateSystemDistance :exec
INSERT INTO system_distances (from_system_id, to_system_id, distance)
VALUES (:from_system_id, :to_system_id, :distance);

-- CreateOrbit creates a new orbit.
--
-- name: CreateOrbit :one
INSERT INTO orbits (star_id, orbit_no, kind, scarcity)
VALUES (:star_id, :orbit_no, :kind, :scarcity)
RETURNING id;

-- DeleteEmptyOrbits deletes all orbits with no planets.
--
-- name: DeleteEmptyOrbits :exec
DELETE
FROM orbits
WHERE kind = 'empty';

-- CreatePlanet creates a new planet.
--
-- name: CreatePlanet :one
INSERT INTO planets (orbit_id, kind, scarcity, habitability)
VALUES (:orbit_id, :kind, :scarcity, :habitability)
RETURNING id;

-- CreateDeposit creates a new deposit.
--
-- name: CreateDeposit :one
INSERT INTO deposits (planet_id, deposit_no, kind, yield_pct, initial_qty, remaining_qty)
VALUES (:planet_id, :deposit_no, :kind, :yield_pct, :initial_qty, :remaining_qty)
RETURNING id;

-- DeleteEmptyOrbits deletes all orbits with no planets.
--
-- name: DeleteEmptyDeposits :exec
DELETE
FROM deposits
WHERE kind = 'none';

-- CreateEmpire creates a new empire.
--
-- name: CreateEmpire :one
INSERT INTO empires (game_id, empire_no, name, home_system_id, home_star_id, home_orbit_id, home_planet_id)
VALUES (:game_id, :empire_no, :name, :home_system_id, :home_star_id, :home_orbit_id, :home_planet_id)
RETURNING id;

-- ReadClusterMap reads the cluster map.
--
-- name: ReadClusterMap :many
SELECT systems.id AS id,
       systems.x as x,
       systems.y as y,
       systems.z as z,
       count(stars.id) AS number_of_stars
FROM games
LEFT JOIN systems on games.id = systems.game_id
LEFT JOIN stars  on systems.id = stars.system_id
WHERE games.code = :game_code
GROUP BY systems.id, systems.x, systems.y, systems.z
ORDER BY systems.id;

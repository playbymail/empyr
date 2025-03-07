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
INSERT INTO orbits (star_id, orbit_no, kind)
VALUES (:star_id, :orbit_no, :kind)
RETURNING id;

-- DeleteEmptyOrbits deletes all orbits with no planets.
--
-- name: DeleteEmptyOrbits :exec
DELETE FROM orbits
WHERE kind = 'empty';

--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreateGame creates a new game.
--
-- name: CreateGame :one
INSERT INTO games (code, name, display_name)
VALUES (:code, :name, :display_name)
RETURNING id;

-- DeleteGameByID deletes an existing game using its ID.
--
-- name: DeleteGameByID :exec
DELETE
FROM games
WHERE id = :game_id;

-- DeleteGameByCode deletes an existing game using its code.
--
-- name: DeleteGameByCode :exec
DELETE
FROM games
WHERE code = :game_code;


-- ReadAllGameInfo returns a list of all games in the database, even the inactive ones.
--
-- name: ReadAllGameInfo :many
SELECT games.id,
       games.code,
       games.name,
       games.display_name,
       games.is_active,
       games.current_turn,
       games.last_empire_no
FROM games
ORDER BY games.code;

-- ReadAllGamesByUser returns all games that the user has a empire in.
--
-- name: ReadAllGamesByUser :many
SELECT games.id,
       games.code,
       games.name,
       games.display_name,
       games.current_turn,
       games.is_active,
       empires.id as empire_id,
       empires.empire_no
FROM empires,
     games
WHERE empires.user_id = :user_id
  AND games.id = empires.game_id
ORDER BY games.code;

-- ReadCurrentTurnByGameID gets the current turn for a game.
--
-- name: ReadCurrentTurnByGameID :one
SELECT current_turn
FROM games
WHERE id = :game_id;

-- ReadCurrentTurnByGameCode gets the current turn for a game.
--
-- name: ReadCurrentTurnByGameCode :one
SELECT current_turn
FROM games
WHERE code = :game_code;

-- ReadGameInfoByCode returns data for a single game using the game code.
--
-- name: ReadGameInfoByCode :one
SELECT games.id,
       games.code,
       games.name,
       games.display_name,
       games.is_active,
       games.current_turn,
       games.last_empire_no
FROM games
WHERE games.code = :game_code;

-- ReadGameSummaryByUserGame returns a game summary for a user in a game.
--
-- name: ReadGameSummaryByEmpire :one
SELECT games.id,
       games.code,
       games.name,
       games.display_name,
       games.is_active,
       games.current_turn,
       empires.id as empire_id,
       empires.empire_no
FROM empires,
     games
WHERE empires.user_id = :user_id
  AND games.id = empires.game_id
  AND games.id = :game_id
ORDER BY code;

-- ReadActiveGameSummariesByUser returns a list of all active games that the user has an empire in.
--
-- name: ReadActiveGameSummariesByUser :many
SELECT games.id,
       games.code,
       games.name,
       games.display_name,
       games.current_turn,
       empires.id as empire_id,
       empires.empire_no
FROM games,
     empires
WHERE empires.user_id = :user_id
  AND games.id = empires.game_id
  AND games.id = :game_id
  AND games.is_active = 1
ORDER BY code;

-- UpdateCurrentTurnByGameID increments the game turn number.
--
-- name: UpdateCurrentTurnByGameID :exec
UPDATE games
SET current_turn = :turn_number
WHERE id = :game_id;

-- UpdateCurrentTurnByGameCode increments the game turn number.
--
-- name: UpdateCurrentTurnByGameCode :exec
UPDATE games
SET current_turn = :turn_number
WHERE code = :game_code;

-- UpdateEmpireCounterByGameID updates the empire metadata in the games table.
--
-- name: UpdateEmpireCounterByGameID :exec
UPDATE games
SET last_empire_no = :empire_no
WHERE id = :game_id;

-- UpdateEmpireCounterByGameCode updates the empire metadata in the games table.
--
-- name: UpdateEmpireCounterByGameCode :exec
UPDATE games
SET last_empire_no = :empire_no
WHERE code = :game_code;



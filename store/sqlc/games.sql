--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- ReadAllGameInfo returns all games in the database, even the inactive ones.
--
-- name: ReadAllGameInfo :many
SELECT games.id,
       games.code,
       games.name,
       games.display_name,
       games.is_active,
       count(empires.id) as empire_count,
       count(empires.id) as player_count,
       games.current_turn,
       games.last_empire_no,
       games.home_system_id,
       games.home_star_id,
       games.home_orbit_id,
       games.home_planet_id
FROM games
         LEFT OUTER JOIN empires ON games.id = empires.game_id
ORDER BY code;

-- ReadUsersGames returns all active games that the user has a player in.
--
-- name: ReadUsersGames :many
SELECT games.id,
       games.code,
       games.name,
       games.display_name,
       games.is_active,
       games.current_turn,
       empires.id as empire_id,
       empires.empire_no
FROM games
         LEFT JOIN empires ON games.id = empires.game_id AND empires.user_id = :user_id
WHERE is_active = 1
ORDER BY code;

-- ReadEmpireGameSummary returns a summary of the empire's game.
--
-- name: ReadEmpireGameSummary :one
SELECT games.id,
       games.code,
       games.name,
       games.display_name,
       games.is_active,
       games.current_turn,
       empires.id as empire_id,
       empires.empire_no
FROM games
         LEFT JOIN empires ON games.id = empires.game_id AND empires.user_id = :user_id
WHERE games.code = :game_code
  AND games.is_active = 1
  AND empires.user_id = :user_id
  AND empires.empire_no = :empire_no
ORDER BY code;

-- GetListOfActiveGamesForUser returns a list of all active games
-- that the user is a player in.
--
-- name: GetListOfActiveGamesForUser :many
SELECT games.id,
       games.code,
       games.name,
       games.display_name,
       games.current_turn,
       empires.id as empire_id,
       empires.empire_no
FROM games,
     empires
WHERE games.is_active = 1
  AND empires.user_id = :user_id
  AND empires.game_id = games.id
ORDER BY code;

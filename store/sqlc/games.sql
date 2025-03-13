--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- ReadAllGames returns all games in the database, even the inactive ones.
--
-- name: ReadAllGames :many
SELECT id,
       code,
       name,
       display_name,
       is_active,
       current_turn,
       last_empire_no,
       home_system_id,
       home_star_id,
       home_orbit_id,
       home_planet_id
FROM games
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
       games.last_empire_no,
       games.home_system_id,
       games.home_star_id,
       games.home_orbit_id,
       games.home_planet_id
FROM games
         LEFT JOIN empires ON games.id = empires.game_id AND empires.user_id = :user_id
WHERE is_active = 1
ORDER BY code;

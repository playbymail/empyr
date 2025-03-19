-- ReadNextEmpireNumberByGame reads the next empire number in a game.
-- Has the side effect of incrementing the last_empire_no in that game by 1.
--
-- name: ReadNextEmpireNumberByGame :one
UPDATE games
SET last_empire_no = last_empire_no + 1
WHERE id = :game_id
RETURNING last_empire_no as next_empire_no;

-- CreateEmpire creates a new empire.
--
-- name: CreateEmpire :one
INSERT INTO empires (game_id,
                     user_id,
                     empire_no,
                     name,
                     home_system_id,
                     home_star_id,
                     home_orbit_id,
                     home_planet_id)
VALUES (:game_id,
        :user_id,
        :empire_no,
        :name,
        :home_system_id,
        :home_star_id,
        :home_orbit_id,
        :home_planet_id)
RETURNING id;


-- ReadAllEmpireCountsByGame returns the number of empires in all games.
--
-- name: ReadAllEmpireCountByGame :one
SELECT game_id, count(id)
FROM empires
GROUP BY game_id
ORDER BY game_id;


-- ReadAllEmpiresByGameCode returns the data for all empires in a game.
--
-- name: ReadAllEmpiresByGameCode :many
SELECT games.id AS game_id, empires.id AS empire_id, empire_no
FROM games,
     empires
WHERE games.code = :game_code
  AND empires.game_id = games.id
ORDER BY empire_no;

-- ReadAllEmpiresByGameID returns the data for all empires in a game.
--
-- name: ReadAllEmpiresByGameID :many
SELECT games.id AS game_id, empires.id AS empire_id, empire_no
FROM games,
     empires
WHERE games.id = :game_id
  AND empires.game_id = games.id
ORDER BY empire_no;

-- ReadEmpireCountByGameID returns the number of empires in a game.
--
-- name: ReadEmpireCountByGameID :one
SELECT count(empire_no) as empire_count
FROM empires
WHERE empires.game_id = :game_id
GROUP BY game_id;

-- ReadEmpireCountByGameCode returns the number of empires in a game.
--
-- name: ReadEmpireCountByGameCode :one
SELECT count(empire_no) as empire_count
FROM empires
WHERE empires.game_id = (SELECT id FROM games WHERE games.code = :game_code);


-- ReadEmpireByID reads an empire by its id.
-- This should only be used by admins. Regular users should use ReadEmpireByUserGame.
--
-- name: ReadEmpireByID :one
SELECT empires.game_id,
       empires.user_id,
       empires.empire_no,
       empires.name,
       empires.home_system_id,
       empires.home_star_id,
       empires.home_orbit_id,
       empires.home_planet_id
FROM empires
WHERE empires.id = :empire_id;

-- ReadEmpireByGameCodeByID returns the data for a single empire in a game.
--
-- name: ReadEmpireByGameCodeByID :one
SELECT games.id AS game_id, empires.id AS empire_id, empire_no
FROM games,
     empires
WHERE games.code = :game_code
  AND empires.game_id = games.id
  AND empires.empire_no = :empire_no;

-- ReadEmpireByGameCodeUserEmpireNo reads an empire for a user in an active game.
--
-- name: ReadEmpireByGameCodeUserEmpireNo :one
SELECT games.id           AS game_id,
       games.code         AS game_code,
       games.name         AS game_name,
       games.display_name AS game_display_name,
       games.current_turn AS game_current_turn,
       empires.user_id,
       empires.id         AS empire_id,
       empires.empire_no,
       empires.is_active  AS empire_is_active,
       empires.name       as empire_name,
       empires.home_system_id,
       empires.home_star_id,
       empires.home_orbit_id,
       empires.home_planet_id
FROM games,
     users,
     empires
WHERE games.code = :game_code
  AND games.is_active = 1
  AND users.id = :user_id
  AND empires.user_id = users.id
  AND empires.empire_no = :empire_no
  AND empires.is_active = 1;


-- ReadEmpireByGameIDByID returns the data for a single empire in a game.
--
-- name: ReadEmpireByGameIDByID :one
SELECT games.id           AS game_id,
       games.code         AS game_code,
       games.name         AS game_name,
       games.display_name AS game_display_name,
       games.current_turn AS game_current_turn,
       users.id           AS user_id,
       users.username     AS username,
       empires.id         AS empire_id,
       empires.empire_no,
       empires.is_active  AS empire_is_active,
       empires.name       AS empire_name,
       empires.home_system_id,
       empires.home_star_id,
       empires.home_orbit_id,
       empires.home_planet_id
FROM games,
     users,
     empires
WHERE games.id = :game_id
  AND empires.game_id = games.id
  AND empires.empire_no = :empire_no
  AND users.id = empires.user_id;

-- ReadEmpireByGameIDUserIDEmpireNo reads an empire for a user in a game.
--
-- name: ReadEmpireByGameIDUserIDEmpireNo :one
SELECT games.id           AS game_id,
       games.code         AS game_code,
       games.name         AS game_name,
       games.display_name AS game_display_name,
       games.current_turn AS game_current_turn,
       users.id           AS user_id,
       empires.id         AS empire_id,
       empires.empire_no,
       empires.is_active  AS empire_is_active,
       empires.name       AS empire_name,
       empires.home_system_id,
       empires.home_star_id,
       empires.home_orbit_id,
       empires.home_planet_id
FROM empires,
     games,
     users
WHERE games.id = :game_id
  AND users.id = :user_id
  AND empires.user_id = users.id
  AND empires.game_id = games.id
  AND empires.empire_no = :empire_no;

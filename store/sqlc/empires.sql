-- ReadNextEmpireNumber reads the next empire number.
--
-- name: ReadNextEmpireNumber :one
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


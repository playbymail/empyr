--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- ReadAllSystems reads the system data for a game.
--
-- name: ReadAllSystems :many
SELECT systems.id      AS id,
       systems.x       as x,
       systems.y       as y,
       systems.z       as z,
       count(stars.id) AS number_of_stars
FROM games
         LEFT JOIN systems on games.id = systems.game_id
         LEFT JOIN stars on systems.id = stars.system_id
WHERE games.id = :game_id
GROUP BY systems.id, systems.x, systems.y, systems.z
ORDER BY systems.id;

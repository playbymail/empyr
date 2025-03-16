--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreateSystem creates a new system.
--
-- name: CreateSystem :one
INSERT INTO systems (cluster_id, x, y, z)
VALUES (:cluster_id, :x, :y, :z)
RETURNING id;


-- ReadAllSystems reads the system data for a game.
--
-- name: ReadAllSystems :many
SELECT systems.id      AS id,
       systems.x       as x,
       systems.y       as y,
       systems.z       as z,
       count(stars.id) AS number_of_stars
FROM games,
     clusters,
     systems,
     stars
WHERE games.id = :game_id
  AND clusters.game_id = games.id
  AND systems.cluster_id = clusters.id
  AND stars.system_id = systems.id
GROUP BY systems.id, systems.x, systems.y, systems.z
ORDER BY systems.id;

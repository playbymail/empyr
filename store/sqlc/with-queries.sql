--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- ReadSystemNeighbors returns a list of systems and their closest neighbors.
--
-- name: ReadSystemNeighbors :many
WITH ranked_neighbors AS (SELECT g.code AS game_code,
                                 s1.id  AS source_id,
                                 s1.x   AS source_x,
                                 s1.y   AS source_y,
                                 s1.z   AS source_z,
                                 s2.id  AS neighbor_id,
                                 s2.x   AS neighbor_x,
                                 s2.y   AS neighbor_y,
                                 s2.z   AS neighbor_z,
                                 sd.distance,
                                 ROW_NUMBER() OVER (
                                     PARTITION BY s1.id
                                     ORDER BY sd.distance ASC
                                     )  AS rank
                          FROM systems s1
                                   JOIN system_distances sd ON s1.id = sd.from_system_id
                                   JOIN systems s2 ON sd.to_system_id = s2.id
                                   JOIN games g ON s1.game_id = g.id
                          WHERE g.code = :code          -- Parameter for game code
                            AND s1.game_id = g.id -- Ensure systems are in the same game
                            AND s2.game_id = g.id -- Ensure systems are in the same game
                            AND s1.id != s2.id -- Ensure systems are different
)
SELECT game_code,
       source_id,
       source_x,
       source_y,
       source_z,
       neighbor_id,
       neighbor_x,
       neighbor_y,
       neighbor_z,
       distance
FROM ranked_neighbors
WHERE rank = 1
ORDER BY distance, source_id;

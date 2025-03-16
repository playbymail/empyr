-- CreateCluster creates a new cluster.
--
-- name: CreateCluster :one
INSERT INTO clusters (game_id, home_system_id, home_star_id, home_orbit_id, home_planet_id)
VALUES (:game_id, :home_system_id, :home_star_id, :home_orbit_id, :home_planet_id)
RETURNING id;

-- ReadClusterMetaByGameCode reads the cluster metadata for a game
--
-- name: ReadClusterMetaByGameCode :one
SELECT clusters.id,
       clusters.home_system_id,
       clusters.home_star_id,
       clusters.home_orbit_id,
       clusters.home_planet_id
FROM clusters
WHERE clusters.game_id = (SELECT id FROM games WHERE code = :game_code);

-- ReadClusterMetaByGameID reads the cluster metadata for a game
--
-- name: ReadClusterMetaByGameID :one
SELECT clusters.id,
       clusters.home_system_id,
       clusters.home_star_id,
       clusters.home_orbit_id,
       clusters.home_planet_id
FROM clusters
WHERE clusters.game_id = :game_id;

-- ReadClusterMapByClusterID reads the cluster map for a game
--
-- name: ReadClusterMapByClusterID :many
SELECT systems.id      AS id,
       systems.x       as x,
       systems.y       as y,
       systems.z       as z,
       count(stars.id) AS number_of_stars
FROM clusters,
     systems,
     stars
WHERE clusters.id = :cluster_id
  AND systems.cluster_id = clusters.id
  AND stars.system_id = systems.id
GROUP BY systems.id, systems.x, systems.y, systems.z
ORDER BY systems.id;

-- ReadClusterMapByGameCode reads the cluster map for a game
--
-- name: ReadClusterMapByGameCode :many
SELECT systems.id      AS id,
       systems.x       as x,
       systems.y       as y,
       systems.z       as z,
       count(stars.id) AS number_of_stars
FROM games,
     clusters,
     systems,
     stars
WHERE games.code = :game_code
  AND clusters.game_id = games.id
  AND systems.cluster_id = clusters.id
  AND stars.system_id = systems.id
GROUP BY systems.id, systems.x, systems.y, systems.z
ORDER BY systems.id;

-- ReadClusterMapByGameID reads the cluster map for a game
--
-- name: ReadClusterMapByGameID :many
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


-- UpdateEmpireMetadataByClusterID updates the empire metadata in the clusters table.
--
-- name: UpdateEmpireMetadataByClusterID :exec
UPDATE clusters
SET home_system_id = :home_system_id,
    home_star_id   = :home_star_id,
    home_orbit_id  = :home_orbit_id,
    home_planet_id = :home_planet_id
WHERE id = :cluster_id;

-- PopulateSystemDistanceByCluster populates the system distance table with the
-- distance between all systems in the cluster.
--
-- name: PopulateSystemDistanceByCluster :exec
INSERT INTO system_distances (from_system_id, to_system_id, distance)
SELECT from_system.id,
       to_system.id,
       ceil(sqrt((from_system.x - to_system.x) * (from_system.x - to_system.x)
           + (from_system.y - to_system.y) * (from_system.y - to_system.y)
           + (from_system.z - to_system.z) * (from_system.z - to_system.z)))
FROM systems from_system,
     systems to_system
WHERE from_system.cluster_id = :cluster_id
  AND to_system.cluster_id = :cluster_id
  and from_system.id != to_system.id;


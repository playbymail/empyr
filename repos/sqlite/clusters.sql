-- ReadClusterMeta reads the cluster metadata for a game
--
-- name: ReadClusterMeta :one
select home_system_id,
       home_star_id,
       home_orbit_id
from games;

-- ReadClusterMap reads the cluster map for a game
--
-- name: ReadClusterMap :many
select id      as system_id,
       x       as x,
       y       as y,
       z       as z,
       nbr_of_stars
from systems
order by systems.id;

-- UpdateEmpireMetadata updates the empire metadata in the games table.
--
-- name: UpdateEmpireMetadata :exec
update games
set home_system_id = :home_system_id,
    home_star_id   = :home_star_id,
    home_orbit_id  = :home_orbit_id;


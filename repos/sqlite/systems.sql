--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreateSystem creates a new system.
--
-- name: CreateSystem :one
insert into systems (x, y, z, system_name, nbr_of_stars)
values (:x, :y, :z, :system_name, :nbr_of_stars)
returning id;

-- ReadAllSystems reads the system data for a game.
--
-- name: ReadAllSystems :many
select id as system_id,
       system_name,
       x  as x,
       y  as y,
       z  as z,
       nbr_of_stars
from systems
order by id;

--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreateStar creates a new star in an existing system.
--
-- name: CreateStar :one
insert into stars (system_id, sequence, star_name, nbr_of_orbits)
values (:system_id, :sequence, :star_name, :nbr_of_orbits)
returning id;

-- ReadStarByID returns a star by its ID.
--
-- name: ReadStarByID :one
select systems.x, systems.y, systems.z, star_name, stars.sequence
from stars,
     systems
where stars.id = :star_id
  and systems.id = stars.system_id;

-- ReadAllStarsInCluster returns a list of all the stars in a cluster.
--
-- name: ReadAllStarsInCluster :many
select systems.id     as system_id,
       system_name,
       stars.id       as star_id,
       stars.sequence as sequence,
       stars.star_name,
       systems.x      as x,
       systems.y      as y,
       systems.z      as z
from systems,
     stars
where stars.system_id = systems.id
order by systems.id, stars.sequence;

-- ReadAllStarsInSystem returns a list of stars in a system.
--
-- name: ReadAllStarsInSystem :many
select stars.id, systems.x, systems.y, systems.z, stars.sequence, stars.star_name
from systems,
     stars
where systems.id = :system_id
  and stars.system_id = systems.id
order by stars.sequence;

-- ReadStarSurvey reads the star survey data for star in a game.
--
-- name: ReadStarSurvey :many
select orbits.id         as orbit_id,
       orbits.orbit_no   as orbit_no,
       orbits.kind       as orbit_kind,
       deposits.kind     as deposit_kind,
       sum(deposits.qty) as quantity
from stars,
     orbits,
     deposits
where stars.id = :star_id
  and orbits.star_id = stars.id
  and deposits.orbit_id = orbits.id
group by orbits.id, orbits.orbit_no, orbits.kind, orbits.kind, deposits.kind
order by orbits.orbit_no, deposits.kind;


-- ReadNextEmpireNumber reads the next empire number in a game.
--
-- name: ReadNextEmpireNumber :one
select min(id) + 0 as empire_id
from empires
where is_active = 0;

-- CreateInactiveEmpires creates a new inactive empire.
--
-- name: CreateInactiveEmpires :exec
insert into empires (id, is_active)
values (:id, 0);

-- CreateEmpire creates a new empire.
--
-- name: CreateEmpire :exec
insert into empire (empire_id,
                    empire_name,
                    username,
                    email,
                    home_system_id,
                    home_star_id,
                    home_orbit_id)
values (:empire_id,
        :empire_name,
        :username,
        :email,
        :home_system_id,
        :home_star_id,
        :home_orbit_id);

-- IsEmpireActive checks if an empire is active.
--
-- name: IsEmpireActive :one
select is_active
from empires
where id = :empire_id;

-- ReadActiveEmpires returns the data for all active empires in a game.
--
-- name: ReadActiveEmpires :many
select id as empire_id
from empires
where is_active = 1
order by id;

-- ReadActiveEmpireCount returns the number of active empires.
--
-- name: ReadActiveEmpireCount :one
select count(id)
from empires
where is_active = 1;


-- ReadEmpireByID reads an empire by its id.
-- This should only be used by admins. Regular users should use ReadEmpireByUser.
--
-- name: ReadEmpireByID :one
select games.code         as game_code,
       games.name         as game_name,
       games.display_name as game_display_name,
       games.current_turn as game_current_turn,
       empire.empire_id,
       empire.empire_name,
       empire.username,
       empire.email,
       empire.home_system_id,
       empire.home_star_id,
       empire.home_orbit_id
from games,
     empires,
     empire
where empires.id = :empire_id
  and empire.empire_id = empires.id
  and is_active = 1;

-- UpdateEmpireStatus updates the status of an empire.
--
-- name: UpdateEmpireStatus :exec
update empires
set is_active = :is_active
where id = :empire_id;
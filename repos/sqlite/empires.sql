-- ReadNextEmpireNumber reads the next empire number in a game.
--
-- name: ReadNextEmpireNumber :one
select min(id) + 0 as empire_id
from empire
where is_active = 0;

-- CreateEmpire creates a new empire and returns the new empire ID.
--
-- name: CreateEmpire :one
insert into empire (home_system_id, home_star_id, home_orbit_id, is_active)
values (:home_system_id, :home_star_id, :home_orbit_id, :is_active)
returning id;

-- CreateEmpireWithID creates a new empire with the given ID.
--
-- name: CreateEmpireWithID :exec
insert into empire (id, home_system_id, home_star_id, home_orbit_id, is_active)
values (:id, :home_system_id, :home_star_id, :home_orbit_id, :is_active);

-- CreateEmpirePlayer creates a new empire player record.
--
-- name: CreateEmpirePlayer :exec
insert into empire_player (empire_id, effdt, enddt, username, email)
values (:empire_id, :effdt, :enddt, :username, :email);

-- CreateEmpireSystemName creates a new empire system name record.
--
-- name: CreateEmpireSystemName :exec
insert into empire_system_name (empire_id, system_id, effdt, enddt, name)
values (:empire_id, :system_id, :effdt, :enddt, :name);

-- UpdateEmpirePlayerEndDt updates the end date for an empire player.
--
-- name: UpdateEmpirePlayerEndDt :exec
update empire_player
set enddt = :enddt
where empire_id = :empire_id
  and effdt = :effdt;

-- IsEmpireActive checks if an empire is active.
--
-- name: IsEmpireActive :one
select is_active
from empire
where id = :empire_id;

-- ReadActiveEmpires returns the data for all active empires in a game.
--
-- name: ReadActiveEmpires :many
select id as empire_id
from empire
where is_active = 1
order by id;

-- ReadActiveEmpireCount returns the number of active empire.
--
-- name: ReadActiveEmpireCount :one
select count(id)
from empire
where is_active = 1;


-- ReadEmpireByID reads an empire by its id.
-- This should only be used by admins. Regular users should use ReadEmpireByUser.
--
-- name: ReadEmpireByID :one
select games.code          as game_code,
       games.name          as game_name,
       games.display_name  as game_display_name,
       games.current_turn  as game_current_turn,
       empire.id           as empire_id,
       empire_name.name    as empire_name,
       empire_player.username,
       empire_player.email,
       empire.home_system_id,
       systems.system_name as home_system_name,
       empire.home_star_id,
       stars.star_name     as home_star_name,
       empire.home_orbit_id,
       orbits.orbit_no
from empire,
     empire_name,
     empire_player,
     systems,
     stars,
     orbits,
     games
where empire.id = :empire_id
  and empire.is_active = 1
  and empire_name.empire_id = empire.id
  and (empire_name.effdt <= :as_of_dt and :as_of_dt <= empire_name.enddt)
  and empire_player.empire_id = empire.id
  and (empire_player.effdt <= :as_of_dt and :as_of_dt <= empire_player.enddt)
  and systems.id = empire.home_system_id
  and stars.id = empire.home_star_id
  and orbits.id = empire.home_orbit_id;

-- UpdateEmpireStatus updates the status of an empire.
--
-- name: UpdateEmpireStatus :exec
update empire
set is_active = :is_active
where id = :empire_id;

-- CreateEmpireName adds a new empire name record.
--
-- name: CreateEmpireName :exec
insert into empire_name(empire_id, effdt, enddt, name)
values (:empire_id, :effdt, :enddt, :name);

-- CorrectEmpireName updates an existing record for an empire name.
--
-- name: CorrectEmpireName :exec
update empire_name
set name = :name
where empire_id = :empire_id
  and effdt = :effdt
  and enddt = :enddt;


-- ExpireEmpireName updates the end date for an empire name.
--
-- name: ExpireEmpireName :exec
update empire_name
set enddt = :enddt
where empire_id = :empire_id
  and effdt = :effdt;

-- ReadEmpireName returns the name of an empire as of the given date.
--
-- name: ReadEmpireName :one
select name, effdt, enddt
from empire_name
where empire_id = :empire_id
  and (empire_name.effdt <= :as_of_dt and :as_of_dt <= empire_name.enddt);



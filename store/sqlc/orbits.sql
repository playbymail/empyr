-- CreateOrbit creates a new orbit.
--
-- name: CreateOrbit :one
insert into orbits (system_id, star_id, orbit_no, kind, habitability, nbr_of_deposits)
values (:system_id, :star_id, :orbit_no, :kind, :habitability, :nbr_of_deposits)
returning id;

-- ReadOrbitStar returns the star for a given orbit.
--
-- name: ReadOrbitStar :one
select systems.id as system_id,
       systems.system_name,
       stars.id   as star_id,
       stars.star_name,
       orbits.orbit_no
from orbits,
     stars,
     systems
where orbits.id = :orbit_id
  and stars.id = orbits.star_id
  and systems.id = orbits.system_id;

-- ReadOrbitSurvey reads the orbit survey data for a game.
--
-- name: ReadOrbitSurvey :many
select systems.id          as system_id,
       systems.system_name,
       stars.id            as star_id,
       stars.star_name,
       orbits.orbit_no     as orbit_no,
       orbit_codes.name    as orbit_kind,
       deposits.id         as deposit_id,
       deposits.deposit_no as deposit_no,
       unit_codes.code     as deposit_kind,
       deposits.qty        as deposit_qty,
       deposits.yield_pct  as yield_pct
from orbits,
     orbit_codes,
     deposits,
     unit_codes,
     stars,
     systems
where orbits.id = :orbit_id
  and deposits.orbit_id = orbits.id
  and unit_codes.code = deposits.kind
  and stars.id = orbits.star_id
  and systems.id = orbits.system_id
order by deposits.deposit_no;

-- UpdateOrbit updates an existing orbit.
--
-- name: UpdateOrbit :exec
update orbits
set kind            = :kind,
    habitability    = :habitability,
    nbr_of_deposits = :nbr_of_deposits
where id = :orbit_id;

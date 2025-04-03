--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- ExportCoverTabByID returns the exportable data for a single empire's cover tab.
--
-- name: ExportCoverTabByID :one
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

-- ExportStarProbes returns the exportable data for all star probes by an empire.
--
-- name: ExportStarProbes :many
select scs.id,
       sc_probe_star_result.probe_id,
       sc_probe_star_result.effdt as as_of_dt,
       sc_probe_star_result.location,
       sc_probe_star_result.nbr_of_orbits
from scs,
     sc_probe_order,
     sc_probe_star_result
where scs.empire_id = :empire_id
  and sc_probe_order.sc_id = scs.id
  and sc_probe_star_result.probe_id = sc_probe_order.id
  and sc_probe_star_result.effdt = :as_of_dt
order by scs.id,
         sc_probe_star_result.location,
         sc_probe_order.id;

-- ExportSystems returns the exportable data for all systems in a game.
--
-- name: ExportSystems :many
select system_name, x, y, z, nbr_of_stars
from systems
order by system_name;
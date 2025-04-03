-- CreateSC creates a new colony and returns its ID.
--
-- name: CreateSC :one
insert into scs (empire_id, sc_cd, sc_tech_level)
values (:empire_id, :sc_cd, :sc_tech_level)
returning id;

-- CreateSCRates creates a new colony.
--
-- name: CreateSCRates :exec
insert into sc_rates (sc_id, effdt, enddt, rations, sol, birth_rate, death_rate)
values (:sc_id, :effdt, :enddt, :rations, :sol, :birth_rate, :death_rate);

-- CreateSCInventory creates a new colony inventory entry.
--
-- name: CreateSCInventory :exec
insert into sc_inventory (sc_id, unit_cd, unit_tech_level, effdt, enddt, qty, mass, volume, is_assembled, is_stored)
values (:sc_id, :unit_cd, :unit_tech_level, :effdt, :enddt, :qty, :mass, :volume, :is_assembled, :is_stored);

-- UpdateSCInventoryEndDt updates the end date for an inventory entry.
--
-- name: UpdateSCInventoryEndDt :exec
update sc_inventory
set enddt = :enddt
where sc_id = :sc_id
  and unit_cd = :unit_cd
  and unit_tech_level = :unit_tech_level
  and effdt = :effdt;

-- CreateSCLocation creates a new colony location entry.
--
-- name: CreateSCLocation :exec
insert into sc_location (sc_id, effdt, enddt, orbit_id, is_on_surface)
values (:sc_id, :effdt, :enddt, :orbit_id, :is_on_surface);

-- UpdateSCLocationEndDt updates the end date for a location entry.
--
-- name: UpdateSCLocationEndDt :exec
update sc_location
set enddt = :enddt
where sc_id = :sc_id
  and effdt = :effdt;

-- CreateSCName creates a new colony name entry.
--
-- name: CreateSCName :exec
insert into sc_name (sc_id, name, effdt, enddt)
values (:sc_id, :name, :effdt, :enddt);

-- UpdateSCNameEndDt updates the end date for a name entry.
--
-- name: UpdateSCNameEndDt :exec
update sc_name
set enddt = :enddt
where sc_id = :sc_id
  and effdt = :effdt;

-- CreateSCPopulation creates a new colony population entry.
--
-- name: CreateSCPopulation :exec
insert into sc_population (sc_id, population_cd, effdt, enddt, qty, pay_rate, rebel_qty)
values (:sc_id, :population_cd, :effdt, :enddt, :qty, :pay_rate, :rebel_qty);

-- UpdateSCPopulationEndDt updates the end date for a population entry.
--
-- name: UpdateSCPopulationEndDt :exec
update sc_population
set enddt = :enddt
where sc_id = :sc_id
  and population_cd = :population_cd
  and effdt = :effdt;

-- CreateSCGroup creates a new ship or colony production group.
--
-- name: CreateSCGroup :exec
insert into sc_group (sc_id, kind, effdt, enddt)
values (:sc_id, :kind, :effdt, :enddt);

-- CreateSCGroupNo creates a new ship or colony production group number.
-- The number must be between 1 and 35 (or 1 and 40 for deposits). The
-- caller must ensure that the number is distinct for the group and the
-- date range.
--
-- name: CreateSCGroupNo :exec
insert into sc_group_no (group_id, effdt, enddt, group_no)
values (:group_id, :effdt, :enddt, :group_no);

-- CreateSCGroupTooling creates a record to change the tooling of a factory
-- group. If the retooled flag is set, the factories will stop producing
-- new items for a few turns.
--
-- name: CreateSCGroupTooling :exec
insert into sc_group_tooling (group_id, effdt, enddt, item_cd, item_tech_level, retooled)
values (:group_id, :effdt, :enddt, :item_cd, :item_tech_level, :retooled);

-- CreateSCGroupUnit creates a new set of units for the group. All the units
-- in a group must have the same technology level. Use this to create a new
-- set of units or change the number of units.
--
-- name: CreateSCGroupUnit :exec
insert into sc_group_unit (group_id, tech_level, effdt, enddt, nbr_of_units)
values (:group_id, :tech_level, :effdt, :enddt, :nbr_of_units);

-- CreateSCGroupUnitProduction creates a record of the resources consumed
-- and created by the units in a single turn.
--
-- name: CreateSCGroupUnitProduction :exec
insert into sc_group_unit_production (group_id, tech_level, production_dt,
                                      fuel_consumed, gold_consumed, mets_consumed, nmts_consumed,
                                      pro_consumed, usk_consumed, aut_consumed,
                                      qty_produced)
values (:group_id, :tech_level, :production_dt,
        :fuel_consumed, :gold_consumed, :mets_consumed, :nmts_consumed,
        :pro_consumed, :usk_consumed, :aut_consumed,
        :qty_produced);

-- CreateSCGroupUnitProductionWIP creates a record of the work in progress for a group
-- of factory units in a single turn.
--
-- name: CreateSCGroupUnitProductionWIP :exec
insert into sc_group_unit_production_wip (group_id, tech_level, production_dt,
                                          wip_25pct_qty, wip_50pct_qty, wip_75pct_qty)
values (:group_id, :tech_level, :production_dt,
        :wip_25pct_qty, :wip_50pct_qty, :wip_75pct_qty);

-- CreateSCGroupProductionSummary stores the totals from the entire group for a single turn.
--
-- name: CreateSCGroupProductionSummary :exec
insert into sc_group_production_summary (group_id, production_dt,
                                         fuel_consumed, gold_consumed, mets_consumed, nmts_consumed,
                                         pro_consumed, usk_consumed, aut_consumed,
                                         qty_produced)
values (:group_id, :production_dt,
        :fuel_consumed, :gold_consumed, :mets_consumed, :nmts_consumed,
        :pro_consumed, :usk_consumed, :aut_consumed,
        :qty_produced);

-- CreateSCGroupProductionWIPSummary stores the total WIP from the entire group for a single turn.
--
-- name: CreateSCGroupProductionWIPSummary :exec
insert into sc_group_production_wip_summary (group_id, production_dt,
                                             wip_25pct_qty, wip_50pct_qty, wip_75pct_qty)
values (:group_id, :production_dt,
        :wip_25pct_qty, :wip_50pct_qty, :wip_75pct_qty);

-- CreateSCMiningSummary creates a summary of resources mined by a colony in a single turn.
--
-- name: CreateSCMiningSummary :exec
insert into sc_mining_summary (sc_id, production_dt,
                               fuel_produced, gold_produced, mets_produced, nmts_produced)
values (:sc_id, :group_no,
        :fuel_produced, :gold_produced, :mets_produced, :nmts_produced);


-- CreateSCProbeOrder creates a new colony probe order.
--
-- name: CreateSCProbeOrder :one
insert into sc_probe_order (sc_id, effdt, target_id, kind)
values (:sc_id, :effdt, :target_id, :kind)
returning id;

-- CreateSCProbeStarResults adds a new result.
--
-- name: CreateSCProbeStarResults :exec
insert into sc_probe_star_result(probe_id, effdt, star_id, location, nbr_of_orbits)
values (:probe_id, :effdt, :star_id, :location, :nbr_of_orbits);

-- DeleteSCProbeStarResult deletes the results of a probe.
--
-- name: DeleteSCProbeStarResult :exec
delete
from sc_probe_star_result
where probe_id = :probe_id;

-- DeleteSCProbeStarResultsByTurn deletes the results of all probes for a given turn.
--
-- name: DeleteSCProbeStarResultsByTurn :exec
delete
from sc_probe_star_result
where effdt = :effdt;

-- CreateSCProbeStarOrbitResults adds a new result.
--
-- name: CreateSCProbeStarOrbitResults :exec
insert into sc_probe_star_orbit_result(probe_id, effdt,
                                       star_id,
                                       orbit_no, orbit_kind,
                                       fuel_est, gold_est, mets_est, nmts_est)
values (:probe_id, :effdt,
        :star_id,
        :orbit_no, :orbit_kind,
        :fuel_est, :gold_est, :mets_est, :nmts_est);

-- DeleteSCProbeStarOrbitResults deletes the results of a probe.
--
-- name: DeleteSCProbeStarOrbitResult :exec
delete
from sc_probe_star_orbit_result
where probe_id = :probe_id;

-- DeleteSCProbeStarOrbitResultsByTurn deletes the results of all probes for a given turn.
--
-- name: DeleteSCProbeStarOrbitResultsByTurn :exec
delete
from sc_probe_star_orbit_result
where effdt = :effdt;

-- CreateSCSurveyOrder creates a new colony survey order.
--
-- name: CreateSCSurveyOrder :one
insert into sc_survey_order (sc_id, effdt, target_id, kind)
values (:sc_id, :effdt, :target_id, :kind)
returning id;

-- CreateSCSurveyOrbitResult adds a new result.
--
-- name: CreateSCSurveyOrbitResult :exec
insert into sc_survey_orbit_result (survey_id, effdt, orbit_id, location, orbit_no, habitability_no, farmland_in_use,
                                    population)
values (:survey_id, :effdt, :orbit_id, :location, :orbit_no, :habitability_no, :farmland_in_use, :population);

-- DeleteSCSurveyOrbitResult deletes the results of a survey.
--
-- name: DeleteSCSurveyOrbitResult :exec
delete
from sc_survey_orbit_result
where survey_id = :survey_id;

-- DeleteSCSurveyOrbitResultsByTurn deletes the results of all surveys for a given turn.
--
-- name: DeleteSCSurveyOrbitResultsByTurn :exec
delete
from sc_survey_orbit_result
where effdt = :effdt;

-- ReadAllColoniesByEmpire returns a list of all colonies for an empire
-- that were active on a given turn.
--
-- name: ReadAllColoniesByEmpire :many
select scs.id        as sc_id,
       systems.id    as system_id,
       systems.system_name,
       stars.id      as star_id,
       stars.star_name,
       orbits.orbit_no,
       sc_codes.name as sc_kind,
       scs.sc_tech_level,
       sc_name.name,
       sc_rates.rations,
       sc_rates.birth_rate,
       sc_rates.death_rate,
       sc_rates.sol
from scs,
     sc_codes,
     sc_name,
     sc_rates,
     sc_location,
     orbits,
     stars,
     systems
where scs.empire_id = :empire_id
  and scs.sc_cd in ('COPN', 'CENC', 'CORB')
  and sc_codes.code = scs.sc_cd
  and sc_name.sc_id = scs.id
  and (sc_name.effdt <= :as_of_dt and :as_of_dt < sc_name.enddt)
  and sc_rates.sc_id = scs.id
  and (sc_rates.effdt <= :as_of_dt and :as_of_dt < sc_rates.enddt)
  and sc_location.sc_id = scs.id
  and (sc_location.effdt <= :as_of_dt and :as_of_dt < sc_location.enddt)
  and orbits.id = sc_location.orbit_id
  and stars.id = orbits.star_id
  and systems.id = orbits.system_id
order by scs.id;

-- ReadAllSurveyOrdersByTurn returns a list of survey orders issued in a given turn of a game.
--
-- name: ReadAllSurveyOrdersByTurn :many
select empire.id as empire_id,
       sc_survey_order.sc_id,
       sc_survey_order.effdt,
       sc_survey_order.target_id
from sc_survey_order,
     scs,
     empire
where sc_survey_order.effdt = :as_of_dt
  and scs.id = sc_survey_order.sc_id
  and empire.id = scs.empire_id
order by empire.id, sc_survey_order.sc_id, sc_survey_order.effdt, sc_survey_order.target_id;

-- ReadSCPopulation returns a list of the population for a given colony.
--
-- name: ReadSCPopulation :many
select sc_population.population_cd,
       population_codes.name as population_kind,
       sc_population.qty,
       sc_population.pay_rate,
       sc_population.rebel_qty
from scs,
     sc_population,
     population_codes
where scs.id = :sc_id
  and sc_population.sc_id = scs.id
  and (sc_population.effdt <= :as_of_dt and :as_of_dt < sc_population.enddt)
  and population_codes.code = sc_population.population_cd
order by population_codes.sort_order;

-- ReadSCInventory returns a list of the inventory for a given colony.
--
-- name: ReadSCInventory :many
select sc_inventory.unit_cd,
       sc_inventory.unit_tech_level,
       unit_codes.name as unit_kind,
       sc_inventory.qty,
       sc_inventory.mass,
       sc_inventory.volume,
       sc_inventory.is_assembled,
       sc_inventory.is_stored
from sc_inventory,
     unit_codes
where sc_inventory.sc_id = :sc_id
  and (sc_inventory.effdt <= :as_of_dt and :as_of_dt < sc_inventory.enddt)
  and unit_codes.code = sc_inventory.unit_cd
order by sc_inventory.unit_cd, sc_inventory.unit_tech_level, sc_inventory.qty;

-- ReadSCGroups returns a list of the groups for a given ship or colony.
--
-- name: ReadSCGroups :many
select sc_group.id as group_id,
       sc_group_no.group_no
from scs,
     sc_group,
     sc_group_no
where scs.id = :sc_id
  and sc_group.sc_id = scs.id
  and sc_group.kind = :kind
  and (sc_group.effdt <= :as_of_dt and :as_of_dt < sc_group.enddt)
  and sc_group_no.group_id = sc_group.id
  and (sc_group_no.effdt <= :as_of_dt and :as_of_dt < sc_group_no.enddt)
order by sc_group_no.group_no;

-- ReadSCGroupTooling returns a list of the factory groups for a given colony.
--
-- name: ReadSCGroupTooling :many
select sc_group.id as group_id,
       sc_group_no.group_no,
       sc_group_tooling.item_cd,
       sc_group_tooling.item_tech_level,
       sc_group_tooling.retooled
from scs,
     sc_group,
     sc_group_no,
     sc_group_tooling
where scs.id = :sc_id
  and sc_group.sc_id = scs.id
  and sc_group.kind = 'factory'
  and (sc_group.effdt <= :as_of_dt and :as_of_dt < sc_group.enddt)
  and sc_group_no.group_id = sc_group.id
  and (sc_group_no.effdt <= :as_of_dt and :as_of_dt < sc_group_no.enddt)
  and sc_group_tooling.group_id = sc_group.id
  and (sc_group_tooling.effdt <= :as_of_dt and :as_of_dt < sc_group_tooling.enddt)
order by sc_group_no.group_no;

-- ReadSCGroupUnits returns a list of the group units for a colony.
--
-- name: ReadSCGroupUnits :many
select sc_group_unit.tech_level,
       nbr_of_units
from sc_group,
     sc_group_unit
where sc_group.sc_id = :sc_id
  and sc_group.kind = :kind
  and (sc_group.effdt <= :as_of_dt and :as_of_dt < sc_group.enddt)
  and sc_group_unit.group_id = sc_group.id
  and (sc_group_unit.effdt <= :as_of_dt and :as_of_dt < sc_group_unit.enddt)
order by sc_group_unit.tech_level, nbr_of_units;

-- ReadSCProbeOrders returns a list of probe orders issued by a colony on a given turn.
--
-- name: ReadSCProbeOrders :exec
select target_id, kind
from sc_probe_order
where sc_id = :sc_id
  and effdt = :effdt
order by target_id, kind;

-- ReadSCSurveyOrders returns a list of survey orders issued by a colony on a given turn.
--
-- name: ReadSCSurveyOrders :exec
select target_id
from sc_survey_order
where sc_id = :sc_id
  and effdt = :effdt
order by target_id;

-- CreateSC creates a new ship or colony.
--
-- name: CreateSC :one
insert into scs (empire_id, sc_cd, sc_tech_level, name, location, is_on_surface, rations, sol, birth_rate, death_rate)
values (:empire_id, :sc_cd, :sc_tech_level, :name, :location, :is_on_surface, :rations, :sol, :birth_rate, :death_rate)
returning id;

-- CreateSCInventory creates a new ship or colony inventory entry.
--
-- name: CreateSCInventory :exec
insert into inventory (sc_id, unit_cd, tech_level, qty, mass, volume, is_assembled, is_stored)
values (:sc_id, :unit_cd, :tech_level, :qty, :mass, :volume, :is_assembled, :is_stored);

-- CreateSCPopulation creates a new ship or colony population entry.
--
-- name: CreateSCPopulation :exec
insert into population (sc_id, population_cd, qty, pay_rate, rebel_qty)
values (:sc_id, :population_cd, :qty, :pay_rate, :rebel_qty);

-- CreateSCFactoryGroup creates a new ship or colony factory group.
--
-- name: CreateSCFactoryGroup :exec
insert into factory_groups (sc_id, group_no, orders_cd, orders_tech_level)
values (:sc_id, :group_no, :orders_cd, :orders_tech_level);

-- CreateSCFactoryGroupUnit creates a new ship or colony factory group unit.
--
-- name: CreateSCFactoryGroupUnit :exec
insert into factory_group (sc_id, group_no, group_tech_level, nbr_of_units,
                           orders_cd, orders_tech_level,
                           wip_25pct_qty, wip_50pct_qty, wip_75pct_qty)
values (:sc_id, :group_no, :group_tech_level, :nbr_of_units,
        :orders_cd, :orders_tech_level,
        :wip_25pct_qty, :wip_50pct_qty, :wip_75pct_qty);

-- CreateSCFarmGroup creates a new ship or colony farm group.
--
-- name: CreateSCFarmGroup :exec
insert into farm_groups (sc_id, group_no)
values (:sc_id, :group_no);

-- CreateSCFarmGroupUnit creates a new ship or colony farm group unit.
--
-- name: CreateSCFarmGroupUnit :exec
insert into farm_group (sc_id, group_no, group_tech_level, nbr_of_units)
values (:sc_id, :group_no, :group_tech_level, :nbr_of_units);

-- CreateSCMiningGroup creates a new ship or colony mining group.
--
-- name: CreateSCMiningGroup :exec
insert into mining_groups (sc_id, group_no, deposit_id)
values (:sc_id, :group_no, :deposit_id);

-- CreateSCMiningGroupUnit creates a new ship or colony mining group unit.
--
-- name: CreateSCMiningGroupUnit :exec
insert into mining_group (sc_id, group_no, group_tech_level, nbr_of_units)
values (:sc_id, :group_no, :group_tech_level, :nbr_of_units);

-- CreateSCProbeOrder creates a new ship or colony probe order.
--
-- name: CreateSCProbeOrder :exec
insert into probe_orders (sc_id, kind, target_id, status)
values (:sc_id, :kind, :target_id, 'ordered');

-- CreateSCSurveyOrder creates a new ship or colony survey order.
--
-- name: CreateSCSurveyOrder :exec
insert into survey_orders (sc_id, target_id, status)
values (:sc_id, :target_id, 'ordered');

-- ReadAllColoniesByEmpire reads the colonies for a given empire in a game.
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
       scs.name,
       scs.rations,
       scs.birth_rate,
       scs.death_rate,
       scs.sol
from scs,
     orbits,
     stars,
     systems,
     sc_codes
where scs.empire_id = :empire_id
  and scs.sc_cd in ('COPN', 'CENC', 'CORB')
  and orbits.id = scs.location
  and stars.id = orbits.star_id
  and systems.id = orbits.system_id
  and sc_codes.code = scs.sc_cd
order by scs.id;

-- ReadSCPopulation reads the population for a given ship or colony.
--
-- name: ReadSCPopulation :many
select population.population_cd,
       population_codes.name as population_kind,
       population.qty,
       population.pay_rate,
       population.rebel_qty
from population,
     population_codes
where population.sc_id = :sc_id
  and population_codes.code = population.population_cd
order by population_codes.sort_order;

-- ReadSCInventory reads the inventory for a given ship or colony.
--
-- name: ReadSCInventory :many
select inventory.unit_cd,
       inventory.tech_level,
       unit_codes.name as unit_kind,
       inventory.qty,
       inventory.mass,
       inventory.volume,
       inventory.is_assembled,
       inventory.is_stored
from inventory,
     unit_codes
where inventory.sc_id = :sc_id
  and unit_codes.code = inventory.unit_cd
order by inventory.unit_cd, inventory.tech_level, inventory.qty;

-- ReadSCFactoryGroups reads the factory groups for a given ship or colony.
--
-- name: ReadSCFactoryGroups :many
SELECT sc_id,
       group_no,
       orders_cd,
       orders_tech_level
FROM factory_groups
WHERE sc_id = :sc_id
ORDER BY sc_id, group_no;

-- ReadSCFactoryGroup reads the factory group for a given ship or colony.
--
-- name: ReadSCFactoryGroup :many
select orders_cd,
       orders_tech_level,
       group_tech_level,
       nbr_of_units,
       wip_25pct_qty,
       wip_50pct_qty,
       wip_75pct_qty
from factory_group
where factory_group.sc_id = :sc_id
  and factory_group.group_no = :group_no
order by orders_cd, orders_tech_level, group_tech_level;

-- ReadSCFactoryGroupRetoolOrder reads the factory group retooling for a given ship or colony.
--
-- name: ReadSCFactoryGroupRetoolOrder :one
select sc_id, group_no, turn_no, orders_cd, orders_tech_level
from factory_group_retool
where sc_id = :sc_id
  and group_no = :group_no;

-- ReadSCFarmGroups reads the farm groups for a given ship or colony.
--
-- name: ReadSCFarmGroups :many
select group_no,
       group_tech_level,
       sum(nbr_of_units) as nbr_of_units
from farm_group
where farm_group.sc_id = :sc_id
group by group_no, group_tech_level
order by group_no, group_tech_level;

-- ReadSCMiningGroups reads the mining groups for a given ship or colony.
--
-- name: ReadSCMiningGroups :many
select mining_groups.sc_id,
       mining_groups.group_no,
       deposits.deposit_no,
       deposits.qty,
       deposits.kind as deposit_kind,
       deposits.yield_pct
from mining_groups,
     deposits
where mining_groups.sc_id = :sc_id
  and deposits.id = mining_groups.deposit_id
order by mining_groups.sc_id, mining_groups.group_no, deposits.deposit_no;

-- ReadSCMiningGroup returns the data for a given mining group.
--
-- name: ReadSCMiningGroup :many
select sc_id,
       group_no,
       group_tech_level,
       sum(mining_group.nbr_of_units) as nbr_of_units
from mining_group
where sc_id = :sc_id
  and group_no = :group_no
group by sc_id, group_no, group_tech_level
order by sc_id, group_no, group_tech_level;

-- ReadSCProbeOrders returns a list of probe orderss issued by a ship or colony on a given turn.
--
-- name: ReadSCProbeOrders :exec
select kind, target_id, status
from probe_orders
where sc_id = :sc_id
order by kind, target_id, status;

-- ReadSCSurveyOrders returns a list of survey orders issued by a ship or colony on a given turn.
--
-- name: ReadSCSurveyOrders :exec
select target_id, status
from survey_orders
where sc_id = :sc_id
order by target_id, status;

-- ReadAllSurveyOrdersGameForTurn returns a list of survey orders issued in a given turn of a game.
--
-- name: ReadAllSurveyOrdersForGameTurn :many
select empires.id as empire_id,
       survey_orders.sc_id,
       survey_orders.target_id,
       survey_orders.status
from empires,
     scs,
     survey_orders
where scs.empire_id = empires.id
  and survey_orders.sc_id = scs.id
order by empires.id, survey_orders.sc_id, survey_orders.target_id, survey_orders.status;

-- ResetProbeOrdersStatus resets the status of all probe orders issued in a given turn.
--
-- name: ResetProbeOrdersStatus :exec
update probe_orders
set status = 'ordered';

-- ResetSurveyOrdersStatus resets the status of all survey orders issued in a given turn.
--
-- name: ResetSurveyOrdersStatus :exec
update survey_orders
set status = 'ordered';

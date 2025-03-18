-- CreateSorC creates a new ship or colony.
--
-- name: CreateSorC :one
INSERT INTO sorcs (empire_id, sorc_cd, tech_level, name, orbit_id, is_on_surface, rations, sol, birth_rate, death_rate)
VALUES (:empire_id, :sorc_cd, :tech_level, :name, :orbit_id, :is_on_surface, :rations, :sol, :birth_rate, :death_rate)
RETURNING id;

-- CreateSorCInventory creates a new ship or colony inventory entry.
--
-- name: CreateSorCInventory :exec
INSERT INTO inventory (sorc_id, unit_cd, tech_level, qty, mass, volume, is_assembled, is_stored)
VALUES (:sorc_id, :unit_cd, :tech_level, :qty, :mass, :volume, :is_assembled, :is_stored);

-- CreateSorCPopulation creates a new ship or colony population entry.
--
-- name: CreateSorCPopulation :exec
INSERT INTO population (sorc_id, population_cd, qty, pay_rate, rebel_qty)
VALUES (:sorc_id, :population_cd, :qty, :pay_rate, :rebel_qty);

-- CreateSorCFarmGroup creates a new ship or colony farm group.
--
-- name: CreateSorCFarmGroup :one
INSERT INTO farm_groups (sorc_id, group_no)
VALUES (:sorc_id, :group_no)
RETURNING id;

-- CreateSorCFarmGroupUnit creates a new ship or colony farm group unit.
--
-- name: CreateSorCFarmGroupUnit :one
INSERT INTO farm_group (farm_group_id, unit_cd, tech_level, nbr_of_units)
VALUES (:farm_group_id, :unit_cd, :tech_level, :nbr_of_units)
RETURNING id;

-- CreateSorCFactoryGroup creates a new ship or colony factory group.
--
-- name: CreateSorCFactoryGroup :one
INSERT INTO factory_groups (sorc_id, group_no, orders_cd, orders_tech_level, retool_turn_no)
VALUES (:sorc_id, :group_no, :orders_cd, :orders_tech_level, :retool_turn_no)
RETURNING id;

-- CreateSorCFactoryGroupUnit creates a new ship or colony factory group unit.
--
-- name: CreateSorCFactoryGroupUnit :one
INSERT INTO factory_group (factory_group_id, unit_cd, tech_level, nbr_of_units, orders_cd, orders_tech_level,
                           wip_25pct_qty, wip_50pct_qty, wip_75pct_qty)
VALUES (:factory_group_id, :unit_cd, :tech_level, :nbr_of_units, :orders_cd, :orders_tech_level, :wip_25pct_qty,
        :wip_50pct_qty, :wip_75pct_qty)
RETURNING id;

-- CreateSorCMiningGroup creates a new ship or colony mining group.
--
-- name: CreateSorCMiningGroup :one
INSERT INTO mining_groups (sorc_id, group_no, deposit_id)
VALUES (:sorc_id, :group_no, :deposit_id)
RETURNING id;

-- CreateSorCMiningGroupUnit creates a new ship or colony mining group unit.
--
-- name: CreateSorCMiningGroupUnit :one
INSERT INTO mining_group (mining_group_id, unit_cd, tech_level, nbr_of_units)
VALUES (:mining_group_id, :unit_cd, :tech_level, :nbr_of_units)
RETURNING id;


-- ReadAllColoniesByEmpire reads the colonies for a given empire in a game.
--
-- name: ReadAllColoniesByEmpire :many
SELECT sorcs.id        AS sorcs_id,
       systems.x,
       systems.y,
       systems.z,
       stars.sequence  AS suffix,
       orbits.orbit_no,
       sorc_codes.name AS sorc_kind,
       sorcs.tech_level,
       sorcs.name,
       sorcs.rations,
       sorcs.birth_rate,
       sorcs.death_rate,
       sorcs.sol
FROM sorcs,
     orbits,
     stars,
     systems,
     sorc_codes
WHERE sorcs.empire_id = :empire_id
  AND sorcs.sorc_cd IN ('COPN', 'CENC', 'CORB')
  AND orbits.id = sorcs.orbit_id
  AND stars.id = orbits.star_id
  AND systems.id = stars.system_id
  AND sorc_codes.code = sorcs.sorc_cd
ORDER BY sorcs.id;

-- ReadSorCPopulation reads the population for a given ship or colony.
--
-- name: ReadSorCPopulation :many
SELECT population.population_cd,
       population_codes.name AS population_kind,
       population.qty,
       population.pay_rate,
       population.rebel_qty
FROM population,
     population_codes
WHERE population.sorc_id = :sorc_id
  AND population_codes.code = population.population_cd
ORDER BY population_codes.sort_order;

-- ReadSorCInventory reads the inventory for a given ship or colony.
--
-- name: ReadSorCInventory :many
SELECT inventory.unit_cd,
       inventory.tech_level,
       unit_codes.name AS unit_kind,
       inventory.qty,
       inventory.mass,
       inventory.volume,
       inventory.is_assembled,
       inventory.is_stored
FROM inventory,
     unit_codes
WHERE inventory.sorc_id = :sorc_id
  AND unit_codes.code = inventory.unit_cd
ORDER BY inventory.unit_cd, inventory.tech_level, inventory.qty;

-- ReadSorCFactoryGroups reads the factory groups for a given ship or colony.
--
-- name: ReadSorCFactoryGroups :many
SELECT id AS group_id,
       group_no,
       orders_cd,
       orders_tech_level,
       retool_turn_no
FROM factory_groups
WHERE sorc_id = :sorc_id
ORDER BY group_no;

-- ReadSorCFactoryGroup reads the factory group for a given ship or colony.
--
-- name: ReadSorCFactoryGroup :many
SELECT orders_cd,
       orders_tech_level,
       tech_level AS factory_tech_level,
       nbr_of_units,
       wip_25pct_qty,
       wip_50pct_qty,
       wip_75pct_qty
FROM factory_group
WHERE factory_group_id = :factory_group_id
ORDER BY orders_cd, tech_level;

-- ReadSorCFarmGroups reads the farm groups for a given ship or colony.
--
-- name: ReadSorCFarmGroups :many
SELECT farm_groups.group_no,
       farm_group.tech_level,
       sum(farm_group.nbr_of_units) AS nbr_of_units
FROM farm_groups,
     farm_group
WHERE farm_groups.sorc_id = :sorc_id
  AND farm_groups.id = farm_group.farm_group_id
GROUP BY farm_groups.group_no, farm_group.tech_level
ORDER BY farm_groups.group_no, farm_group.tech_level;

-- ReadSorCMiningGroups reads the mining groups for a given ship or colony.
--
-- name: ReadSorCMiningGroups :many
SELECT mining_groups.id AS group_id,
       mining_groups.group_no,
       deposits.deposit_no,
       deposits.remaining_qty,
       deposits.kind    AS deposit_kind,
       deposits.yield_pct
FROM mining_groups,
     deposits
WHERE mining_groups.sorc_id = :sorc_id
  AND deposits.id = mining_groups.deposit_id
ORDER BY mining_groups.group_no, deposits.deposit_no;

-- ReadSorCMiningGroup returns the data for a given mining group.
--
-- name: ReadSorCMiningGroup :many
SELECT mining_group.tech_level,
       sum(mining_group.nbr_of_units) AS nbr_of_units
FROM mining_group
WHERE mining_group_id = :group_id
GROUP BY tech_level
ORDER BY tech_level;
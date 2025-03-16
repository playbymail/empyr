-- CreateSorC creates a new ship or colony.
--
-- name: CreateSorC :one
INSERT INTO sorcs (empire_id, sorc_cd, tech_level, name, orbit_id, is_on_surface, rations, sol, birth_rate, death_rate)
VALUES (:empire_id, :sorc_cd, :tech_level, :name, :orbit_id, :is_on_surface, :rations, :sol, :birth_rate, :death_rate)
RETURNING id;

-- CreateSorCInventory creates a new ship or colony inventory entry.
--
-- name: CreateSorCInventory :exec
INSERT INTO inventory (sorc_id, unit_cd, qty)
VALUES (:sorc_id, :unit_cd, :qty);

-- CreateSorCPopulation creates a new ship or colony population entry.
--
-- name: CreateSorCPopulation :exec
INSERT INTO population (sorc_id, population_cd, qty, pay_rate, rebel_qty)
VALUES (:sorc_id, :population_cd, :qty, :pay_rate, :rebel_qty);


-- ReadAllColoniesByEmpire reads the colonies for a given empire in a game.
--
-- name: ReadAllColoniesByEmpire :many
SELECT sorcs.id        AS sorcs_id,
       systems.x,
       systems.y,
       systems.z,
       stars.sequence  as suffix,
       orbits.orbit_no,
       sorc_codes.name AS sorc_kind,
       sorcs.tech_level,
       sorcs.name,
       sorcs.rations,
       sorcs.birth_rate,
       sorcs.death_rate,
       sorcs.sol
FROM sorcs,
     sorc_codes,
     orbits,
     stars,
     systems
WHERE sorcs.empire_id = :empire_id
  AND sorcs.sorc_cd IN ('COPN', 'CENC', 'CORB')
  AND sorcs.sorc_cd = sorc_codes.code
  AND orbits.star_id = stars.id
  AND systems.id = stars.system_id
ORDER BY sorcs.id;
--  Copyright (c) 2025 Michael D Henderson. All rights reserved.
--

-- CreateGame creates a new game.
--
-- name: CreateGame :one
INSERT INTO games (code, name, display_name)
VALUES (:code, :name, :display_name)
RETURNING id;

-- ReadGameByCode gets a game by its code.
--
-- name: ReadGameByCode :one
SELECT id,
       code,
       name,
       display_name,
       current_turn,
       last_empire_no,
       home_system_id,
       home_star_id,
       home_orbit_id,
       home_planet_id
FROM games
WHERE code = :code;

-- DeleteGame deletes an existing game
--
-- name: DeleteGame :exec
DELETE
FROM games
WHERE code = :code;

-- UpdateGameTurn increments the game turn number.
--
-- name: UpdateGameTurn :exec
UPDATE games
SET current_turn = :turn_number
WHERE id = :game_id;

-- GetCurrentGameTurn gets the current game turn.
--
-- name: GetCurrentGameTurn :one
SELECT current_turn
FROM games
WHERE id = :game_id;

-- CreateSystem creates a new system.
--
-- name: CreateSystem :one
INSERT INTO systems (game_id, x, y, z, scarcity)
VALUES (:game_id, :x, :y, :z, :scarcity)
RETURNING id;

-- CreateStar creates a new star in an existing system.
--
-- name: CreateStar :one
INSERT INTO stars (system_id, sequence, scarcity)
VALUES (:system_id, :sequence, :scarcity)
RETURNING id;

-- CreateSystemDistance inserts the distance between two systems.
--
-- name: CreateSystemDistance :exec
INSERT INTO system_distances (from_system_id, to_system_id, distance)
VALUES (:from_system_id, :to_system_id, :distance);

-- CreateOrbit creates a new orbit.
--
-- name: CreateOrbit :one
INSERT INTO orbits (star_id, orbit_no, kind, scarcity)
VALUES (:star_id, :orbit_no, :kind, :scarcity)
RETURNING id;

-- DeleteEmptyOrbits deletes all orbits with no planets.
--
-- name: DeleteEmptyOrbits :exec
DELETE
FROM orbits
WHERE kind = 0;

-- CreatePlanet creates a new planet.
--
-- name: CreatePlanet :one
INSERT INTO planets (orbit_id, kind, scarcity, habitability)
VALUES (:orbit_id, :kind, :scarcity, :habitability)
RETURNING id;

-- CreateDeposit creates a new deposit.
--
-- name: CreateDeposit :one
INSERT INTO deposits (planet_id, deposit_no, kind, initial_qty, remaining_qty, yield_pct)
VALUES (:planet_id, :deposit_no, :kind, :initial_qty, :remaining_qty, :yield_pct)
RETURNING id;

-- DeleteEmptyOrbits deletes all orbits with no planets.
--
-- name: DeleteEmptyDeposits :exec
DELETE
FROM deposits
WHERE kind = 0;

-- UpdateGameEmpireMetadata updates the empire metadata in the games table.
--
-- name: UpdateGameEmpireMetadata :exec
UPDATE games
SET last_empire_no = :empire_no,
    home_system_id = :home_system_id,
    home_star_id   = :home_star_id,
    home_orbit_id  = :home_orbit_id,
    home_planet_id = :home_planet_id
WHERE id = :game_id;

-- UpdateGameEmpireCounter updates the last empire number in the games table.
--
-- name: UpdateGameEmpireCounter :exec
UPDATE games
SET last_empire_no = :empire_no
WHERE id = :game_id;

-- ReadNextEmpireNumber reads the next empire number.
--
-- name: ReadNextEmpireNumber :one
UPDATE games
SET last_empire_no = last_empire_no + 1
WHERE id = :game_id
RETURNING last_empire_no as next_empire_no;

-- CreateEmpire creates a new empire.
--
-- name: CreateEmpire :one
INSERT INTO empires (game_id, empire_no, name, handle, home_system_id, home_star_id, home_orbit_id, home_planet_id)
VALUES (:game_id, :empire_no, :name, :handle, :home_system_id, :home_star_id, :home_orbit_id, :home_planet_id)
RETURNING id;

-- ReadClusterMap reads the cluster map.
--
-- name: ReadClusterMap :many
SELECT systems.id      AS id,
       systems.x       as x,
       systems.y       as y,
       systems.z       as z,
       count(stars.id) AS number_of_stars
FROM games
         LEFT JOIN systems on games.id = systems.game_id
         LEFT JOIN stars on systems.id = stars.system_id
WHERE games.code = :game_code
GROUP BY systems.id, systems.x, systems.y, systems.z
ORDER BY systems.id;

-- ReadGameEmpire returns the data for a single empire in a game.
--
-- name: ReadGameEmpire :one
SELECT game_id, id AS empire_id, empire_no
FROM empires
WHERE empire_no = :empire_no
  AND game_id = (SELECT id FROM games WHERE code = :game_code);

-- ReadGameEmpires returns the data for all empires in a game.
--
-- name: ReadGameEmpires :many
SELECT game_id, id AS empire_id, empire_no
FROM empires
WHERE game_id = (SELECT id FROM games WHERE code = :game_code);

-- CreateSorC creates a new ship or colony.
--
-- name: CreateSorC :one
INSERT INTO sorcs (empire_id, kind)
VALUES (:empire_id, :kind)
RETURNING id;

-- CreateSorCDetails creates a new ship or colony details entry.
--
-- name: CreateSorCDetails :one
INSERT INTO sorc_details (sorc_id, turn_no, tech_level, name, uem_qty, uem_pay, usk_qty, usk_pay, pro_qty, pro_pay,
                          sld_qty, sld_pay, cnw_qty, spy_qty, rations, birth_rate, death_rate, sol, orbit_id,
                          is_on_surface)
VALUES (:sorc_id, :turn_no, :tech_level, :name, :uem_qty, :uem_pay, :usk_qty, :usk_pay, :pro_qty, :pro_pay, :sld_qty,
        :sld_pay, :cnw_qty, :spy_qty, :rations, :birth_rate, :death_rate, :sol, :orbit_id, :is_on_surface)
RETURNING id;

-- CreateSorCInfrastructure creates a new ship or colony infrastructure entry.
--
-- name: CreateSorCInfrastructure :exec
INSERT INTO sorc_infrastructure (sorc_detail_id, kind, tech_level, qty)
VALUES (:sorc_id, :kind, :tech_level, :qty);

-- CreateSorCInventory creates a new ship or colony inventory entry.
--
-- name: CreateSorCInventory :exec
INSERT INTO sorc_inventory (sorc_detail_id, kind, tech_level, qty_assembled, qty_stored)
VALUES (:sorc_id, :kind, :tech_level, :qty_assembled, :qty_stored);

-- CreateSorCSuperstructure creates a new ship or colony infrastructure entry.
--
-- name: CreateSorCSuperstructure :exec
INSERT INTO sorc_superstructure (sorc_detail_id, kind, tech_level, qty)
VALUES (:sorc_id, :kind, :tech_level, :qty);

-- ReadEmpiresInGame reads the empires in a game.
--
-- name: ReadEmpiresInGame :many
SELECT empire_no, id
FROM empires
WHERE game_id = :game_id
ORDER by empire_no;

-- ReadEmpireAllColoniesForTurn reads the colonies for a given empire and turn in a game.
--
-- name: ReadEmpireAllColoniesForTurn :many
SELECT sorc_id,
       sorcs.kind,
       tech_level,
       name,
       uem_qty,
       uem_pay,
       usk_qty,
       usk_pay,
       pro_qty,
       pro_pay,
       sld_qty,
       sld_pay,
       cnw_qty,
       spy_qty,
       rations,
       birth_rate,
       death_rate,
       sol,
       systems.x,
       systems.y,
       systems.z,
       stars.sequence as suffix,
       orbits.orbit_no,
       is_on_surface
FROM sorcs
         LEFT JOIN sorc_details ON sorcs.id = sorc_details.sorc_id AND sorc_details.turn_no = :turn_no
         LEFT JOIN orbits ON orbits.id = sorc_details.orbit_id
         LEFT JOIN stars ON stars.id = orbits.star_id
         LEFT JOIN systems ON systems.id = stars.system_id
WHERE sorcs.empire_id = :empire_id
  AND sorcs.kind in (2, 3, 4)
ORDER BY sorcs.id, sorcs.kind;
--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- foreign keys must be disabled to drop tables
PRAGMA foreign_keys = OFF;
DROP TABLE IF EXISTS meta_migrations;
DROP TABLE IF EXISTS clusters;
DROP TABLE IF EXISTS deposits;
DROP TABLE IF EXISTS empires;
DROP TABLE IF EXISTS factory_group;
DROP TABLE IF EXISTS factory_groups;
DROP TABLE IF EXISTS farm_group;
DROP TABLE IF EXISTS farm_groups;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS mining_group;
DROP TABLE IF EXISTS mining_groups;
DROP TABLE IF EXISTS orbit_codes;
DROP TABLE IF EXISTS orbits;
DROP TABLE IF EXISTS planet_codes;
DROP TABLE IF EXISTS planets;
DROP TABLE IF EXISTS population;
DROP TABLE IF EXISTS population_codes;
DROP TABLE IF EXISTS report_probes;
DROP TABLE IF EXISTS report_spies;
DROP TABLE IF EXISTS report_surveys;
DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS sorc_codes;
DROP TABLE IF EXISTS sorcs;
DROP TABLE IF EXISTS stars;
DROP TABLE IF EXISTS system_distances;
DROP TABLE IF EXISTS systems;
DROP TABLE IF EXISTS unit_codes;
DROP TABLE IF EXISTS users;
-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

-- create the table for managing migrations
CREATE TABLE meta_migrations
(
    version    INTEGER  NOT NULL UNIQUE,
    comment    TEXT     NOT NULL,
    script     TEXT     NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE games
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    code           TEXT     NOT NULL UNIQUE,
    name           TEXT     NOT NULL UNIQUE,
    display_name   TEXT     NOT NULL UNIQUE,
    current_turn   INTEGER  NOT NULL DEFAULT 0,
    last_empire_no INTEGER  NOT NULL DEFAULT 0,
    is_active      INTEGER  NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE clusters
(
    game_id        INTEGER NOT NULL,
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    home_system_id INTEGER NOT NULL DEFAULT 0,
    home_star_id   INTEGER NOT NULL DEFAULT 0,
    home_orbit_id  INTEGER NOT NULL DEFAULT 0,
    home_planet_id INTEGER NOT NULL DEFAULT 0,
    CONSTRAINT fk_game_id
        FOREIGN KEY (game_id)
            REFERENCES games (id)
            ON DELETE CASCADE
);

CREATE TABLE systems
(
    cluster_id INTEGER NOT NULL,
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    x          INTEGER NOT NULL CHECK (x BETWEEN 0 AND 30),
    y          INTEGER NOT NULL CHECK (y BETWEEN 0 AND 30),
    z          INTEGER NOT NULL CHECK (z BETWEEN 0 AND 30),
    UNIQUE (cluster_id, x, y, z),
    CONSTRAINT fk_cluster_id
        FOREIGN KEY (cluster_id)
            REFERENCES clusters (id)
            ON DELETE CASCADE
);

CREATE TABLE stars
(
    system_id INTEGER NOT NULL,
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    sequence  TEXT    NOT NULL CHECK (LENGTH(sequence) = 1 AND sequence BETWEEN 'A' AND 'Z'),
    UNIQUE (system_id, sequence),
    CONSTRAINT fk_system_id
        FOREIGN KEY (system_id)
            REFERENCES systems (id)
            ON DELETE CASCADE
);

CREATE TABLE orbit_codes
(
    code TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
INSERT INTO orbit_codes (code, name)
VALUES ('EMPTY', 'empty');
INSERT INTO orbit_codes (code, name)
VALUES ('ASTR', 'asteroid belt');
INSERT INTO orbit_codes (code, name)
VALUES ('ERTH', 'earth-like planet');
INSERT INTO orbit_codes (code, name)
VALUES ('GASG', 'gas giant');
INSERT INTO orbit_codes (code, name)
VALUES ('ICEG', 'ice giant');
INSERT INTO orbit_codes (code, name)
VALUES ('RCKY', 'rocky planet');

CREATE TABLE orbits
(
    star_id  INTEGER NOT NULL,
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    orbit_no INTEGER NOT NULL CHECK (orbit_no BETWEEN 1 AND 10),
    kind     TEXT    NOT NULL,
    UNIQUE (star_id, orbit_no),
    CONSTRAINT fk_star_id
        FOREIGN KEY (star_id)
            REFERENCES stars (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_orbit_kind_cd
        FOREIGN KEY (kind)
            REFERENCES orbit_codes (code)
            ON DELETE CASCADE
);

CREATE TABLE planet_codes
(
    code TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
INSERT INTO planet_codes (code, name)
VALUES ('NONE', 'no planet');
INSERT INTO planet_codes (code, name)
VALUES ('ASTR', 'asteroid belt');
INSERT INTO planet_codes (code, name)
VALUES ('GASG', 'gas giant');
INSERT INTO planet_codes (code, name)
VALUES ('TERR', 'terrestrial planet');

CREATE TABLE planets
(
    orbit_id     INTEGER NOT NULL,
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    kind         TEXT    NOT NULL,
    habitability INTEGER NOT NULL CHECK (habitability BETWEEN 0 AND 25),
    CONSTRAINT fk_orbit_id
        FOREIGN KEY (orbit_id)
            REFERENCES orbits (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_planet_kind_cd
        FOREIGN KEY (kind)
            REFERENCES planet_codes (code)
            ON DELETE CASCADE
);

CREATE TABLE unit_codes
(
    code           TEXT    NOT NULL PRIMARY KEY,
    name           TEXT    NOT NULL UNIQUE,
    category       TEXT    NOT NULL,
    is_operational INTEGER NOT NULL CHECK (is_operational IN (0, 1)),
    is_consumable  INTEGER NOT NULL CHECK (is_consumable IN (0, 1)),
    is_resource    INTEGER NOT NULL CHECK (is_resource IN (0, 1)),
    aliases        TEXT
);

INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('ANM', 'Anti-Missiles', 'Vehicles', 0, 1, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('ASC', 'Assault Craft', 'Vehicles', 0, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('ASW', 'Assault Weapons', 'Vehicles', 0, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('AUT', 'Automation', 'Assembly', 1, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('CNGD', 'Consumer Goods', 'Consumables', 0, 1, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('ESH', 'Energy Shields', 'Assembly', 1, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('EWP', 'Energy Weapons', 'Assembly', 1, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('FCT', 'Factories', 'Assembly', 1, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('FOOD', 'Food', 'Consumables', 0, 1, 1);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('FRM', 'Farms', 'Assembly', 1, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('FUEL', 'Fuel', 'Consumables', 0, 1, 1);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('GOLD', 'Gold', 'Consumables', 0, 1, 1);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('HEN', 'Hyper Engines', 'Assembly', 1, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('LAB', 'Laboratories', 'Assembly', 1, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('LFS', 'Life Supports', 'Assembly', 1, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('METS', 'Metals', 'Consumables', 0, 1, 1);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('MIN', 'Mines', 'Assembly', 1, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('MSL', 'Missile Launchers', 'Assembly', 1, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('MSS', 'Missiles', 'Vehicles', 0, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('MTBT', 'Military Robots', 'Bots', 0, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('MTSP', 'Military Supplies', 'Consumables', 0, 1, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('NMTS', 'Non-Metals', 'Consumables', 0, 1, 1);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('PWP', 'Power Plants', 'Assembly', 1, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('RPV', 'Robot Probe Vehicles', 'Bots', 0, 1, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('RSCH', 'Research', 'Consumables', 0, 1, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('SEN', 'Sensors', 'Assembly', 1, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('SLS', 'Light Structure', 'Assembly', 1, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('SPD', 'Space Drives', 'Assembly', 1, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('STU', 'Structure', 'Assembly', 1, 0, 0);
INSERT INTO unit_codes (code, name, category, is_operational, is_consumable, is_resource)
VALUES ('TPT', 'Transports', 'Vehicles', 0, 0, 0);

UPDATE unit_codes
SET aliases = 'ASWP'
WHERE code = 'ASW';
UPDATE unit_codes
SET aliases = 'FACT, FCTU, FU'
WHERE code = 'FCT';
UPDATE unit_codes
SET aliases = 'FARM, FRMU'
WHERE code = 'FRM';
UPDATE unit_codes
SET aliases = 'MSS'
WHERE code = 'MTSP';
UPDATE unit_codes
SET aliases = 'MINU, MU'
WHERE code = 'MIN';
UPDATE unit_codes
SET aliases = 'STUN'
WHERE code = 'STU';

CREATE TABLE deposits
(
    planet_id     INTEGER NOT NULL,
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    deposit_no    INTEGER NOT NULL CHECK (deposit_no BETWEEN 1 AND 35),
    kind          TEXT    NOT NULL,
    remaining_qty INTEGER NOT NULL CHECK (remaining_qty BETWEEN 0 AND 99000000),
    yield_pct     INTEGER NOT NULL CHECK (yield_pct BETWEEN 0 AND 100),
    UNIQUE (planet_id, deposit_no),
    CONSTRAINT fk_planet_id
        FOREIGN KEY (planet_id)
            REFERENCES planets (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_deposit_kind_cd
        FOREIGN KEY (kind)
            REFERENCES unit_codes (code)
            ON DELETE CASCADE
);

CREATE TABLE users
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    username        TEXT     NOT NULL UNIQUE,
    email           TEXT     NOT NULL UNIQUE,
    hashed_password TEXT     NOT NULL UNIQUE,
    is_active       INTEGER  NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    is_admin        INTEGER  NOT NULL DEFAULT 0 CHECK (is_admin IN (0, 1)),
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE empires
(
    game_id        INTEGER NOT NULL,
    user_id        INTEGER NOT NULL,
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    is_active      INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    empire_no      INTEGER NOT NULL CHECK (empire_no BETWEEN 1 AND 256),
    name           TEXT    NOT NULL,
    home_system_id INTEGER NOT NULL,
    home_star_id   INTEGER NOT NULL,
    home_orbit_id  INTEGER NOT NULL,
    home_planet_id INTEGER NOT NULL,
    UNIQUE (game_id, user_id),
    UNIQUE (game_id, empire_no),
    CONSTRAINT fk_game_id
        FOREIGN KEY (game_id)
            REFERENCES games (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_system_id
        FOREIGN KEY (home_system_id)
            REFERENCES systems (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_star_id
        FOREIGN KEY (home_star_id)
            REFERENCES stars (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_orbit_id
        FOREIGN KEY (home_orbit_id)
            REFERENCES orbits (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_home_planet_id
        FOREIGN KEY (home_planet_id)
            REFERENCES planets (id)
            ON DELETE CASCADE
);

CREATE TABLE sorc_codes
(
    code TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
INSERT INTO sorc_codes (code, name)
VALUES ('NONE', 'none');
INSERT INTO sorc_codes (code, name)
VALUES ('SHIP', 'Ship');
INSERT INTO sorc_codes (code, name)
VALUES ('COPN', 'Open Surface Colony');
INSERT INTO sorc_codes (code, name)
VALUES ('CENC', 'Enclosed Surface Colony');
INSERT INTO sorc_codes (code, name)
VALUES ('CORB', 'Orbital Colony');

CREATE TABLE sorcs
(
    empire_id     INTEGER NOT NULL,
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    sorc_cd       TEXT    NOT NULL,
    tech_level    INTEGER NOT NULL CHECK (tech_level BETWEEN 0 AND 10),
    name          TEXT    NOT NULL,
    orbit_id      INTEGER NOT NULL,
    is_on_surface INTEGER NOT NULL DEFAULT 0 CHECK (is_on_surface IN (0, 1)),
    rations       REAL    NOT NULL CHECK (rations >= 0),
    sol           REAL    NOT NULL CHECK (sol >= 0),
    birth_rate    REAL    NOT NULL CHECK (birth_rate >= 0),
    death_rate    REAL    NOT NULL CHECK (death_rate >= 0),
    CONSTRAINT fk_empire_id
        FOREIGN KEY (empire_id)
            REFERENCES empires (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_sorc_cd
        FOREIGN KEY (sorc_cd)
            REFERENCES sorc_codes (code)
            ON DELETE CASCADE,
    CONSTRAINT fk_orbit_id
        FOREIGN KEY (orbit_id)
            REFERENCES orbits (id)
            ON DELETE CASCADE
);

CREATE TABLE inventory
(
    sorc_id      INTEGER NOT NULL,
    unit_cd      TEXT    NOT NULL,
    tech_level   INTEGER NOT NULL CHECK (tech_level BETWEEN 0 AND 10),
    qty          INTEGER NOT NULL CHECK (qty >= 0),
    mass         REAL    NOT NULL CHECK (mass >= 0),
    volume       REAL    NOT NULL CHECK (volume >= 0),
    is_assembled INTEGER NOT NULL DEFAULT 0 CHECK (is_assembled IN (0, 1)),
    is_stored    INTEGER NOT NULL DEFAULT 0 CHECK (is_stored IN (0, 1)),
    CONSTRAINT fk_sorc_id
        FOREIGN KEY (sorc_id)
            REFERENCES sorcs (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_unit_cd
        FOREIGN KEY (unit_cd)
            REFERENCES unit_codes (code)
            ON DELETE CASCADE
);

CREATE TABLE population_codes
(
    code          TEXT    NOT NULL PRIMARY KEY,
    name          TEXT    NOT NULL UNIQUE,
    base_pay_rate REAL    NOT NULL CHECK (base_pay_rate >= 0),
    sort_order    INTEGER NOT NULL
);
INSERT INTO population_codes (code, name, base_pay_rate, sort_order)
VALUES ('CNW', 'Construction Worker', 0.5000, 5);
INSERT INTO population_codes (code, name, base_pay_rate, sort_order)
VALUES ('PLC', 'Police', 0.2500, 7);
INSERT INTO population_codes (code, name, base_pay_rate, sort_order)
VALUES ('PRO', 'Professional', 0.3750, 3);
INSERT INTO population_codes (code, name, base_pay_rate, sort_order)
VALUES ('SAG', 'Special Agents', 0.6250, 8);
INSERT INTO population_codes (code, name, base_pay_rate, sort_order)
VALUES ('SLD', 'Soldier', 0.2500, 4);
INSERT INTO population_codes (code, name, base_pay_rate, sort_order)
VALUES ('SPY', 'Spy', 0.6250, 6);
INSERT INTO population_codes (code, name, base_pay_rate, sort_order)
VALUES ('TRN', 'Trainees', 0.1250, 9);
INSERT INTO population_codes (code, name, base_pay_rate, sort_order)
VALUES ('UEM', 'Unemployable', 0.0000, 1);
INSERT INTO population_codes (code, name, base_pay_rate, sort_order)
VALUES ('USK', 'Unskilled', 0.1250, 2);

CREATE TABLE population
(
    sorc_id       INTEGER NOT NULL,
    population_cd TEXT    NOT NULL,
    qty           INTEGER NOT NULL CHECK (qty >= 0),
    pay_rate      REAL    NOT NULL CHECK (pay_rate >= 0),
    rebel_qty     INTEGER NOT NULL CHECK (rebel_qty >= 0),
    CONSTRAINT fk_sorc_id
        FOREIGN KEY (sorc_id)
            REFERENCES sorcs (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_population_cd
        FOREIGN KEY (population_cd)
            REFERENCES population_codes (code)
            ON DELETE CASCADE
);

CREATE TABLE farm_groups
(
    sorc_id  INTEGER NOT NULL,
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    group_no INTEGER NOT NULL CHECK (group_no BETWEEN 1 AND 30),
    UNIQUE (sorc_id, group_no),
    CONSTRAINT fk_sorc_id
        FOREIGN KEY (sorc_id)
            REFERENCES sorcs (id)
            ON DELETE CASCADE
);

CREATE TABLE farm_group
(
    farm_group_id INTEGER NOT NULL,
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    unit_cd       TEXT    NOT NULL,
    tech_level    INTEGER NOT NULL CHECK (tech_level BETWEEN 0 AND 10),
    nbr_of_units  INTEGER NOT NULL CHECK (nbr_of_units >= 0),
    CONSTRAINT fk_farm_group_id
        FOREIGN KEY (farm_group_id)
            REFERENCES farm_groups (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_unit_cd
        FOREIGN KEY (unit_cd)
            REFERENCES unit_codes (code)
            ON DELETE CASCADE
);

CREATE TABLE factory_groups
(
    sorc_id           INTEGER NOT NULL,
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    group_no          INTEGER NOT NULL CHECK (group_no BETWEEN 1 AND 30),
    orders_cd         TEXT    NOT NULL,
    orders_tech_level INTEGER NOT NULL CHECK (orders_tech_level BETWEEN 0 AND 10),
    retool_turn_no    INTEGER,
    UNIQUE (sorc_id, group_no),
    CONSTRAINT fk_sorc_id
        FOREIGN KEY (sorc_id)
            REFERENCES sorcs (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_orders_cd
        FOREIGN KEY (orders_cd)
            REFERENCES unit_codes (code)
            ON DELETE CASCADE
);

CREATE TABLE factory_group
(
    factory_group_id  INTEGER NOT NULL,
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    unit_cd           TEXT    NOT NULL,
    tech_level        INTEGER NOT NULL CHECK (tech_level BETWEEN 0 AND 10),
    nbr_of_units      INTEGER NOT NULL CHECK (nbr_of_units >= 0),
    orders_cd         TEXT    NOT NULL,
    orders_tech_level INTEGER NOT NULL CHECK (orders_tech_level BETWEEN 0 AND 10),
    wip_25pct_qty     INTEGER NOT NULL CHECK (wip_25pct_qty >= 0),
    wip_50pct_qty     INTEGER NOT NULL CHECK (wip_50pct_qty >= 0),
    wip_75pct_qty     INTEGER NOT NULL CHECK (wip_75pct_qty >= 0),
    CONSTRAINT fk_factory_group_id
        FOREIGN KEY (factory_group_id)
            REFERENCES factory_groups (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_unit_cd
        FOREIGN KEY (unit_cd)
            REFERENCES unit_codes (code)
            ON DELETE CASCADE
);

CREATE TABLE mining_groups
(
    sorc_id    INTEGER NOT NULL,
    deposit_id INTEGER NOT NULL,
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    group_no   INTEGER NOT NULL CHECK (group_no BETWEEN 1 AND 35),
    UNIQUE (sorc_id, group_no),
    CONSTRAINT fk_sorc_id
        FOREIGN KEY (sorc_id)
            REFERENCES sorcs (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_deposit_id
        FOREIGN KEY (deposit_id)
            REFERENCES deposits (id)
            ON DELETE CASCADE
);

CREATE TABLE mining_group
(
    mining_group_id INTEGER NOT NULL,
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    unit_cd         TEXT    NOT NULL,
    tech_level      INTEGER NOT NULL CHECK (tech_level BETWEEN 0 AND 10),
    nbr_of_units    INTEGER NOT NULL CHECK (nbr_of_units >= 0),
    CONSTRAINT fk_mining_group_id
        FOREIGN KEY (mining_group_id)
            REFERENCES mining_groups (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_unit_cd
        FOREIGN KEY (unit_cd)
            REFERENCES unit_codes (code)
            ON DELETE CASCADE
);


CREATE TABLE probe_orders
(
    sorc_id    INTEGER NOT NULL,
    turn_no    INTEGER NOT NULL CHECK (turn_no >= 0),
    tech_level INTEGER NOT NULL CHECK (tech_level BETWEEN 0 AND 10),
    kind       TEXT    NOT NULL CHECK (kind in ('system', 'star', 'orbit', 'sorc')),
    target_id  INTEGER NOT NULL,
    status     TEXT,
    CONSTRAINT fk_sorc_id
        FOREIGN KEY (sorc_id)
            REFERENCES sorcs (id)
            ON DELETE CASCADE
);

CREATE TABLE survey_orders
(
    sorc_id    INTEGER NOT NULL,
    turn_no    INTEGER NOT NULL CHECK (turn_no >= 0),
    tech_level INTEGER NOT NULL CHECK (tech_level BETWEEN 0 AND 10),
    orbit_id   INTEGER NOT NULL,
    status     TEXT,
    CONSTRAINT fk_sorc_id
        FOREIGN KEY (sorc_id)
            REFERENCES sorcs (id)
            ON DELETE CASCADE
);

CREATE TABLE reports
(
    sorc_id INTEGER NOT NULL,
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    turn_no INTEGER NOT NULL CHECK (turn_no >= 0),
    UNIQUE (sorc_id, turn_no),
    CONSTRAINT fk_sorc_id
        FOREIGN KEY (sorc_id)
            REFERENCES sorcs (id)
            ON DELETE CASCADE
);

CREATE TABLE report_production_inputs
(
    report_id  INTEGER NOT NULL,
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    category   TEXT    NOT NULL,
    fuel       INTEGER NOT NULL CHECK (fuel >= 0),
    gold       INTEGER NOT NULL CHECK (gold >= 0),
    metals     INTEGER NOT NULL CHECK (metals >= 0),
    non_metals INTEGER NOT NULL CHECK (non_metals >= 0),
    CONSTRAINT fk_report_id
        FOREIGN KEY (report_id)
            REFERENCES reports (id)
            ON DELETE CASCADE
);

CREATE TABLE report_production_outputs
(
    report_id    INTEGER NOT NULL,
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    category     TEXT    NOT NULL,
    unit_cd      TEXT    NOT NULL,
    tech_level   INTEGER NOT NULL CHECK (tech_level BETWEEN 0 AND 10),
    farmed       INTEGER NOT NULL CHECK (farmed >= 0),
    mined        INTEGER NOT NULL CHECK (mined >= 0),
    manufactured INTEGER NOT NULL CHECK (manufactured >= 0),
    CONSTRAINT fk_report_id
        FOREIGN KEY (report_id)
            REFERENCES reports (id)
            ON DELETE CASCADE
);

CREATE TABLE report_surveys
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    report_id INTEGER NOT NULL,
    orbit_id  INTEGER NOT NULL,
    CONSTRAINT fk_report_id
        FOREIGN KEY (report_id)
            REFERENCES reports (id)
            ON DELETE CASCADE
);

CREATE TABLE report_survey_deposits
(
    report_id         INTEGER NOT NULL,
    deposit_no        INTEGER NOT NULL,
    deposit_qty       INTEGER NOT NULL,
    deposit_kind      TEXT    NOT NULL,
    deposit_yield_pct INTEGER NOT NULL,
    CONSTRAINT fk_report_id
        FOREIGN KEY (report_id)
            REFERENCES report_surveys (id)
            ON DELETE CASCADE
);

CREATE TABLE report_probes
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    report_id      INTEGER NOT NULL,
    orbit_id       INTEGER NOT NULL,
    habitability   INTEGER NOT NULL,
    fuel_qty       INTEGER NOT NULL,
    gold_qty       INTEGER NOT NULL,
    metals_qty     INTEGER NOT NULL,
    non_metals_qty INTEGER NOT NULL,
    CONSTRAINT fk_report_id
        FOREIGN KEY (report_id)
            REFERENCES reports (id)
            ON DELETE CASCADE
);

CREATE TABLE report_probe_sorcs
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    report_id INTEGER NOT NULL,
    empire_id INTEGER NOT NULL, -- fk to empire that controls the planet being reported on
    sorc_id   INTEGER NOT NULL, -- fk to sorc being reported on
    sorc_cd   TEXT    NOT NULL,
    sorc_mass INTEGER NOT NULL, -- estimated mass of sorc
    CONSTRAINT fk_report_id
        FOREIGN KEY (report_id)
            REFERENCES report_probes (id)
            ON DELETE CASCADE
);

CREATE TABLE report_spies
(
    report_id INTEGER NOT NULL,
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    CONSTRAINT fk_report_id
        FOREIGN KEY (report_id)
            REFERENCES reports (id)
            ON DELETE CASCADE
);

CREATE TABLE system_distances
(
    from_system_id INTEGER NOT NULL,
    to_system_id   INTEGER NOT NULL,
    distance       INTEGER NOT NULL,
    UNIQUE (from_system_id, to_system_id),
    CONSTRAINT fk_from_system_id
        FOREIGN KEY (from_system_id)
            REFERENCES systems (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_to_system_id
        FOREIGN KEY (to_system_id)
            REFERENCES systems (id)
            ON DELETE CASCADE
);

-- create table deposits_summary as
-- select deposits.planet_id,
--        case when kind = 'FUEL' then sum(remaining_qty) else 0 end as fuel,
--        case when kind = 'GOLD' then sum(remaining_qty) else 0 end as gold,
--        case when kind = 'METS' then sum(remaining_qty) else 0 end as mets,
--        case when kind = 'NMTS' then sum(remaining_qty) else 0 end as nmts
-- from deposits
-- group by planet_id;
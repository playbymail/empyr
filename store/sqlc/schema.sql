--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- foreign keys must be disabled to drop tables
PRAGMA foreign_keys = OFF;

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
    code           TEXT     NOT NULL,
    name           TEXT     NOT NULL,
    display_name   TEXT     NOT NULL,
    current_turn   INTEGER  NOT NULL DEFAULT 0,
    home_system_id INTEGER  NOT NULL DEFAULT 0,
    home_star_id   INTEGER  NOT NULL DEFAULT 0,
    home_orbit_id  INTEGER  NOT NULL DEFAULT 0,
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE systems
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    x            INTEGER NOT NULL CHECK (x BETWEEN 0 AND 30),
    y            INTEGER NOT NULL CHECK (y BETWEEN 0 AND 30),
    z            INTEGER NOT NULL CHECK (z BETWEEN 0 AND 30),
    system_name  TEXT    NOT NULL UNIQUE,
    nbr_of_stars INTEGER NOT NULL DEFAULT 0,
    UNIQUE (x, y, z)
);

CREATE TABLE stars
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    system_id     INTEGER NOT NULL,
    sequence      TEXT    NOT NULL CHECK (LENGTH(sequence) = 1 AND sequence BETWEEN 'A' AND 'Z'),
    star_name     TEXT    NOT NULL UNIQUE,
    nbr_of_orbits INTEGER NOT NULL DEFAULT 0,
    UNIQUE (system_id, sequence),
    CONSTRAINT fk_system_id
        FOREIGN KEY (system_id)
            REFERENCES systems (id)
);

CREATE TABLE orbit_codes
(
    code TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
INSERT INTO orbit_codes (code, name)
VALUES ('NONE', 'no planet');
INSERT INTO orbit_codes (code, name)
VALUES ('ASTR', 'asteroid belt');
INSERT INTO orbit_codes (code, name)
VALUES ('GASG', 'gas giant');
INSERT INTO orbit_codes (code, name)
VALUES ('TERR', 'terrestrial planet');

CREATE TABLE orbits
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    system_id       INTEGER NOT NULL,
    star_id         INTEGER NOT NULL,
    orbit_no        INTEGER NOT NULL CHECK (orbit_no BETWEEN 1 AND 10),
    kind            TEXT    NOT NULL,
    habitability    INTEGER NOT NULL CHECK (habitability BETWEEN 0 AND 25),
    nbr_of_deposits INTEGER NOT NULL DEFAULT 0,
    UNIQUE (star_id, orbit_no),
    CONSTRAINT fk_system_id
        FOREIGN KEY (system_id)
            REFERENCES systems (id),
    CONSTRAINT fk_star_id
        FOREIGN KEY (star_id)
            REFERENCES stars (id),
    CONSTRAINT fk_orbit_kind_cd
        FOREIGN KEY (kind)
            REFERENCES orbit_codes (code)
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
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    system_id  INTEGER NOT NULL,
    star_id    INTEGER NOT NULL,
    orbit_id   INTEGER NOT NULL,
    deposit_no INTEGER NOT NULL CHECK (deposit_no BETWEEN 1 AND 35),
    kind       TEXT    NOT NULL,
    qty        INTEGER NOT NULL CHECK (qty BETWEEN 0 AND 99000000),
    yield_pct  INTEGER NOT NULL CHECK (yield_pct BETWEEN 0 AND 100),
    UNIQUE (orbit_id, deposit_no),
    CONSTRAINT fk_system_id
        FOREIGN KEY (system_id)
            REFERENCES systems (id),
    CONSTRAINT fk_star_id
        FOREIGN KEY (star_id)
            REFERENCES stars (id),
    CONSTRAINT fk_orbit_id
        FOREIGN KEY (orbit_id)
            REFERENCES orbits (id),
    CONSTRAINT fk_deposit_kind_cd
        FOREIGN KEY (kind)
            REFERENCES unit_codes (code)
);

-- deposits_summary is a temporary table that is used to summarize the
-- deposits table for reporting. It holds the total quantity of each
-- resource type from all the deposits on a planet. It also holds the
-- log10 estimated value of the resources.
create table deposits_summary
(
    deposit_id   integer not null,
    eff_turn     integer not null,
    end_turn     integer not null,
    fuel_qty     integer not null,
    fuel_est_qty integer not null,
    gold_qty     integer not null,
    gold_est_qty integer not null,
    mets_qty     integer not null,
    mets_est_qty integer not null,
    nmts_qty     integer not null,
    nmts_est_qty integer not null,
    CONSTRAINT fk_deposit_id
        FOREIGN KEY (deposit_id)
            REFERENCES deposits (id)
);

CREATE TABLE empires
(
    id        INTEGER NOT NULL UNIQUE CHECK (id BETWEEN 1 AND 256),
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1))
);

CREATE TABLE empire
(
    empire_id      INTEGER NOT NULL UNIQUE,
    empire_name    TEXT    NOT NULL UNIQUE,
    username       TEXT    NOT NULL UNIQUE,
    email          TEXT    NOT NULL UNIQUE,
    home_system_id INTEGER NOT NULL,
    home_star_id   INTEGER NOT NULL,
    home_orbit_id  INTEGER NOT NULL,
    CONSTRAINT fk_empire_id
        FOREIGN KEY (empire_id)
            REFERENCES empires (id),
    CONSTRAINT fk_system_id
        FOREIGN KEY (home_system_id)
            REFERENCES systems (id),
    CONSTRAINT fk_star_id
        FOREIGN KEY (home_star_id)
            REFERENCES stars (id),
    CONSTRAINT fk_orbit_id
        FOREIGN KEY (home_orbit_id)
            REFERENCES orbits (id)
);

CREATE TABLE sc_codes
(
    code TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
INSERT INTO sc_codes (code, name)
VALUES ('NONE', 'none');
INSERT INTO sc_codes (code, name)
VALUES ('SHIP', 'Ship');
INSERT INTO sc_codes (code, name)
VALUES ('COPN', 'Open Surface Colony');
INSERT INTO sc_codes (code, name)
VALUES ('CENC', 'Enclosed Surface Colony');
INSERT INTO sc_codes (code, name)
VALUES ('CORB', 'Orbital Colony');

CREATE TABLE scs
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    empire_id     INTEGER NOT NULL,
    sc_cd         TEXT    NOT NULL,
    sc_tech_level INTEGER NOT NULL CHECK (sc_tech_level BETWEEN 1 AND 10),
    name          TEXT    NOT NULL,
    location      INTEGER NOT NULL,
    is_on_surface INTEGER NOT NULL DEFAULT 0 CHECK (is_on_surface IN (0, 1)),
    rations       REAL    NOT NULL CHECK (rations >= 0),
    sol           REAL    NOT NULL CHECK (sol >= 0),
    birth_rate    REAL    NOT NULL CHECK (birth_rate >= 0),
    death_rate    REAL    NOT NULL CHECK (death_rate >= 0),
    CONSTRAINT fk_empire_id
        FOREIGN KEY (empire_id)
            REFERENCES empires (id),
    CONSTRAINT fk_sc_cd
        FOREIGN KEY (sc_cd)
            REFERENCES sc_codes (code),
    CONSTRAINT fk_orbit_id
        FOREIGN KEY (location)
            REFERENCES orbits (id)
);

CREATE TABLE inventory
(
    sc_id        INTEGER NOT NULL,
    unit_cd      TEXT    NOT NULL,
    tech_level   INTEGER NOT NULL CHECK (tech_level BETWEEN 0 AND 10),
    qty          INTEGER NOT NULL CHECK (qty >= 0),
    mass         REAL    NOT NULL CHECK (mass >= 0),
    volume       REAL    NOT NULL CHECK (volume >= 0),
    is_assembled INTEGER NOT NULL DEFAULT 0 CHECK (is_assembled IN (0, 1)),
    is_stored    INTEGER NOT NULL DEFAULT 0 CHECK (is_stored IN (0, 1)),
    CONSTRAINT fk_sc_id
        FOREIGN KEY (sc_id)
            REFERENCES scs (id),
    CONSTRAINT fk_unit_cd
        FOREIGN KEY (unit_cd)
            REFERENCES unit_codes (code)
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
    sc_id         INTEGER NOT NULL,
    population_cd TEXT    NOT NULL,
    qty           INTEGER NOT NULL CHECK (qty >= 0),
    pay_rate      REAL    NOT NULL CHECK (pay_rate >= 0),
    rebel_qty     INTEGER NOT NULL CHECK (rebel_qty >= 0),
    PRIMARY KEY (sc_id, population_cd),
    CONSTRAINT fk_sc_id
        FOREIGN KEY (sc_id)
            REFERENCES scs (id),
    CONSTRAINT fk_population_cd
        FOREIGN KEY (population_cd)
            REFERENCES population_codes (code)
);

CREATE TABLE factory_groups
(
    sc_id             INTEGER NOT NULL,
    group_no          INTEGER NOT NULL CHECK (group_no BETWEEN 1 AND 30),
    orders_cd         TEXT    NOT NULL,
    orders_tech_level INTEGER NOT NULL CHECK (orders_tech_level BETWEEN 0 AND 10),
    PRIMARY KEY (sc_id, group_no),
    CONSTRAINT fk_sc_id
        FOREIGN KEY (sc_id)
            REFERENCES scs (id),
    CONSTRAINT fk_orders_cd
        FOREIGN KEY (orders_cd)
            REFERENCES unit_codes (code)
            ON DELETE CASCADE
);

CREATE TABLE factory_group
(
    sc_id             INTEGER NOT NULL,
    group_no          INTEGER NOT NULL,
    group_tech_level  INTEGER NOT NULL CHECK (group_tech_level BETWEEN 1 AND 10),
    nbr_of_units      INTEGER NOT NULL CHECK (nbr_of_units >= 0),
    orders_cd         TEXT    NOT NULL,
    orders_tech_level INTEGER NOT NULL CHECK (orders_tech_level BETWEEN 0 AND 10),
    wip_25pct_qty     INTEGER NOT NULL CHECK (wip_25pct_qty >= 0),
    wip_50pct_qty     INTEGER NOT NULL CHECK (wip_50pct_qty >= 0),
    wip_75pct_qty     INTEGER NOT NULL CHECK (wip_75pct_qty >= 0),
    PRIMARY KEY (sc_id, group_no, group_tech_level),
    CONSTRAINT fk_factory_group_id
        FOREIGN KEY (sc_id, group_no)
            REFERENCES factory_groups (sc_id, group_no)
);

CREATE TABLE factory_group_retool
(
    sc_id             INTEGER NOT NULL,
    group_no          INTEGER NOT NULL,
    turn_no           INTEGER NOT NULL,
    orders_cd         TEXT    NOT NULL,
    orders_tech_level INTEGER NOT NULL CHECK (orders_tech_level BETWEEN 0 AND 10),
    PRIMARY KEY (sc_id, group_no),
    CONSTRAINT fk_sc_id
        FOREIGN KEY (sc_id)
            REFERENCES scs (id),
    CONSTRAINT fk_orders_cd
        FOREIGN KEY (orders_cd)
            REFERENCES unit_codes (code)
);

CREATE TABLE farm_groups
(
    sc_id    INTEGER NOT NULL,
    group_no INTEGER NOT NULL CHECK (group_no BETWEEN 1 AND 30),
    PRIMARY KEY (sc_id, group_no),
    CONSTRAINT fk_sc_id
        FOREIGN KEY (sc_id)
            REFERENCES scs (id)
);

CREATE TABLE farm_group
(
    sc_id            INTEGER NOT NULL,
    group_no         INTEGER NOT NULL,
    group_tech_level INTEGER NOT NULL CHECK (group_tech_level BETWEEN 1 AND 10),
    nbr_of_units     INTEGER NOT NULL CHECK (nbr_of_units >= 0),
    PRIMARY KEY (sc_id, group_no, group_tech_level),
    CONSTRAINT fk_farm_groups_id
        FOREIGN KEY (sc_id, group_no)
            REFERENCES farm_groups (sc_id, group_no)
);

CREATE TABLE mining_groups
(
    sc_id      INTEGER NOT NULL,
    group_no   INTEGER NOT NULL CHECK (group_no BETWEEN 1 AND 35),
    deposit_id INTEGER NOT NULL,
    PRIMARY KEY (sc_id, group_no),
    CONSTRAINT fk_sc_id
        FOREIGN KEY (sc_id)
            REFERENCES scs (id),
    CONSTRAINT fk_deposit_id
        FOREIGN KEY (deposit_id)
            REFERENCES deposits (id)
);

CREATE TABLE mining_group
(
    sc_id            INTEGER NOT NULL,
    group_no         INTEGER NOT NULL CHECK (group_no BETWEEN 1 AND 35),
    group_tech_level INTEGER NOT NULL CHECK (group_tech_level BETWEEN 1 AND 10),
    nbr_of_units     INTEGER NOT NULL CHECK (nbr_of_units >= 0),
    PRIMARY KEY (sc_id, group_no, group_tech_level),
    CONSTRAINT fk_mining_group_id
        FOREIGN KEY (sc_id, group_no)
            REFERENCES mining_groups (sc_id, group_no)
);

CREATE TABLE probe_orders
(
    sc_id     INTEGER NOT NULL,
    kind      TEXT    NOT NULL CHECK (kind in ('colony', 'orbit', 'ship', 'star', 'system')),
    target_id INTEGER NOT NULL,
    status    TEXT    NOT NULL,
    CONSTRAINT fk_sc_id
        FOREIGN KEY (sc_id)
            REFERENCES scs (id)
);

CREATE TABLE survey_orders
(
    sc_id     INTEGER NOT NULL,
    target_id INTEGER NOT NULL,
    status    TEXT    NOT NULL,
    CONSTRAINT fk_sc_id
        FOREIGN KEY (sc_id)
            REFERENCES scs (id),
    CONSTRAINT fk_target_id
        FOREIGN KEY (target_id)
            REFERENCES orbits (id)
);

-- sc_farming_summary is a summary of the farming activity of a single group
-- on a ship/colony for a single turn.
create table sc_farming_summary
(
    sc_id         integer not null,
    group_no      integer not null,
    turn_no       integer not null,
    fuel_consumed integer not null,
    pro_consumed  integer not null,
    usk_consumed  integer not null,
    aut_consumed  integer not null,
    food_produced integer not null,
    primary key (sc_id, group_no, turn_no),
    constraint fk_sc_id
        foreign key (sc_id)
            references scs (id)
);

-- sc_manufacturing_summary is a summary of the manufacturing activity of a single group
-- on a ship/colony for a single turn. The primary key is the sc_id, group_no, and turn_no.
-- This works because each factory group can only create one type of unit per turn.
create table sc_manufacturing_summary
(
    sc_id           integer not null,
    group_no        integer not null,
    turn_no         integer not null,
    fuel_consumed   integer not null,
    mets_consumed   integer not null,
    nmts_consumed   integer not null,
    pro_consumed    integer not null,
    usk_consumed    integer not null,
    aut_consumed    integer not null,
    unit_cd         text    not null,
    unit_tech_level integer not null,
    units_produced  integer not null,
    primary key (sc_id, group_no, turn_no),
    constraint fk_sc_id
        foreign key (sc_id)
            references scs (id)
);

-- sc_mining_summary is a summary of the mining activity of a single group
-- on a ship/colony for a single turn. The primary key is the sc_id, group_no,
-- and turn_no. This works because each mining group can only mine one deposit
-- per turn.
create table sc_mining_summary
(
    sc_id         integer not null,
    group_no      integer not null,
    turn_no       integer not null,
    fuel_consumed integer not null,
    pro_consumed  integer not null,
    usk_consumed  integer not null,
    aut_consumed  integer not null,
    fuel_produced integer not null,
    gold_produced integer not null,
    mets_produced integer not null,
    nmts_produced integer not null,
    primary key (sc_id, group_no, turn_no),
    constraint fk_sc_id
        foreign key (sc_id)
            references scs (id)
);
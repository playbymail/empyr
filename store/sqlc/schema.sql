--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- foreign keys must be disabled to drop tables
PRAGMA foreign_keys = OFF;
DROP TABLE IF EXISTS meta_migrations;
DROP TABLE IF EXISTS codes;
DROP TABLE IF EXISTS colonies;
DROP TABLE IF EXISTS sorc_details;
DROP TABLE IF EXISTS sorc_infrastructure;
DROP TABLE IF EXISTS sorc_inventory;
DROP TABLE IF EXISTS sorc_population;
DROP TABLE IF EXISTS sorc_superstructure;
DROP TABLE IF EXISTS deposit;
DROP TABLE IF EXISTS deposits;
DROP TABLE IF EXISTS empires;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS planets;
DROP TABLE IF EXISTS players;
DROP TABLE IF EXISTS orbits;
DROP TABLE IF EXISTS stars;
DROP TABLE IF EXISTS system_distances;
DROP TABLE IF EXISTS system_stars;
DROP TABLE IF EXISTS systems;
DROP TABLE IF EXISTS units;

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

CREATE TABLE codes
(
    category TEXT    NOT NULL,
    code     TEXT    NOT NULL,
    value    INTEGER NOT NULL,
    display  TEXT    NOT NULL,
    UNIQUE (category, code),
    UNIQUE (category, value)
);

CREATE TABLE players
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    handle     TEXT     NOT NULL UNIQUE,
    magic_link TEXT     NOT NULL UNIQUE,
    is_active  INTEGER  NOT NULL CHECK (is_active IN (0, 1)),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE units
(
    code           TEXT    NOT NULL,
    mass           INTEGER NOT NULL CHECK (mass BETWEEN 0 AND 100),
    is_operational INTEGER NOT NULL CHECK (is_operational IN (0, 1)),
    UNIQUE (code)
);
-- insert into units (code, mass, is_operational) values ('fighter', 1, 1);
-- -- operational units are: Space drives, sensors, automation units, life support units, energy weapons, energy shields, mining units, factory units, farms, hyper engines, structural units, light structural units, and missile launchers

CREATE TABLE games
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    code           TEXT     NOT NULL UNIQUE,
    name           TEXT     NOT NULL UNIQUE,
    display_name   TEXT     NOT NULL UNIQUE,
    current_turn   INTEGER  NOT NULL DEFAULT 0,
    last_empire_no INTEGER  NOT NULL DEFAULT 0,
    home_system_id INTEGER  NOT NULL DEFAULT 0,
    home_star_id   INTEGER  NOT NULL DEFAULT 0,
    home_orbit_id  INTEGER  NOT NULL DEFAULT 0,
    home_planet_id INTEGER  NOT NULL DEFAULT 0,
    created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE systems
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id  INTEGER NOT NULL,
    x        INTEGER NOT NULL CHECK (x BETWEEN 0 AND 30),
    y        INTEGER NOT NULL CHECK (y BETWEEN 0 AND 30),
    z        INTEGER NOT NULL CHECK (z BETWEEN 0 AND 30),
    scarcity INTEGER NOT NULL CHECK (scarcity BETWEEN 0 AND 3),
    CONSTRAINT fk_game_id
        FOREIGN KEY (game_id)
            REFERENCES games (id)
            ON DELETE CASCADE
);

CREATE TABLE system_distances
(
    from_system_id INTEGER NOT NULL,
    to_system_id   INTEGER NOT NULL,
    distance       INTEGER NOT NULL CHECK (distance BETWEEN 0 AND 100),
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

CREATE TABLE stars
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    system_id INTEGER NOT NULL,
    sequence  TEXT    NOT NULL CHECK (LENGTH(sequence) = 1 AND sequence BETWEEN 'A' AND 'Z'),
    scarcity  INTEGER NOT NULL CHECK (scarcity BETWEEN 0 AND 3),
    UNIQUE (system_id, sequence),
    CONSTRAINT fk_system_id
        FOREIGN KEY (system_id)
            REFERENCES systems (id)
            ON DELETE CASCADE
);

CREATE TABLE orbits
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    star_id  INTEGER NOT NULL,
    orbit_no INTEGER NOT NULL CHECK (orbit_no BETWEEN 1 AND 10),
    kind     INTEGER NOT NULL CHECK (kind BETWEEN 0 AND 5),
    scarcity INTEGER NOT NULL CHECK (scarcity BETWEEN 0 AND 3),
    UNIQUE (star_id, orbit_no),
    CONSTRAINT fk_star_id
        FOREIGN KEY (star_id)
            REFERENCES stars (id)
            ON DELETE CASCADE
);

CREATE TABLE planets
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    orbit_id     INTEGER NOT NULL,
    kind         INTEGER NOT NULL CHECK (kind BETWEEN 0 AND 4),
    habitability INTEGER NOT NULL CHECK (habitability BETWEEN 0 AND 25),
    scarcity     INTEGER NOT NULL CHECK (scarcity BETWEEN 0 AND 3),
    CONSTRAINT fk_orbit_id
        FOREIGN KEY (orbit_id)
            REFERENCES orbits (id)
            ON DELETE CASCADE
);

CREATE TABLE deposits
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    planet_id     INTEGER NOT NULL,
    deposit_no    INTEGER NOT NULL CHECK (deposit_no BETWEEN 1 AND 35),
    kind          INTEGER NOT NULL CHECK (kind BETWEEN 0 AND 4),
    initial_qty   INTEGER NOT NULL CHECK (initial_qty BETWEEN 0 AND 99000000),
    remaining_qty INTEGER NOT NULL CHECK (remaining_qty BETWEEN 0 AND 99000000),
    yield_pct     INTEGER NOT NULL CHECK (yield_pct BETWEEN 0 AND 100),
    UNIQUE (planet_id, deposit_no),
    CONSTRAINT fk_planet_id
        FOREIGN KEY (planet_id)
            REFERENCES planets (id)
            ON DELETE CASCADE
);

CREATE TABLE deposit
(
    deposit_id INTEGER NOT NULL,
    quantity   INTEGER NOT NULL CHECK (quantity BETWEEN 0 AND 99000000),
    yield_pct  INTEGER NOT NULL CHECK (yield_pct BETWEEN 0 AND 100),
    eff_turn   INTEGER NOT NULL,
    end_turn   INTEGER NOT NULL,
    active     INTEGER NOT NULL CHECK (active IN (0, 1)),
    CONSTRAINT fk_deposit_id
        FOREIGN KEY (deposit_id)
            REFERENCES deposits (id)
            ON DELETE CASCADE
);

CREATE TABLE empires
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id        INTEGER NOT NULL,
    player_id      INTEGER NOT NULL,
    empire_no      INTEGER NOT NULL CHECK (empire_no BETWEEN 1 AND 256),
    name           TEXT    NOT NULL,
    home_system_id INTEGER NOT NULL,
    home_star_id   INTEGER NOT NULL,
    home_orbit_id  INTEGER NOT NULL,
    home_planet_id INTEGER NOT NULL,
    UNIQUE (game_id, player_id),
    UNIQUE (game_id, empire_no),
    CONSTRAINT fk_game_id
        FOREIGN KEY (game_id)
            REFERENCES games (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_player_id
        FOREIGN KEY (player_id)
            REFERENCES players (id)
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

CREATE TABLE sorcs
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    empire_id INTEGER NOT NULL,
    kind      INTEGER NOT NULL CHECK (kind BETWEEN 1 AND 4),
    CONSTRAINT fk_empire_id
        FOREIGN KEY (empire_id)
            REFERENCES empires (id)
            ON DELETE CASCADE
);

CREATE TABLE sorc_details
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    sorc_id       INTEGER NOT NULL,
    turn_no       INTEGER NOT NULL CHECK (turn_no >= 0),
    tech_level    INTEGER NOT NULL CHECK (tech_level BETWEEN 0 AND 10),
    name          TEXT    NOT NULL,
    uem_qty       INTEGER NOT NULL CHECK (uem_qty >= 0),
    uem_pay       REAL    NOT NULL CHECK (uem_pay >= 0),
    usk_qty       INTEGER NOT NULL CHECK (usk_qty >= 0),
    usk_pay       REAL    NOT NULL CHECK (usk_pay >= 0),
    pro_qty       INTEGER NOT NULL CHECK (pro_qty >= 0),
    pro_pay       REAL    NOT NULL CHECK (pro_pay >= 0),
    sld_qty       INTEGER NOT NULL CHECK (sld_qty >= 0),
    sld_pay       REAL    NOT NULL CHECK (sld_pay >= 0),
    cnw_qty       INTEGER NOT NULL CHECK (cnw_qty >= 0),
    spy_qty       INTEGER NOT NULL CHECK (spy_qty >= 0),
    rations       REAL    NOT NULL CHECK (rations >= 0),
    birth_rate    REAL    NOT NULL CHECK (birth_rate >= 0),
    death_rate    REAL    NOT NULL CHECK (death_rate >= 0),
    sol           REAL    NOT NULL CHECK (sol >= 0),
    orbit_id      INTEGER NOT NULL,
    is_on_surface INTEGER NOT NULL CHECK (is_on_surface IN (0, 1)),
    UNIQUE (sorc_id, turn_no),
    CONSTRAINT fk_sorc_id
        FOREIGN KEY (sorc_id)
            REFERENCES sorcs (id)
            ON DELETE CASCADE
);

CREATE TABLE sorc_infrastructure
(
    sorc_detail_id INTEGER NOT NULL,
    kind           TEXT    NOT NULL,
    tech_level     INTEGER NOT NULL CHECK (tech_level BETWEEN 0 AND 10),
    qty            INTEGER NOT NULL CHECK (qty >= 0),
    UNIQUE (sorc_detail_id, kind),
    CONSTRAINT fk_sorc_detail_id
        FOREIGN KEY (sorc_detail_id)
            REFERENCES sorc_details (id)
            ON DELETE CASCADE
);

CREATE TABLE sorc_inventory
(
    sorc_detail_id INTEGER NOT NULL,
    kind           TEXT    NOT NULL,
    tech_level     INTEGER NOT NULL CHECK (tech_level BETWEEN 0 AND 10),
    qty_assembled  INTEGER NOT NULL CHECK (qty_assembled >= 0),
    qty_stored     INTEGER NOT NULL CHECK (qty_stored >= 0),
    UNIQUE (sorc_detail_id, kind),
    CONSTRAINT fk_sorc_detail_id
        FOREIGN KEY (sorc_detail_id)
            REFERENCES sorc_details (id)
            ON DELETE CASCADE
);

CREATE TABLE sorc_population
(
    sorc_detail_id INTEGER NOT NULL,
    kind           TEXT    NOT NULL,
    qty            INTEGER NOT NULL CHECK (qty >= 0),
    UNIQUE (sorc_detail_id, kind),
    CONSTRAINT fk_sorc_detail_id
        FOREIGN KEY (sorc_detail_id)
            REFERENCES sorc_details (id)
            ON DELETE CASCADE
);

CREATE TABLE sorc_superstructure
(
    sorc_detail_id INTEGER NOT NULL,
    kind           TEXT    NOT NULL,
    tech_level     INTEGER NOT NULL CHECK (tech_level BETWEEN 0 AND 10),
    qty            INTEGER NOT NULL CHECK (qty >= 0),
    UNIQUE (sorc_detail_id, kind),
    CONSTRAINT fk_sorc_detail_id
        FOREIGN KEY (sorc_detail_id)
            REFERENCES sorc_details (id)
            ON DELETE CASCADE
);

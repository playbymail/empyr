--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- foreign keys must be disabled to drop tables
PRAGMA foreign_keys = OFF;
DROP TABLE IF EXISTS meta_migrations;
DROP TABLE IF EXISTS colonies;
DROP TABLE IF EXISTS deposit;
DROP TABLE IF EXISTS deposits;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS planets;
DROP TABLE IF EXISTS orbits;
DROP TABLE IF EXISTS stars;
DROP TABLE IF EXISTS system_distances;
DROP TABLE IF EXISTS system_stars;
DROP TABLE IF EXISTS systems;

-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

-- -- create the table for managing migrations
-- CREATE TABLE meta_migrations
-- (
--     version    INTEGER  NOT NULL UNIQUE,
--     comment    TEXT     NOT NULL,
--     script     TEXT     NOT NULL UNIQUE,
--     created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
-- );

-- -- update the migrations table
-- INSERT INTO meta_migrations (version, comment, script)
-- VALUES (202502110915, 'initial migration', '202502110915_initial.sql');

CREATE TABLE games
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    code         TEXT     NOT NULL UNIQUE,
    name         TEXT     NOT NULL UNIQUE,
    display_name TEXT     NOT NULL UNIQUE,
    current_turn INTEGER  NOT NULL DEFAULT 0,
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
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
    kind     TEXT    NOT NULL CHECK (kind IN ('empty', 'gas-giant', 'terrestrial', 'asteroid-belt')),
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
    kind         TEXT    NOT NULL CHECK (kind IN ('empty', 'gas-giant', 'terrestrial', 'asteroid-belt')),
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
    kind          TEXT    NOT NULL CHECK (kind IN ('none', 'gold', 'fuel', 'metallic', 'non-metallic')),
    yield_pct     INTEGER NOT NULL CHECK (yield_pct BETWEEN 0 AND 100),
    initial_qty   INTEGER NOT NULL CHECK (initial_qty BETWEEN 0 AND 99000000),
    remaining_qty INTEGER NOT NULL CHECK (remaining_qty BETWEEN 0 AND 99000000),
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

CREATE TABLE colonies
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    planet_id INTEGER NOT NULL REFERENCES planets (id),
    kind      TEXT    NOT NULL CHECK (kind IN ('open', 'enclosed')),
    location  TEXT    NOT NULL CHECK (location IN ('surface', 'orbital'))
);

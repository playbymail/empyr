--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- foreign keys must be disabled to drop tables
PRAGMA foreign_keys = OFF;
DROP TABLE IF EXISTS meta_migrations;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS systems;
DROP TABLE IF EXISTS stars;
DROP TABLE IF EXISTS planets;
DROP TABLE IF EXISTS planet_deposits;
DROP TABLE IF EXISTS planet_deposit;
DROP TABLE IF EXISTS colonies;

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
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id INTEGER NOT NULL REFERENCES games (id),
    x       INTEGER NOT NULL CHECK (x BETWEEN 0 AND 30),
    y       INTEGER NOT NULL CHECK (y BETWEEN 0 AND 30),
    z       INTEGER NOT NULL CHECK (z BETWEEN 0 AND 30)
);

CREATE TABLE stars
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    system_id INTEGER NOT NULL REFERENCES systems (id),
    sequence  TEXT    NOT NULL CHECK (LENGTH(sequence) = 1 AND sequence BETWEEN 'A' AND 'Z'),
    UNIQUE (system_id, sequence)
);

CREATE TABLE planets
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    star_id      INTEGER NOT NULL REFERENCES stars (id),
    orbit        INTEGER NOT NULL CHECK (orbit BETWEEN 1 AND 10),
    kind         TEXT    NOT NULL CHECK (kind IN ('gas-giant', 'terrestrial', 'asteroid')),
    habitability INTEGER NOT NULL CHECK (habitability BETWEEN 0 AND 25),
    UNIQUE (star_id, orbit)
);

CREATE TABLE planet_deposits
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    planet_id  INTEGER NOT NULL REFERENCES planets (id),
    deposit_no INTEGER NOT NULL CHECK (deposit_no BETWEEN 1 AND 35),
    kind       TEXT    NOT NULL CHECK (kind IN ('gold', 'fuel', 'metallic', 'non-metallic')),
    UNIQUE (planet_id, deposit_no)
);

CREATE TABLE planet_deposit
(
    deposit_id INTEGER NOT NULL REFERENCES planet_deposits (id),
    quantity   INTEGER NOT NULL CHECK (quantity BETWEEN 0 AND 99000000),
    yield_pct  INTEGER NOT NULL CHECK (yield_pct BETWEEN 0 AND 100),
    eff_turn   INTEGER NOT NULL,
    end_turn   INTEGER NOT NULL,
    active     INTEGER NOT NULL CHECK (active IN (0, 1))
);

CREATE TABLE colonies
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    planet_id INTEGER NOT NULL REFERENCES planets (id),
    kind      TEXT    NOT NULL CHECK (kind IN ('open', 'enclosed')),
    location  TEXT    NOT NULL CHECK (location IN ('surface', 'orbital'))
);

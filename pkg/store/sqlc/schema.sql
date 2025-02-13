--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

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

-- update the migrations table
INSERT INTO meta_migrations (version, comment, script)
VALUES (202502110915, 'initial migration', '202502110915_initial.sql');
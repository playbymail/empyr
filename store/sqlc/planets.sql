--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreatePlanet creates a new planet.
--
-- name: CreatePlanet :one
INSERT INTO planets (orbit_id, kind, habitability)
VALUES (:orbit_id, :kind, :habitability)
RETURNING id;


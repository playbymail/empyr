-- CreateOrbit creates a new orbit.
--
-- name: CreateOrbit :one
INSERT INTO orbits (star_id, orbit_no, kind)
VALUES (:star_id, :orbit_no, :kind)
RETURNING id;

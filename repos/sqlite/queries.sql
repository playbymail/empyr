-- InsertMetaMigration inserts a new migration.
--
-- name: InsertMetaMigration :exec
INSERT INTO meta_migrations (version, comment, script)
VALUES (:version, :comment, :script);



--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- GetSchemaVersion gets the current schema version.
--
-- name: GetSchemaVersion :exec
SELECT max(version)
FROM meta_migrations;
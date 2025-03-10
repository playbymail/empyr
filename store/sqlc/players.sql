--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- ReadPlayerByHandle gets a player by its handle.
--
-- name: ReadPlayerByHandle :one
SELECT id, magic_link, is_active
FROM players
WHERE handle = :handle;
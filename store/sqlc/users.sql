--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- ReadUserByHandle gets a player by its handle.
--
-- name: ReadUserByHandle :one
SELECT id, magic_link, is_active
FROM users
WHERE handle = :handle;

-- ReadUserByMagicKey gets a player by its magic key.
--
-- name: ReadUserByMagicKey :one
SELECT id, handle, is_active
FROM users
WHERE magic_link = :magic_key;


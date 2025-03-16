--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- ReadUserByHandle gets a user by its handle.
--
-- name: ReadUserByHandle :one
SELECT id, magic_link, is_active, is_admin
FROM users
WHERE handle = :handle;

-- ReadUserByID gets a user by its id.
--
-- name: ReadUserByID :one
SELECT id, handle, is_active, is_admin
FROM users
WHERE id = :id;

-- ReadUserByMagicKey gets a user by its magic key.
--
-- name: ReadUserByMagicKey :one
SELECT id, is_active, is_admin
FROM users
WHERE handle = :handle
  AND magic_link = :magic_key;

-- ReadUsersInGame gets a list of users in a game.
--
-- name: ReadUsersInGame :many
SELECT users.id, users.handle, users.is_active, users.is_admin
FROM users,
     empires
WHERE empires.game_id = :game_id
  AND empires.user_id = users.id;


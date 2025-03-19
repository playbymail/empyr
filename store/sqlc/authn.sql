--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreateUser registers a new user.
--
-- name: CreateUser :one
INSERT INTO users (username, email, hashed_password, is_active, is_admin)
VALUES (:username, :email, :hashed_password, :is_active, :is_admin)
RETURNING id;

-- ReadUserByEmail gets a user by its email address.
--
-- name: ReadUserByEmail :one
SELECT id, username, email, hashed_password, is_active, is_admin
FROM users
WHERE email = :email;

-- ReadUserByID gets a user by its id.
--
-- name: ReadUserByID :one
SELECT id, username, email, hashed_password, is_active, is_admin
FROM users
WHERE id = :user_id;

-- ReadUserByUsername gets a user by its name.
--
-- name: ReadUserByUsername :one
SELECT id, username, email, hashed_password, is_active, is_admin
FROM users
WHERE username = :username;


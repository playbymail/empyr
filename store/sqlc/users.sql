--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- ReadUsersInGame gets a list of users in a game.
--
-- name: ReadUsersInGame :many
SELECT users.id, users.username, users.email, users.is_active, users.is_admin
FROM users,
     empires
WHERE empires.game_id = :game_id
  AND users.id = empires.user_id;


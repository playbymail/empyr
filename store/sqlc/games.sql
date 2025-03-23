--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- CreateGame creates a new game.
--
-- name: CreateGame :exec
insert into games (code, name, display_name)
values (:code, :name, :display_name);

-- ReadAllGameInfo returns a list of all games in the database, even the inactive ones.
--
-- name: ReadAllGameInfo :one
select code,
       name,
       display_name,
       current_turn,
       home_system_id,
       home_star_id,
       home_orbit_id
from games;

-- ReadCurrentTurn gets the current turn for a game.
--
-- name: ReadCurrentTurn :one
select current_turn
from games;

-- UpdateCurrentTurn increments the game turn number.
--
-- name: UpdateCurrentTurn :exec
update games
set current_turn = :turn_number;

-- UpdateGameHomeSystems updates the home system for a game.
--
-- name: UpdateGameHomeSystems :exec
update games
set home_system_id = :home_system_id,
    home_star_id   = :home_star_id,
    home_orbit_id  = :home_orbit_id;

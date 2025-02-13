// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"time"
)

// Player is an instance of a User in a specific Game.
type Player struct {
	ID      string
	Name    string // unique within a game, non-unique globally
	Created time.Time
	// Factories
	// Population
	// Weapons
	// Transportation
	//   Number
	//   Level
	//   Type
	// HomeNationID
	// HomePlanetID
	// RaceID
	// RaceName
}

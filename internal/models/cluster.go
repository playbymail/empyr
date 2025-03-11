// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package models

import "time"

// some notes to myself on the model:
// - this needs to implement the minimal viable product

// Player represents a single player in the game.
// A player is a human or AI or just an NPC.
// Each player controls a single empire in any given game.
type Player struct {
	ID int // unique identifier for the player
}

// Game represents a single game.
type Game struct {
	ID          int // unique identifier for the game
	Code        string
	Name        string
	DisplayName string
	CurrentTurn int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Cluster     *Cluster

	Empires []*Empire // list of empires in the game
}

// Empire is a single empire in the game.
// They may be controlled by a human player or an AI.
type Empire struct {
	ID       int
	PlayerID int // player that controls this empire
}

// Cluster defines the cluster of systems in the game.
// Every game has a single cluster of systems.
type Cluster struct {
	ID int // unique identifier for the cluster

	Systems []*System // list of systems in the cluster
}

// System is a stellar system containing one or more star systems.
// All systems use a 3D coordinate system. The location of a
// system is used to identify it in reports.
type System struct {
	ID int // unique identifier for the system

	// X, Y, Z are the coordinates of the system in 3D space.
	X, Y, Z int

	Stars []*Star // list of stars in the system
}

// Star is a single star in a stellar system.
// All stars have 10 orbits; orbits may be empty or contain a single planet.
type Star struct {
	ID       int // unique identifier for the star
	Sequence int // 1, 2, 3, etc

	// Orbits is a list of planets in the star system.
	// The zero index is always nil.
	Orbits [11]*Planet
}

type Planet struct {
	ID   int // unique identifier for the planet
	Kind Planet_e

	// Habitability is a rating for how many colonists can live on the surface
	// of a planet in open colonies. It ranges from 0 to 25.
	Habitability int

	// NaturalResources is a list of natural resources on the planet.
	// The zero index is always nil.
	NaturalResources [35]*Deposit
}

type Planet_e int

const (
	ASTEROID_BELT Planet_e = iota + 1
	GAS_GIANT
	TERRESTRIAL
)

type Deposit struct {
	ID       int        // unique identifier for the deposit
	Kind     Resource_e // the type of resource
	Quantity int        // current quantity of the resource
	YieldPct int        // percentage of each unit mined that can be refined
}

type Resource_e int

const (
	FUEL Resource_e = iota + 1
	GOLD
	METALLICS
	NONMETALLICS
)

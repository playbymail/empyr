// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import "github.com/playbymail/empyr/store"

type Engine_t struct {
	Store *store.Store
}

type Resource_e int64

const (
	GOLD Resource_e = iota
	FUEL
	METALLICS
	NON_METALLICS
)

type Deposit_t struct {
	Id       int64
	Resource Resource_e
	Quantity int64
	Yield    int64
}

type Scarcity_e int64

const (
	TYPICAL Scarcity_e = iota
	RICH
	POOR
)

type Cluster_t struct {
	Systems []*System_t
	Stars   []*Star_t
	Orbits  []*Orbit_t
	Planets []*Planet_t
}

type Point_t struct {
	X, Y, Z int64
}

type System_t struct {
	Id          int64
	Coordinates Point_t
	Scarcity    Scarcity_e
	Stars       []int64 // index into Stars
}

type Star_t struct {
	Id       int64
	SystemId int64  // index into Systems
	Sequence string // A ... D for the four stars in the system
	Scarcity Scarcity_e
	Orbits   [11]int64 // index into Orbits
}

type Orbit_t struct {
	Id      int64
	StarId  int64 // index into Stars
	OrbitNo int64 // value from 1 to 10 for this orbit
	Kind    Orbit_e
}

type Orbit_e int64

const (
	EmptyOrbit Orbit_e = iota
	AsteroidBelt
	EarthlikePlant
	GasGiant
	IceGiant
	RockyPlanet
)

type Planet_t struct {
	Id   int64
	Star int64 // index into Stars
	Kind Planet_e
}

type Planet_e int64

const (
	NoPlanet Planet_e = iota
	AsteroidBeltPlanet
	GasGiantPlanet
	TerrestrialPlanet
)

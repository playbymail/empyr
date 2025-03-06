// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

type Resource_e int

const (
	GOLD Resource_e = iota
	FUEL
	METALLICS
	NON_METALLICS
)

type Deposit_t struct {
	Id       int
	Resource Resource_e
	Quantity int
	Yield    int
}

type Scarcity_e int

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
	X, Y, Z int
}

type System_t struct {
	Id          int
	Coordinates Point_t
	Scarcity    Scarcity_e
	Stars       []int // index into Stars
}

type Star_t struct {
	Id       int
	System   int    // index into Systems
	Sequence string // A ... D for the four stars in the system
	Scarcity Scarcity_e
	Orbits   [11]int // index into Orbits
}

type Orbit_t struct {
	Id      int
	Star    int // index into Stars
	OrbitNo int // value from 1 to 10 for this orbit
	Kind    Orbit_e
}

type Orbit_e int

const (
	EmptyOrbit Orbit_e = iota
	AsteroidBelt
	EarthlikePlant
	GasGiant
	IceGiant
	RockyPlanet
)

type Planet_t struct {
	Id   int
	Star int // index into Stars
	Kind Planet_e
}

type Planet_e int

const (
	NoPlanet Planet_e = iota
	AsteroidBeltPlanet
	GasGiantPlanet
	TerrestrialPlanet
)

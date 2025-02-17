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
	Systems [101]*System_t
}

type Point_t struct {
	X, Y, Z int
}

type System_t struct {
	Id          int
	Coordinates Point_t
	Stars       []*Star_t
}

type Star_t struct {
	Id       int
	Sequence string // A ... D for the four stars in the system
	Orbits   [11]*Orbit_t
}

type Orbit_t struct {
	Id     int
	Planet *Planet_t
}

type Planet_t struct {
	Id   int
	Kind Planet_e
}

type Planet_e int

const (
	TERRESTRIAL Planet_e = iota
	GAS_GIANT
	ASTEROID_BELT
)

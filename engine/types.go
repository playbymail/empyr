// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import "github.com/playbymail/empyr/store"

type Engine_t struct {
	Store *store.Store
}

type Resource_e int64

const (
	NONE Resource_e = iota
	GOLD
	FUEL
	METALLICS
	NON_METALLICS
)

type Deposit_t struct {
	Id        int64
	Planet    *Planet_t
	DepositNo int64
	Resource  Resource_e
	Quantity  int64
	Yield     int64
}

// Less is a helper for sorting deposits.
// It is used to sort deposits by resource type.
// It sorts NONE last, and everything else by resource type.
func (d *Deposit_t) Less(d2 *Deposit_t) bool {
	if d.Resource == NONE {
		return false
	} else if d2.Resource == NONE {
		return true
	}
	return d.Resource < d2.Resource
}

type Scarcity_e int64

const (
	TYPICAL Scarcity_e = iota
	RICH
	POOR
)

type Cluster_t struct {
	Systems  []*System_t
	Stars    []*Star_t
	Orbits   []*Orbit_t
	Planets  []*Planet_t
	Deposits []*Deposit_t
	Empires  []*Empire_t
}

type Point_t struct {
	X, Y, Z int64
}

type System_t struct {
	Id          int64
	Coordinates Point_t
	Scarcity    Scarcity_e
	Stars       []*Star_t
}

type Star_t struct {
	Id       int64
	System   *System_t
	Sequence string // A ... D for the four stars in the system
	Scarcity Scarcity_e
	Orbits   [11]*Orbit_t
}

type Orbit_t struct {
	Id       int64
	System   *System_t
	Star     *Star_t
	OrbitNo  int64 // value from 1 to 10 for this orbit
	Kind     Orbit_e
	Scarcity Scarcity_e
	Planet   *Planet_t
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
	Id           int64
	System       *System_t
	Star         *Star_t
	Orbit        *Orbit_t
	Kind         Planet_e
	Habitability int64 // 0..25
	Scarcity     Scarcity_e
	Deposits     [36]*Deposit_t
}

type Planet_e int64

const (
	NoPlanet Planet_e = iota
	AsteroidBeltPlanet
	GasGiantPlanet
	TerrestrialPlanet
)

type Empire_t struct {
	Id         int64
	EmpireNo   int64
	Name       string
	HomeSystem *System_t
	HomeStar   *Star_t
	HomeOrbit  *Orbit_t
	HomePlanet *Planet_t
}

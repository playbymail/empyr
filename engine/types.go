// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"fmt"
	"github.com/playbymail/empyr/store"
)

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

func (r Resource_e) Code() string {
	switch r {
	case NONE:
		return "NONE"
	case GOLD:
		return "GOLD"
	case FUEL:
		return "FUEL"
	case METALLICS:
		return "METS"
	case NON_METALLICS:
		return "NMTS"
	}
	return fmt.Sprintf("Resource_e(%d)", r)
}
func (r Resource_e) Descr() string {
	switch r {
	case NONE:
		return "None"
	case GOLD:
		return "Gold"
	case FUEL:
		return "Fuel"
	case METALLICS:
		return "Metallics"
	case NON_METALLICS:
		return "Non-metallics"
	}
	return fmt.Sprintf("Resource_e(%d)", r)
}

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
	Stars       []*Star_t
}

type Star_t struct {
	Id       int64
	System   *System_t
	Sequence string // A ... D for the four stars in the system
	Orbits   [11]*Orbit_t
}

type Orbit_t struct {
	Id      int64
	System  *System_t
	Star    *Star_t
	OrbitNo int64 // value from 1 to 10 for this orbit
	Kind    Orbit_e
	Planet  *Planet_t
}

type Orbit_e int64

const (
	EmptyOrbit Orbit_e = iota
	AsteroidBelt
	EarthlikePlanet
	GasGiant
	IceGiant
	RockyPlanet
)

func (e Orbit_e) Code() string {
	switch e {
	case EmptyOrbit:
		return "EMPTY"
	case AsteroidBelt:
		return "ASTR"
	case EarthlikePlanet:
		return "ERTH"
	case GasGiant:
		return "GASG"
	case IceGiant:
		return "ICEG"
	case RockyPlanet:
		return "RCKY"
	}
	return fmt.Sprintf("Orbit_e(%d)", e)
}

func (e Orbit_e) Descr() string {
	switch e {
	case EmptyOrbit:
		return "empty"
	case AsteroidBelt:
		return "asteroid belt"
	case EarthlikePlanet:
		return "earth-like planet"
	case GasGiant:
		return "gas giant"
	case IceGiant:
		return "ice giant"
	case RockyPlanet:
		return "rocky planet"
	}
	return fmt.Sprintf("Orbit_e(%d)", e)
}

type Planet_t struct {
	Id           int64
	System       *System_t
	Star         *Star_t
	Orbit        *Orbit_t
	Kind         Planet_e
	Habitability int64 // 0..25
	Deposits     [36]*Deposit_t
}

type Planet_e int64

const (
	NoPlanet Planet_e = iota
	AsteroidBeltPlanet
	GasGiantPlanet
	TerrestrialPlanet
)

func (e Planet_e) Code() string {
	switch e {
	case NoPlanet:
		return "NONE"
	case AsteroidBeltPlanet:
		return "ASTR"
	case GasGiantPlanet:
		return "GASG"
	case TerrestrialPlanet:
		return "TERR"
	}
	return fmt.Sprintf("Planet_e(%d)", e)
}
func (e Planet_e) Descr() string {
	switch e {
	case NoPlanet:
		return "no planet"
	case AsteroidBeltPlanet:
		return "asteroid belt"
	case GasGiantPlanet:
		return "gas giant"
	case TerrestrialPlanet:
		return "terrestrial planet"
	}
	return fmt.Sprintf("Planet_e(%d)", e)
}

type Empire_t struct {
	Id         int64
	EmpireNo   int64
	Name       string
	HomeSystem *System_t
	HomeStar   *Star_t
	HomeOrbit  *Orbit_t
	HomePlanet *Planet_t
}

type SorC_e int64

const (
	SCShip                  SorC_e = 1
	SCOpenSurfaceColony            = 2
	SCEnclosedSurfaceColony        = 3
	SCOrbitalColony                = 4
)

func (e SorC_e) Code() string {
	switch e {
	case SCShip:
		return "SHIP"
	case SCOpenSurfaceColony:
		return "COPN"
	case SCEnclosedSurfaceColony:
		return "CENC"
	case SCOrbitalColony:
		return "CORB"
	default:
		return "NONE"
	}
}

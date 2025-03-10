// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"log"
	"math/rand/v2"
)

type Game_t struct {
	Id int64 // database key

	Home struct {
		System *System_t
		Star   *Star_t
		Orbit  *Orbit_t
		Planet *Planet_t
	}

	Systems  map[int64]*System_t
	Stars    map[int64]*Star_t
	Orbits   map[int64]*Orbit_t
	Planets  map[int64]*Planet_t
	Deposits map[int64]*Deposit_t
	Empires  map[int64]*Empire_t
}

func newGame(r *rand.Rand) (*Game_t, error) {
	// create an empty cluster with about 100 systems
	cluster := emptyCluster(r)

	// use templates to populate the cluster
	for _, system := range cluster.Systems {
		if system.Id == 1 {
			homeSystemTemplate(system, r)
			log.Printf("home system: %p", system)
		} else {
			coreSystemTemplate(system, r)
		}
	}

	//// print the cluster map
	//origin := Point_t{X: 15, Y: 15, Z: 15}
	//for _, system := range cluster.Systems {
	//	if system.Id == 1 {
	//		log.Printf("system %3d: %02d %02d %02d d %8.4f stars %2d\n", system.Id, system.Coordinates.X, system.Coordinates.Y, system.Coordinates.Z, origin.DistanceTo(system.Coordinates), len(system.Stars))
	//		for _, star := range system.Stars {
	//			log.Printf("  star %3d: seq %s orbits %2d\n", star.Id, star.Sequence, len(star.Orbits))
	//			for _, orbit := range star.Orbits {
	//				if orbit != nil {
	//					planet := orbit.Planet
	//					log.Printf("    orbit %3d: orbit %2d kind %s\n", orbit.Id, orbit.OrbitNo, orbit.Kind.Code())
	//					log.Printf("   planet %3d: orbit %2d kind %s habitability %3d\n", planet.Id, orbit.OrbitNo, planet.Kind.Code(), planet.Habitability)
	//					log.Printf("               dep kind quantity______ yield\n")
	//					for _, deposit := range orbit.Planet.Deposits {
	//						if deposit != nil && deposit.Resource != NONE {
	//							log.Printf("               %3d %4s %14s %4d%%\n", deposit.DepositNo, deposit.Resource.Code(), commas(deposit.Quantity), deposit.Yield)
	//						}
	//					}
	//				}
	//			}
	//		}
	//	}
	//}

	g := &Game_t{
		Systems:  make(map[int64]*System_t),
		Stars:    make(map[int64]*Star_t),
		Orbits:   make(map[int64]*Orbit_t),
		Planets:  make(map[int64]*Planet_t),
		Deposits: make(map[int64]*Deposit_t),
		Empires:  make(map[int64]*Empire_t),
	}
	for _, system := range cluster.Systems {
		g.Systems[system.Id] = system
	}
	for _, star := range cluster.Stars {
		g.Stars[star.Id] = star
	}
	for _, orbit := range cluster.Orbits {
		g.Orbits[orbit.Id] = orbit
	}
	for _, planet := range cluster.Planets {
		g.Planets[planet.Id] = planet
	}
	for _, deposit := range cluster.Deposits {
		g.Deposits[deposit.Id] = deposit
	}
	g.Home.System = g.Systems[1]
	g.Home.Star = g.Home.System.Stars[0]
	g.Home.Orbit = g.Home.Star.Orbits[3]
	g.Home.Planet = g.Home.Orbit.Planet

	return g, nil
}

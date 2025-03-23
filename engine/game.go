// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

type Game_t struct {
	Code string // database key

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

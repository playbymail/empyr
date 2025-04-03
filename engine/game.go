// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import "sort"

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
	Entities map[int64]*Entity_t
}

// getInitialMiningSites is a helper function that returns a list of mining sites
// for the initial mining setup. It returns the four "best" mining sites, one for
// each resource type. Best is defined as the mining site with the highest yield.
func getInitialMiningSites(o *Orbit_t, d map[int64]*Deposit_t) (gold, fuel, mets, nmts *Deposit_t) {
	var deposits []*Deposit_t
	for _, v := range d {
		if v == nil || v.Resource == NONE || v.Planet.Orbit.Id != o.Id {
			continue
		}
		deposits = append(deposits, v)
	}
	// sort deposits by resource type and yield
	sort.Slice(deposits, func(i, j int) bool {
		if deposits[i].Resource == deposits[j].Resource {
			if deposits[i].Yield == deposits[j].Yield {
				return deposits[i].Quantity > deposits[j].Quantity
			}
			return deposits[i].Yield > deposits[j].Yield
		}
		return deposits[i].Resource < deposits[j].Resource
	})

	// get the deposit with the highest yield for each resource type
	for _, d := range deposits {
		switch d.Resource {
		case GOLD:
			if gold == nil {
				gold = d
			}
		case FUEL:
			if fuel == nil {
				fuel = d
			}
		case METALLICS:
			if mets == nil {
				mets = d
			}
		case NON_METALLICS:
			if nmts == nil {
				nmts = d
			}
		default:
			// ignore
		}
	}

	return gold, fuel, mets, nmts
}

// nextAvailableEmpireID returns the next available empire ID.
// That is the smallest ID that is not currently in use.
func (g *Game_t) nextAvailableEmpireID() int64 {
	// start looking from 1 and continue until we find an ID that isn't in the map
	var id int64 = 1
	for {
		if _, exists := g.Empires[id]; exists {
			break
		}
		id++
	}
	return id
}

// nextAvailableEntityID returns the next available entity ID.
// That is always one more than the largest ID in the map.
func (g *Game_t) nextAvailableEntityID() int64 {
	// start looking from 1 and continue until we find an ID that isn't in the map
	var id int64 = 0
	for k := range g.Entities {
		if k > id {
			id = k
		}
	}
	return id + 1
}

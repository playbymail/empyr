// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"fmt"
	"math/rand/v2"
)

// generateRings assigns entries to orbits naturally to a star system
// with enforced constraints on the number and location of planets.
//
// note: if the number of planets is not specified, it is randomly chosen.
//
// the slice returned uses asteroid belt, gas giant, earth-like, and rocky.
func generateRings(r *rand.Rand, numPlanets int) (orbits [11]string) {
	// Ensure valid range
	if numPlanets <= 0 {
		// assign a random number of planets, from 1 to 10
		numPlanets = r.IntN(5) + r.IntN(6) + 1
	}
	if numPlanets > 10 {
		numPlanets = 10
	}

	// helper to assign a kind to an orbit if possible
	assignOrbit := func(t bool, kind string, possibleRings ...int) {
		if numPlanets == 0 {
			return
		}
		var rings []int
		for _, ring := range possibleRings {
			if orbits[ring] == "" {
				rings = append(rings, ring)
			}
		}
		if len(rings) == 0 { // don't assign a planet if there are no rings available
			return
		}
		if !t {
			return
		}
		r.Shuffle(len(rings), func(i, j int) { rings[i], rings[j] = rings[j], rings[i] })
		orbits[rings[0]] = kind
		numPlanets = numPlanets - 1
	}

	// always create a gas giant
	assignOrbit(true, "gas giant", 6, 7, 8, 9)
	// create at least one asteroid belt, if possible
	assignOrbit(numPlanets > 1, "asteroid belt", 3, 4, 5)
	// 50% chance of another gas giant
	assignOrbit(numPlanets > 4 && r.IntN(2) == 0, "gas giant", 7, 8, 9)
	// 50% chance of another asteroid belt
	assignOrbit(numPlanets > 4 && r.IntN(2) == 0, "asteroid belt", 1, 4, 5, 6, 10)
	// 33% chance habitable, earth like planets (four attempts!)
	assignOrbit(r.IntN(3) == 0, "earth-like", 2, 3, 4, 5, 6)
	assignOrbit(r.IntN(3) == 0, "earth-like", 2, 3, 4, 5, 6)
	assignOrbit(r.IntN(3) == 0, "earth-like", 2, 3, 4, 5, 6)
	assignOrbit(r.IntN(3) == 0, "earth-like", 2, 3, 4, 5, 6)

	// any remaining orbits are kind of random
	var rings []int // will hold the rings that are not assigned
	for i := 1; i <= 10; i++ {
		if orbits[i] == "" {
			rings = append(rings, i)
		}
	}
	r.Shuffle(len(rings), func(i, j int) { rings[i], rings[j] = rings[j], rings[i] })

	// assign the remaining planets to the remaining orbits
	for ; numPlanets > 0 && len(rings) > 0; rings, numPlanets = rings[1:], numPlanets-1 {
		if ring := rings[0]; ring == 1 {
			orbits[ring] = "rocky"
		} else if ring < 6 {
			if r.IntN(10) == 0 {
				orbits[ring] = "earth-like"
			} else if r.IntN(10) < 2 {
				orbits[ring] = "asteroid belt"
			} else {
				orbits[ring] = "rocky"
			}
		} else if ring < 10 {
			if r.IntN(10) == 0 {
				orbits[ring] = "gas giant"
			} else if r.IntN(10) < 2 {
				orbits[ring] = "asteroid belt"
			} else {
				orbits[ring] = "rocky"
			}
		} else {
			if r.IntN(4) == 0 {
				orbits[ring] = "asteroid belt"
			} else {
				orbits[ring] = "rocky"
			}
		}
	}

	// this is a hack to ensure that there are no asteroid belts next to each other
	for i := 1; i <= 10; i++ {
		if orbits[i] == "asteroid belt" && orbits[i-1] == "asteroid belt" {
			if 2 <= i && i <= 5 && r.IntN(10) == 0 {
				orbits[i] = "earth-like"
			} else {
				orbits[i] = "rocky"
			}
		}
	}

	return orbits
}

// print system for debugging
func printSystem(orbits [11]string) {
	for i, orbit := range orbits {
		fmt.Printf("Orbit %d: %s\n", i+1, orbit)
	}
}

//func main() {
//	r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))
//	orbits := generateOrbits(r, r.IntN(5)+r.IntN(6)+1) // Example with 2..10 planets
//	printSystem(orbits)
//}

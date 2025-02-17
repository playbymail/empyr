// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"fmt"
	"math/rand/v2"
)

// this file implements the commands to create assets such as games, systems, and planets.

func CreateCluster(r *rand.Rand) (any, error) {
	// create a slice of points to randomly place most of the systems
	var points []Point_t
	var point Point_t
	// the location of the 4- and 3--star systems are fixed
	for _, point := range []Point_t{
		{X: 0, Y: 0, Z: 0},    // 4-star system
		{X: 1, Y: 1, Z: 1},    // 3-star system
		{X: -1, Y: -1, Z: -1}, // 3-star system
		{X: -2, Y: 2, Z: 0},   // 3-star system
	} {
		points = append(points, point)
	}
	// the 2-star systems are always within 7 units of the origin
	for n := 1; n <= 9; n++ {
		point = randomPoint(7)
		for maxAttempts := 0; tooClose(point, points) && maxAttempts < 1_000; maxAttempts++ {
			point = randomPoint(7)
		}
		points = append(points, point)
	}
	// the remaining single-star systems are always within the 31-unit cube centered at the origin
	for len(points) < 100 {
		point = randomCubePoint()
		for maxAttempts := 0; tooClose(point, points) && maxAttempts < 1_000; maxAttempts++ {
			point = randomCubePoint()
		}
		points = append(points, point)
	}

	// create a cluster
	var cluster Cluster_t
	// create 100 systems for the cluster
	for id := 1; len(points) != 0 && id < len(cluster.Systems); id++ {
		var numberOfStars int
		switch id {
		case 1: // 1 4-star system
			numberOfStars = 4
		case 2, 3, 4: // 3 3-star systems
			numberOfStars = 3
		case 5, 6, 7, 8, 9, 10, 11, 12, 13: // 9 2-star systems
			numberOfStars = 2
		default: // remaining are all 1-star systems
			numberOfStars = 1
		}
		// grab the pre-allocated point from the slice of points
		point, points = points[0], points[1:]
		// create the system using the point and number of stars
		system, err := createSystem(r, point, numberOfStars, TYPICAL)
		if err != nil {
			return nil, fmt.Errorf("system: %w", err)
		}
		cluster.Systems[id] = system
	}

	if len(points) != 0 {
		// we should have consumed all the points!
		panic(fmt.Sprintf("assert(len(points) != %d)", len(points)))
	}

	// todo: find a better place to generate the various id values
	return &cluster, nil
}

// createSystem creates a system with the given number of stars.
// The number of stars must be between 1 and 5.
func createSystem(r *rand.Rand, point Point_t, numberOfStars int, scarcity Scarcity_e) (*System_t, error) {
	if !(1 <= numberOfStars && numberOfStars <= 5) {
		return nil, fmt.Errorf("number of stars must be between 1 and 5")
	}
	// create the system
	system := &System_t{Coordinates: point}
	for i := 0; i < numberOfStars; i++ {
		star, err := createStar(r, scarcity)
		if err != nil {
			return nil, fmt.Errorf("star: %w", err)
		}
		star.Sequence = string(rune(65 + i))
		system.Stars = append(system.Stars, star)
	}
	return system, nil
}

// createStar creates a star with the given number of planets.
// The number of planets must be between 0 and 10.
func createStar(r *rand.Rand, scarcity Scarcity_e) (*Star_t, error) {
	star := &Star_t{}
	numberOfPlanets := r.IntN(5) + r.IntN(6) + 1 // normalRandInRange(r, 1, 10)
	rings := generateRings(r, numberOfPlanets)
	for i := 1; i <= 10; i++ {
		orbit, err := createOrbit(r)
		if err != nil {
			return nil, fmt.Errorf("orbit: %w", err)
		}
		star.Orbits[i] = orbit
		if rings[i] == "" {
			continue
		}
		orbit.Planet, err = createPlanet(r, rings[i], scarcity)
		if err != nil {
			return nil, fmt.Errorf("planet: %w", err)
		}
	}
	return star, nil
}

func createOrbit(r *rand.Rand) (*Orbit_t, error) {
	orbit := &Orbit_t{}
	return orbit, nil
}

// createPlanet creates a planet.
func createPlanet(r *rand.Rand, ring string, scarcity Scarcity_e) (*Planet_t, error) {
	planet := &Planet_t{}
	if ring == "asteroid belt" {
		planet.Kind = ASTEROID_BELT
	} else if ring == "gas giant" {
		planet.Kind = GAS_GIANT
	} else {
		planet.Kind = TERRESTRIAL
	}
	return planet, nil
}

// createDeposits creates natural resource deposits for a planet.
func createDeposits(r *rand.Rand, scarcity Scarcity_e) ([]*Deposit_t, error) {
	var numberOfDeposits int
	switch scarcity {
	case TYPICAL:
		numberOfDeposits = normalRandInRange(r, 1, 35)
	case RICH:
		numberOfDeposits = normalRandInRange(r, 16, 35)
	case POOR:
		numberOfDeposits = normalRandInRange(r, 1, 17)
	}
	var deposits []*Deposit_t
	for i := 0; i < numberOfDeposits; i++ {
		deposit, err := createDeposit(r, scarcity)
		if err != nil {
			return nil, fmt.Errorf("deposit: %w", err)
		}
		deposits = append(deposits, deposit)
	}
	return deposits, nil
}

// createDeposit creates a natural resource deposit.
func createDeposit(r *rand.Rand, scarcity Scarcity_e) (deposit *Deposit_t, err error) {
	var resource Resource_e
	switch scarcity {
	case TYPICAL:
		switch n := r.IntN(100); true {
		case n < 3:
			resource = GOLD
		case n < 22:
			resource = FUEL
		case n < 56:
			resource = METALLICS
		default:
			resource = NON_METALLICS
		}
	case RICH:
		switch n := r.IntN(100); true {
		case n < 5:
			resource = GOLD
		case n < 25:
			resource = FUEL
		case n < 58:
			resource = METALLICS
		default:
			resource = NON_METALLICS
		}
	case POOR:
		switch n := r.IntN(100); true {
		case n < 2:
			resource = FUEL
		case n < 85:
			resource = METALLICS
		default:
			resource = NON_METALLICS
		}
	}
	var minQuantity, maxQuantity, minYield, maxYield int
	switch resource {
	case GOLD:
		switch scarcity {
		case TYPICAL:
			minQuantity, maxQuantity = 1_000_000, 35_000_000
			minYield, maxYield = 1, 9
		case RICH:
			minQuantity, maxQuantity = 10_000_000, 35_000_000
			minYield, maxYield = 3, 9
		case POOR:
			minQuantity, maxQuantity = 1_000_000, 15_000_000
			minYield, maxYield = 1, 3
		}
	case FUEL:
		switch scarcity {
		case TYPICAL:
			minQuantity, maxQuantity = 1_000_000, 50_000_000
			minYield, maxYield = 1, 12
		case RICH:
			minQuantity, maxQuantity = 20_000_000, 50_000_000
			minYield, maxYield = 4, 12
		case POOR:
			minQuantity, maxQuantity = 1_000_000, 25_000_000
			minYield, maxYield = 1, 8
		}
	case METALLICS:
		switch scarcity {
		case TYPICAL:
			minQuantity, maxQuantity = 1_000_000, 99_000_000
			minYield, maxYield = 1, 36
		case RICH:
			minQuantity, maxQuantity = 25_000_000, 99_000_000
			minYield, maxYield = 12, 36
		case POOR:
			minQuantity, maxQuantity = 1_000_000, 75_000_000
			minYield, maxYield = 1, 18
		}
	case NON_METALLICS:
		switch scarcity {
		case TYPICAL:
			minQuantity, maxQuantity = 1_000_000, 99_000_000
			minYield, maxYield = 1, 36
		case RICH:
			minQuantity, maxQuantity = 25_000_000, 99_000_000
			minYield, maxYield = 12, 36
		case POOR:
			minQuantity, maxQuantity = 1_000_000, 75_000_000
			minYield, maxYield = 1, 18
		}
	}
	return &Deposit_t{
		Resource: resource,
		Quantity: normalRandInRange(r, minQuantity, maxQuantity),
		Yield:    normalRandInRange(r, minYield, maxYield),
	}, nil
}

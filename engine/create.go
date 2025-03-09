// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"errors"
	"fmt"
	"github.com/playbymail/empyr/store"
	"github.com/playbymail/empyr/store/sqlc"
	"log"
	"math"
	"math/rand/v2"
	"sort"
)

// this file implements the commands to create assets such as games, systems, and planets.

func Open(db *store.Store) (*Engine_t, error) {
	return &Engine_t{Store: db}, nil
}

func Close(e *Engine_t) error {
	if e != nil && e.Store != nil {
		err := e.Store.Close()
		e.Store = nil
		return err
	}
	return nil
}

func (e *Engine_t) CreateGame(code, name, displayName string, numberOfEmpires int64, calculateSystemDistances bool, r *rand.Rand, forceCreate bool) (int64, error) {
	cluster, err := e.CreateCluster(r, numberOfEmpires)
	if err != nil {
		return 0, errors.Join(fmt.Errorf("create cluster"), err)
	}

	if forceCreate {
		if err := e.DeleteGame(code); err != nil {
			return 0, errors.Join(fmt.Errorf("force delete game"), err)
		}
	}

	q, tx, err := e.Store.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var homeSystem *System_t
	var homeStar *Star_t
	var homeOrbit *Orbit_t
	var homePlanet *Planet_t

	id, err := q.CreateGame(e.Store.Context, sqlc.CreateGameParams{
		Code:        code,
		Name:        name,
		DisplayName: displayName,
	})
	if err != nil {
		return 0, errors.Join(fmt.Errorf("create game"), err)
	}
	log.Printf("create: game: %d: %s\n", id, code)

	log.Printf("create: systems: %8d systems\n", len(cluster.Systems)-1)
	for _, system := range cluster.Systems {
		if system == nil {
			continue
		}
		if homeSystem == nil {
			homeSystem = system
		}
		systemId, err := q.CreateSystem(e.Store.Context, sqlc.CreateSystemParams{
			GameID:   id,
			X:        system.Coordinates.X,
			Y:        system.Coordinates.Y,
			Z:        system.Coordinates.Z,
			Scarcity: int64(system.Scarcity),
		})
		if err != nil {
			return 0, errors.Join(fmt.Errorf("create system"), err)
		}
		// update the system with the id from the database
		system.Id = systemId
	}
	if homeSystem == nil {
		return 0, errors.New("home system not found")
	}

	log.Printf("create: stars: %8d stars\n", len(cluster.Stars)-1)
	for _, star := range cluster.Stars {
		if star == nil {
			continue
		}
		if homeStar == nil && star.System == homeSystem {
			homeStar = star
		}
		starId, err := q.CreateStar(e.Store.Context, sqlc.CreateStarParams{
			SystemID: star.System.Id,
			Sequence: star.Sequence,
			Scarcity: int64(star.Scarcity),
		})
		if err != nil {
			return 0, errors.Join(fmt.Errorf("create star"), err)
		}
		// update the star with the id from the database
		star.Id = starId
	}
	if homeStar == nil {
		return 0, errors.New("home star not found")
	}

	log.Printf("create: orbits: %8d orbits\n", len(cluster.Orbits)-1)
	for _, orbit := range cluster.Orbits {
		if orbit == nil {
			continue
		}
		if homeOrbit == nil && orbit.Star == homeStar && orbit.Kind == EarthlikePlanet {
			homeOrbit = orbit
		}
		var kind string
		switch orbit.Kind {
		case EmptyOrbit:
			kind = "empty"
		case AsteroidBelt:
			kind = "asteroid-belt"
		case EarthlikePlanet:
			kind = "terrestrial"
		case GasGiant:
			kind = "gas-giant"
		case IceGiant:
			kind = "gas-giant"
		case RockyPlanet:
			kind = "terrestrial"
		default:
			panic(fmt.Sprintf("assert(orbit.kind != %d)", orbit.Kind))
		}
		orbitId, err := q.CreateOrbit(e.Store.Context, sqlc.CreateOrbitParams{
			StarID:   orbit.Star.Id,
			OrbitNo:  orbit.OrbitNo,
			Kind:     kind,
			Scarcity: int64(orbit.Scarcity),
		})
		if err != nil {
			return 0, errors.Join(fmt.Errorf("create orbit"), err)
		}
		// update the orbit with the id from the database
		orbit.Id = orbitId
	}
	if homeOrbit == nil {
		return 0, errors.New("home orbit not found")
	}

	log.Printf("create: planets: %8d planets\n", len(cluster.Planets)-1)
	for _, planet := range cluster.Planets {
		if planet == nil {
			continue
		}
		if homePlanet == nil && planet.Orbit == homeOrbit {
			homePlanet = planet
		}
		var kind string
		switch planet.Kind {
		case NoPlanet:
			kind = "empty"
		case AsteroidBeltPlanet:
			kind = "asteroid-belt"
		case GasGiantPlanet:
			kind = "gas-giant"
		case TerrestrialPlanet:
			kind = "terrestrial"
		default:
			panic(fmt.Sprintf("assert(planet.kind != %d)", planet.Kind))
		}
		planetId, err := q.CreatePlanet(e.Store.Context, sqlc.CreatePlanetParams{
			OrbitID:      planet.Orbit.Id,
			Kind:         kind,
			Scarcity:     int64(planet.Scarcity),
			Habitability: planet.Habitability,
		})
		if err != nil {
			return 0, errors.Join(fmt.Errorf("create planet"), err)
		}
		// update the planet with the id from the database
		planet.Id = planetId
	}
	if homePlanet == nil {
		return 0, errors.New("home planet not found")
	}

	log.Printf("create: deposits: %8d deposits\n", len(cluster.Deposits))
	for _, deposit := range cluster.Deposits {
		if deposit == nil {
			continue
		}
		var resource string
		switch deposit.Resource {
		case NONE:
			resource = "none"
		case FUEL:
			resource = "fuel"
		case GOLD:
			resource = "gold"
		case METALLICS:
			resource = "metallic"
		case NON_METALLICS:
			resource = "non-metallic"
		default:
			panic(fmt.Sprintf("assert(deposit.resource != %d)", deposit.Resource))
		}
		depositId, err := q.CreateDeposit(e.Store.Context, sqlc.CreateDepositParams{
			PlanetID:     deposit.Planet.Id,
			DepositNo:    deposit.DepositNo,
			Kind:         resource,
			YieldPct:     int64(deposit.Yield),
			InitialQty:   int64(deposit.Quantity),
			RemainingQty: int64(deposit.Quantity),
		})
		if err != nil {
			return 0, errors.Join(fmt.Errorf("create deposit"), err)
		}
		// update the deposit with the id from the database
		deposit.Id = depositId
	}

	log.Printf("create: empires: %8d empires\n", len(cluster.Empires)-1)
	for _, empire := range cluster.Empires {
		if empire == nil {
			continue
		}
		parms := sqlc.CreateEmpireParams{
			GameID:       id,
			EmpireNo:     empire.EmpireNo,
			Name:         empire.Name,
			HomeSystemID: empire.HomeSystem.Id,
			HomeStarID:   empire.HomeStar.Id,
			HomeOrbitID:  empire.HomeOrbit.Id,
			HomePlanetID: empire.HomePlanet.Id,
		}
		empireId, err := q.CreateEmpire(e.Store.Context, parms)
		if err != nil {
			return 0, errors.Join(fmt.Errorf("create empire"), err)
		}
		// update the empire with the id from the database
		empire.Id = empireId
	}

	// clean up the deposits table. we added empty deposits to keep players
	// from guessing system resources based on the number of deposits they have
	// seen.
	err = q.DeleteEmptyDeposits(e.Store.Context)
	if err != nil {
		return 0, errors.Join(fmt.Errorf("delete empty deposits"), err)
	}

	// clean up the orbits table. we added empty orbits to keep players from
	// guessing system resources based on the number of orbits they have seen.
	// if constraints are implemented properly, this should also delete the
	// planets and deposits.
	err = q.DeleteEmptyOrbits(e.Store.Context)
	if err != nil {
		return 0, errors.Join(fmt.Errorf("delete empty orbits"), err)
	}

	// calculate the system distances to help reporting
	if calculateSystemDistances {
		for _, from := range cluster.Systems {
			for _, to := range cluster.Systems {
				if from == nil || to == nil {
					continue
				}
				distance := int64(0)
				if from.Id != to.Id {
					dx := from.Coordinates.X - to.Coordinates.X
					dy := from.Coordinates.Y - to.Coordinates.Y
					dz := from.Coordinates.Z - to.Coordinates.Z
					distance = int64(math.Ceil(math.Sqrt(float64(dx*dx + dy*dy + dz*dz))))
				}
				err := q.CreateSystemDistance(e.Store.Context, sqlc.CreateSystemDistanceParams{
					FromSystemID: from.Id,
					ToSystemID:   to.Id,
					Distance:     distance,
				})
				if err != nil {
					return 0, errors.Join(fmt.Errorf("create system distance"), err)
				}
			}
		}
	}

	// update some game meta data
	err = q.UpdateGameEmpireMetadata(e.Store.Context, sqlc.UpdateGameEmpireMetadataParams{
		GameID:       id,
		EmpireNo:     0,
		HomeSystemID: homeSystem.Id,
		HomeStarID:   homeStar.Id,
		HomeOrbitID:  homeOrbit.Id,
		HomePlanetID: homePlanet.Id,
	})
	if err != nil {
		return 0, errors.Join(fmt.Errorf("update game empire metadata"), err)
	}

	return id, tx.Commit()
}

func (e *Engine_t) CreateCluster(r *rand.Rand, numberOfEmpires int64) (*Cluster_t, error) {
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
	cluster := &Cluster_t{
		Systems: []*System_t{nil},
		Stars:   []*Star_t{nil},
		Orbits:  []*Orbit_t{nil},
		Planets: []*Planet_t{nil},
		Empires: []*Empire_t{nil},
	}
	// create 100 systems for the cluster
	isHomeSystem := true
	for i := int64(1); len(points) != 0; i++ {
		var numberOfStars int64
		var scarcity Scarcity_e
		switch i {
		case 1: // 1 4-star system
			numberOfStars, scarcity = 4, TYPICAL
		case 2, 3, 4: // 3 3-star systems
			numberOfStars, scarcity = 3, TYPICAL
		case 5, 6, 7, 8, 9, 10, 11, 12, 13: // 9 2-star systems
			numberOfStars, scarcity = 2, TYPICAL
		default: // remaining are all 1-star systems
			numberOfStars, scarcity = 1, TYPICAL
		}
		// grab the pre-allocated point from the slice of points
		point, points = points[0], points[1:]
		// create the system using the point and number of stars
		system := &System_t{Id: i, Coordinates: Point_t{X: point.X + 15, Y: point.Y + 15, Z: point.Z + 15}}
		// create the stars for the system
		for j := int64(0); j < numberOfStars; j++ {
			star := &Star_t{Id: int64(len(cluster.Stars)), System: system, Sequence: string(rune(65 + j)), Scarcity: scarcity}
			var err error
			var orbits [11]*Orbit_t
			if isHomeSystem {
				orbits, err = generateHomeSystemOrbits(r, star)
				if err != nil {
					return nil, errors.Join(fmt.Errorf("generate home system orbits"), err)
				}
				isHomeSystem = false
			} else {
				orbits, err = generateSystemOrbits(r, star)
				if err != nil {
					return nil, errors.Join(fmt.Errorf("generate system orbits"), err)
				}
			}
			for _, orbit := range orbits {
				if orbit == nil {
					continue
				}
				star.Orbits[orbit.OrbitNo] = orbit
				cluster.Planets = append(cluster.Planets, orbit.Planet)
				cluster.Orbits = append(cluster.Orbits, orbit)
				for _, deposit := range orbit.Planet.Deposits {
					if deposit == nil {
						continue
					}
					cluster.Deposits = append(cluster.Deposits, deposit)
				}
			}
			cluster.Stars = append(cluster.Stars, star)
			system.Stars = append(system.Stars, star)
		}
		cluster.Systems = append(cluster.Systems, system)
	}

	if true {
		log.Printf("create: cluster: skipping empire creation!\n")
	} else {
		for i := int64(1); i <= numberOfEmpires; i++ {
			empire := &Empire_t{
				EmpireNo:   i,
				Name:       fmt.Sprintf("Empire %03d", i),
				HomeSystem: cluster.Systems[1],
				HomeStar:   cluster.Stars[1],
				HomeOrbit:  cluster.Orbits[3],
				HomePlanet: cluster.Planets[3],
			}
			cluster.Empires = append(cluster.Empires, empire)
		}
	}

	if len(points) != 0 {
		// we should have consumed all the points!
		panic(fmt.Sprintf("assert(len(points) != %d)", len(points)))
	}

	return cluster, nil
}

func (e *Engine_t) DeleteGame(code string) error {
	err := e.Store.Queries.DeleteGame(e.Store.Context, code)
	if err != nil {
		return errors.Join(fmt.Errorf("delete game"), err)
	}
	return nil
}

// createPlanet creates a planet.
func createPlanet(r *rand.Rand, orbit *Orbit_t) (*Planet_t, error) {
	planet := &Planet_t{System: orbit.System, Star: orbit.Star, Orbit: orbit, Scarcity: orbit.Scarcity}
	switch orbit.Kind {
	case EmptyOrbit:
		planet.Kind = NoPlanet
	case AsteroidBelt:
		planet.Kind = AsteroidBeltPlanet
	case EarthlikePlanet, RockyPlanet:
		planet.Kind = TerrestrialPlanet
	case GasGiant, IceGiant:
		planet.Kind = GasGiantPlanet
	default:
		panic(fmt.Sprintf("assert(orbit.Kind != %d)", orbit.Kind))
	}
	return planet, nil
}

// createDeposits creates natural resource deposits for a planet.
func createDeposits(r *rand.Rand, planet *Planet_t, isHomeSystem bool) error {
	var numberOfDeposits int64
	if isHomeSystem {
		numberOfDeposits = 35
	} else {
		switch planet.Scarcity {
		case TYPICAL:
			numberOfDeposits = normalRandInRange(r, 1, 35)
		case RICH:
			numberOfDeposits = normalRandInRange(r, 16, 35)
		case POOR:
			numberOfDeposits = normalRandInRange(r, 1, 17)
		}
	}
	for i := int64(1); i <= numberOfDeposits; i++ {
		deposit, err := createDeposit(r, planet)
		if err != nil {
			return fmt.Errorf("deposit: %w", err)
		}
		planet.Deposits[i] = deposit
	}

	// sort deposits by resource type
	// create a temporary slice to hold non-nil deposits (excluding index 0)
	var deposits []*Deposit_t
	for i := int64(1); i < 36; i++ {
		if deposit := planet.Deposits[i]; deposit != nil {
			deposits = append(deposits, deposit)
		}
	}
	// sort the slice by resource type. assumes that we have not yet inserted any NONE deposits!
	sort.Slice(deposits, func(i, j int) bool {
		if deposits[i].Resource == deposits[j].Resource {
			if deposits[i].Quantity == deposits[j].Quantity {
				return deposits[i].Yield > deposits[j].Yield
			}
			return deposits[i].Quantity > deposits[j].Quantity
		}
		return deposits[i].Resource < deposits[j].Resource
	})
	// number the deposits and put them back into the planet's deposits
	for n, deposit := range deposits {
		deposit.DepositNo = int64(n + 1)
		planet.Deposits[deposit.DepositNo] = deposit
	}
	// fill in the remaining deposits with NONE
	for i := int64(1); i <= 35; i++ {
		if planet.Deposits[i] == nil {
			planet.Deposits[i] = &Deposit_t{Planet: planet, Resource: NONE}
		}
		planet.Deposits[i].DepositNo = i
	}
	return nil
}

// createDeposit creates a natural resource deposit.
func createDeposit(r *rand.Rand, planet *Planet_t) (deposit *Deposit_t, err error) {
	var resource Resource_e
	switch planet.Scarcity {
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
	case NONE: // should never happen
		return &Deposit_t{Planet: planet, Resource: resource}, nil
	case GOLD:
		switch planet.Scarcity {
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
		switch planet.Scarcity {
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
		switch planet.Scarcity {
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
		switch planet.Scarcity {
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
		Planet:   planet,
		Resource: resource,
		Quantity: normalRandInRange(r, minQuantity, maxQuantity),
		Yield:    normalRandInRange(r, minYield, maxYield),
	}, nil
}

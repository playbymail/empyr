// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"errors"
	"fmt"
	"github.com/playbymail/empyr/store"
	"github.com/playbymail/empyr/store/sqlc"
	"log"
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

func (e *Engine_t) CreateGame(code, name, displayName string, includeEmptyResources, calculateSystemDistances bool, r *rand.Rand, forceCreate bool) (*Game_t, error) {
	g := &Game_t{
		Systems:  make(map[int64]*System_t),
		Stars:    make(map[int64]*Star_t),
		Orbits:   make(map[int64]*Orbit_t),
		Planets:  make(map[int64]*Planet_t),
		Deposits: make(map[int64]*Deposit_t),
		Empires:  make(map[int64]*Empire_t),
	}

	// create an empty cluster with about 100 systems and use templates to populate the cluster
	cluster := emptyCluster(r)
	for _, system := range cluster.Systems {
		if system.Id == 1 {
			homeSystemTemplate(system, r)
		} else {
			coreSystemTemplate(system, r)
		}
	}
	// copy the cluster to the game
	for _, system := range cluster.Systems {
		g.Systems[system.Id] = system
		for _, star := range system.Stars {
			g.Stars[star.Id] = star
			for _, orbit := range star.Orbits {
				if orbit != nil && (orbit.Kind != EmptyOrbit || includeEmptyResources) {
					g.Orbits[orbit.Id] = orbit
					planet := orbit.Planet
					if planet != nil {
						g.Planets[planet.Id] = planet
						for _, deposit := range planet.Deposits {
							if deposit != nil {
								if deposit.Resource != NONE || includeEmptyResources {
									g.Deposits[deposit.Id] = deposit
								}
							}
						}
					}
				}
			}
		}
	}
	g.Home.System = g.Systems[1]
	g.Home.Star = g.Home.System.Stars[0]
	g.Home.Orbit = g.Home.Star.Orbits[3]
	g.Home.Planet = g.Home.Orbit.Planet

	if forceCreate {
		if err := e.DeleteGame(code); err != nil {
			return nil, errors.Join(fmt.Errorf("force delete game"), err)
		}
	}

	q, tx, err := e.Store.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	id, err := q.CreateGame(e.Store.Context, sqlc.CreateGameParams{
		Code:        code,
		Name:        name,
		DisplayName: displayName,
	})
	if err != nil {
		return nil, errors.Join(fmt.Errorf("create game"), err)
	}
	g.Id = id
	log.Printf("create: game: %d: %s\n", g.Id, code)

	clusterID, err := q.CreateCluster(e.Store.Context, sqlc.CreateClusterParams{
		GameID:       id,
		HomeSystemID: g.Home.System.Id,
		HomeStarID:   g.Home.Star.Id,
		HomeOrbitID:  g.Home.Orbit.Id,
		HomePlanetID: g.Home.Planet.Id,
	})
	if err != nil {
		return nil, errors.Join(fmt.Errorf("create cluster"), err)
	}

	var systems []*System_t
	for _, system := range g.Systems {
		if system != nil {
			systems = append(systems, system)
		}
	}
	sort.Slice(systems, func(i, j int) bool {
		return systems[i].Id < systems[j].Id
	})
	var stars []*Star_t
	for _, star := range g.Stars {
		if star != nil {
			stars = append(stars, star)
		}
	}
	sort.Slice(stars, func(i, j int) bool {
		return stars[i].Id < stars[j].Id
	})
	var orbits []*Orbit_t
	for _, orbit := range g.Orbits {
		if orbit != nil {
			orbits = append(orbits, orbit)
		}
	}
	sort.Slice(orbits, func(i, j int) bool {
		return orbits[i].Id < orbits[j].Id
	})
	var planets []*Planet_t
	for _, planet := range g.Planets {
		if planet != nil {
			planets = append(planets, planet)
		}
	}
	sort.Slice(planets, func(i, j int) bool {
		return planets[i].Id < planets[j].Id
	})
	var deposits []*Deposit_t
	for _, deposit := range g.Deposits {
		if deposit != nil {
			deposits = append(deposits, deposit)
		}
	}
	sort.Slice(deposits, func(i, j int) bool {
		return deposits[i].Id < deposits[j].Id
	})

	log.Printf("create: systems: %8d systems\n", len(g.Systems))
	g.Systems = make(map[int64]*System_t)
	for _, system := range systems {
		if system == nil {
			continue
		}
		systemId, err := q.CreateSystem(e.Store.Context, sqlc.CreateSystemParams{
			ClusterID: clusterID,
			X:         system.Coordinates.X,
			Y:         system.Coordinates.Y,
			Z:         system.Coordinates.Z,
		})
		if err != nil {
			return nil, errors.Join(fmt.Errorf("create system"), err)
		}
		// update the system with the id from the database
		system.Id = systemId
		// update the id numbers in the game map
		g.Systems[system.Id] = system
	}

	log.Printf("create: stars: %8d stars\n", len(g.Stars))
	g.Stars = make(map[int64]*Star_t)
	for _, star := range stars {
		if star == nil {
			continue
		}
		starId, err := q.CreateStar(e.Store.Context, sqlc.CreateStarParams{
			SystemID: star.System.Id,
			Sequence: star.Sequence,
		})
		if err != nil {
			return nil, errors.Join(fmt.Errorf("create star"), err)
		}
		// update the star with the id from the database
		star.Id = starId
		// update the id numbers in the game map
		g.Stars[star.Id] = star
	}

	log.Printf("create: orbits: %8d orbits\n", len(g.Orbits))
	g.Orbits = make(map[int64]*Orbit_t)
	for _, orbit := range orbits {
		if orbit == nil {
			continue
		}
		orbitId, err := q.CreateOrbit(e.Store.Context, sqlc.CreateOrbitParams{
			StarID:  orbit.Star.Id,
			OrbitNo: orbit.OrbitNo,
			Kind:    orbit.Kind.Code(),
		})
		if err != nil {
			return nil, errors.Join(fmt.Errorf("create orbit"), err)
		}
		// update the orbit with the id from the database
		orbit.Id = orbitId
		// update the id numbers in the game map
		g.Orbits[orbit.Id] = orbit
	}

	log.Printf("create: planets: %8d planets\n", len(g.Planets))
	g.Planets = make(map[int64]*Planet_t)
	for _, planet := range planets {
		if planet == nil {
			continue
		}
		planetId, err := q.CreatePlanet(e.Store.Context, sqlc.CreatePlanetParams{
			OrbitID:      planet.Orbit.Id,
			Kind:         planet.Kind.Code(),
			Habitability: planet.Habitability,
		})
		if err != nil {
			return nil, errors.Join(fmt.Errorf("create planet"), err)
		}
		// update the planet with the id from the database
		planet.Id = planetId
		// update the id numbers in the game map
		g.Planets[planet.Id] = planet
	}

	log.Printf("create: deposits: %8d deposits\n", len(g.Deposits))
	g.Deposits = make(map[int64]*Deposit_t)
	for _, deposit := range deposits {
		if deposit == nil {
			continue
		}
		depositParams := sqlc.CreateDepositParams{
			PlanetID:     deposit.Planet.Id,
			DepositNo:    deposit.DepositNo,
			Kind:         deposit.Resource.Code(),
			YieldPct:     deposit.Yield,
			RemainingQty: deposit.Quantity,
		}
		depositId, err := q.CreateDeposit(e.Store.Context, depositParams)
		if err != nil {
			log.Printf("engine: createGame: create deposit: %+v\n", depositParams)
			log.Printf("engine: createGame: create deposit: %v\n", err)
			return nil, errors.Join(fmt.Errorf("create deposit"), err)
		}
		// update the deposit with the id from the database
		deposit.Id = depositId
		// update the id numbers in the game map
		g.Deposits[deposit.Id] = deposit
	}

	log.Printf("create: empires: %8d empires\n", len(g.Empires))
	for _, empire := range g.Empires {
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
			return nil, errors.Join(fmt.Errorf("create empire"), err)
		}
		// update the empire with the id from the database
		empire.Id = empireId
	}
	// update the id numbers in the game map
	if src := g.Empires; len(src) != 0 {
		g.Empires = make(map[int64]*Empire_t)
		for _, empire := range src {
			g.Empires[empire.Id] = empire
		}
	}

	// calculate the system distances to help reporting
	if calculateSystemDistances {
		err = q.PopulateSystemDistanceByCluster(e.Store.Context, clusterID)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("populate system distance"), err)
		}
	}

	return g, tx.Commit()
}

func (e *Engine_t) ObsoleteCreateCluster(r *rand.Rand, numberOfEmpires int64) (*Cluster_t, error) {
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
		switch i {
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
		system := &System_t{Id: i, Coordinates: Point_t{X: point.X + 15, Y: point.Y + 15, Z: point.Z + 15}}
		// create the stars for the system
		for j := int64(0); j < numberOfStars; j++ {
			star := &Star_t{Id: int64(len(cluster.Stars)), System: system, Sequence: string(rune(65 + j))}
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
	err := e.Store.Queries.DeleteGameByCode(e.Store.Context, code)
	if err != nil {
		return errors.Join(fmt.Errorf("delete game"), err)
	}
	return nil
}

// createPlanet creates a planet.
func createPlanet(r *rand.Rand, orbit *Orbit_t) (*Planet_t, error) {
	planet := &Planet_t{System: orbit.System, Star: orbit.Star, Orbit: orbit}
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
		numberOfDeposits = normalRandInRange(r, 1, 35)
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
	var minQuantity, maxQuantity, minYield, maxYield int
	if n := r.IntN(100); n < 3 {
		resource = GOLD
		minQuantity, maxQuantity = 1_000_000, 35_000_000
		minYield, maxYield = 1, 9
	} else if n < 22 {
		resource = FUEL
		minQuantity, maxQuantity = 1_000_000, 50_000_000
		minYield, maxYield = 1, 12
	} else if n < 56 {
		resource = METALLICS
		minQuantity, maxQuantity = 1_000_000, 99_000_000
		minYield, maxYield = 1, 36
	} else {
		resource = NON_METALLICS
		minQuantity, maxQuantity = 1_000_000, 99_000_000
		minYield, maxYield = 1, 36
	}
	return &Deposit_t{
		Planet:   planet,
		Resource: resource,
		Quantity: normalRandInRange(r, minQuantity, maxQuantity),
		Yield:    normalRandInRange(r, minYield, maxYield),
	}, nil
}

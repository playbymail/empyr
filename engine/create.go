// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"errors"
	"fmt"
	"github.com/playbymail/empyr/store"
	"github.com/playbymail/empyr/store/sqlc"
	"math"
	"math/rand/v2"
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

func (e *Engine_t) CreateGame(code, name, displayName string, r *rand.Rand) (int64, error) {
	cluster, err := e.CreateCluster(r)
	if err != nil {
		return 0, errors.Join(fmt.Errorf("create cluster"), err)
	}
	//for _, obj := range []struct {
	//	name string
	//	ptr  any
	//}{
	//	{name: "systems", ptr: cluster.Systems[1:]},
	//	{name: "stars", ptr: cluster.Stars[1:]},
	//	{name: "orbits", ptr: cluster.Orbits[1:]},
	//	{name: "planets", ptr: cluster.Planets[1:]},
	//} {
	//	if data, err := json.MarshalIndent(obj.ptr, "", "  "); err != nil {
	//		return 0, err
	//	} else if err = os.WriteFile(obj.name+".json", data, 0644); err != nil {
	//		log.Fatalf("cluster: %s: %v\n", obj.name, err)
	//	} else {
	//		log.Printf("cluster: %s: wrote %s\n", obj.name, obj.name+".json")
	//	}
	//}

	q, tx, err := e.Store.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	err = q.DeleteGame(e.Store.Context, code)
	if err != nil {
		return 0, errors.Join(fmt.Errorf("delete game"), err)
	}

	id, err := q.CreateGame(e.Store.Context, sqlc.CreateGameParams{
		Code:        code,
		Name:        name,
		DisplayName: displayName,
	})
	if err != nil {
		return 0, errors.Join(fmt.Errorf("create game"), err)
	}

	dbSystemIds := map[int64]int64{}
	for _, system := range cluster.Systems {
		if system == nil {
			continue
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
		dbSystemIds[system.Id] = systemId
	}

	dbStarIds := map[int64]int64{}
	for _, star := range cluster.Stars {
		if star == nil {
			continue
		}
		starId, err := q.CreateStar(e.Store.Context, sqlc.CreateStarParams{
			SystemID: dbSystemIds[star.SystemId],
			Sequence: star.Sequence,
			Scarcity: int64(star.Scarcity),
		})
		if err != nil {
			return 0, errors.Join(fmt.Errorf("create star"), err)
		}
		dbStarIds[star.Id] = starId
	}

	for _, star := range cluster.Stars {
		if star == nil {
			continue
		}
		err := q.CreateSystemStarLink(e.Store.Context, sqlc.CreateSystemStarLinkParams{
			SystemID: dbSystemIds[star.SystemId],
			StarID:   dbStarIds[star.Id],
		})
		if err != nil {
			return 0, errors.Join(fmt.Errorf("link system with star"), err)
		}
	}

	dbOrbitIds := map[int64]int64{}
	for _, orbit := range cluster.Orbits {
		if orbit == nil {
			continue
		}
		var kind string
		switch orbit.Kind {
		case EmptyOrbit:
			kind = "empty"
		case AsteroidBelt:
			kind = "asteroid-belt"
		case EarthlikePlant:
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
			StarID:  dbStarIds[orbit.StarId],
			OrbitNo: orbit.OrbitNo,
			Kind:    kind,
		})
		if err != nil {
			return 0, errors.Join(fmt.Errorf("create orbit"), err)
		}
		dbOrbitIds[orbit.Id] = orbitId
	}

	for _, orbit := range cluster.Orbits {
		if orbit == nil {
			continue
		}
		err := q.CreateStarOrbitLink(e.Store.Context, sqlc.CreateStarOrbitLinkParams{
			StarID:  dbStarIds[orbit.StarId],
			OrbitID: dbOrbitIds[orbit.Id],
		})
		if err != nil {
			return 0, errors.Join(fmt.Errorf("link star with orbit"), err)
		}
	}

	// clean up the orbits table. we added empty orbits to keep players from
	// guessing system resources based on the number of orbits they have seen.
	// if constraints are implemented properly, this should also delete the planets.
	err = q.DeleteEmptyOrbits(e.Store.Context)
	if err != nil {
		return 0, errors.Join(fmt.Errorf("delete empty orbits"), err)
	}

	// calculate the system distances to help reporting
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
				FromSystemID: dbSystemIds[from.Id],
				ToSystemID:   dbSystemIds[to.Id],
				Distance:     distance,
			})
			if err != nil {
				return 0, errors.Join(fmt.Errorf("create system distance"), err)
			}
		}
	}

	return id, tx.Commit()
}

func (e *Engine_t) CreateCluster(r *rand.Rand) (*Cluster_t, error) {
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
	}
	// create 100 systems for the cluster
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
			star := &Star_t{Id: int64(len(cluster.Stars)), SystemId: system.Id, Sequence: string(rune(65 + j)), Scarcity: scarcity}
			// all stars have 10 orbits but not all orbits have planets
			numberOfPlanets := r.IntN(5) + r.IntN(6) + 1 // normalRandInRange(r, 1, 10)
			// generate the rings for the star based on the number of planets
			rings := generateRings(r, numberOfPlanets)
			// generate the planet for each orbit
			for k := int64(1); k <= 10; k++ {
				orbit := &Orbit_t{Id: int64(len(cluster.Orbits)), StarId: star.Id, OrbitNo: k, Kind: rings[k]}
				cluster.Orbits = append(cluster.Orbits, orbit)
				planet := &Planet_t{Id: orbit.Id, Star: star.Id}
				switch orbit.Kind {
				case EmptyOrbit:
					planet.Kind = NoPlanet
				case AsteroidBelt:
					planet.Kind = AsteroidBeltPlanet
				case EarthlikePlant, RockyPlanet:
					planet.Kind = TerrestrialPlanet
				case GasGiant, IceGiant:
					planet.Kind = GasGiantPlanet
				default:
					panic(fmt.Sprintf("assert(orbit.Kind != %d)", orbit.Kind))
				}
				//orbit.Planet, err = createPlanet(r, rings[k], scarcity)
				//if err != nil {
				//	return nil, fmt.Errorf("planet: %w", err)
				//}
				cluster.Planets = append(cluster.Planets, planet)
			}
			cluster.Stars = append(cluster.Stars, star)
			system.Stars = append(system.Stars, star.Id)
		}
		cluster.Systems = append(cluster.Systems, system)
	}

	if len(points) != 0 {
		// we should have consumed all the points!
		panic(fmt.Sprintf("assert(len(points) != %d)", len(points)))
	}

	return cluster, nil
}

// createDeposits creates natural resource deposits for a planet.
func createDeposits(r *rand.Rand, scarcity Scarcity_e) ([]*Deposit_t, error) {
	var numberOfDeposits int64
	switch scarcity {
	case TYPICAL:
		numberOfDeposits = normalRandInRange(r, 1, 35)
	case RICH:
		numberOfDeposits = normalRandInRange(r, 16, 35)
	case POOR:
		numberOfDeposits = normalRandInRange(r, 1, 17)
	}
	var deposits []*Deposit_t
	for i := int64(0); i < numberOfDeposits; i++ {
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

// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"math/rand/v2"
	"sort"
)

// emptyCluster creates a cluster using 100 random locations for the systems.
// semi-random because the first system is always at the origin and the next
// three are with one unit of the origin.
//
// the systems will have 1 to 4 stars each. each star will have all orbits populated,
// all orbits will have planets populated, and all planets will have deposits populated.
//
// NB: it's not obvious, but we fully load every system, star, orbit, and planet
// to prevent players from using the id values to guess the location of resources.
func emptyCluster(r *rand.Rand) *Cluster_t {
	var c Cluster_t

	var starsPerSystem = []int{4, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2, 2}

	// create the systems, offset the points to move the origin to (15, 15, 15)
	for _, point := range randomPoints(r) {
		numberOfStarsInSystem := 1 // default to at least one star
		if len(starsPerSystem) != 0 {
			// if we have more stars per system, use and consume the next value
			numberOfStarsInSystem = starsPerSystem[0]
			starsPerSystem = starsPerSystem[1:]
		}
		c.Systems = append(c.Systems, emptySystem(Point_t{X: point.X + 15, Y: point.Y + 15, Z: point.Z + 15}, numberOfStarsInSystem))
	}

	// add id and links to the stars, orbits, planets, and deposits
	systemId, starId, orbitId, planetId, depositId := int64(1), int64(1), int64(1), int64(1), int64(1)
	for _, system := range c.Systems {
		if system != nil {
			system.Id, systemId = systemId, systemId+1
			for _, star := range system.Stars {
				if star != nil {
					c.Stars = append(c.Stars, star)
					star.Id, starId = starId, starId+1
					for ono, orbit := range star.Orbits {
						if orbit != nil {
							c.Orbits = append(c.Orbits, orbit)
							orbit.Id, orbit.OrbitNo, orbitId = orbitId, int64(ono), orbitId+1
							planet := orbit.Planet
							if planet != nil {
								c.Planets = append(c.Planets, planet)
								planet.Id, planetId = planetId, planetId+1
								for n, deposit := range planet.Deposits {
									if deposit != nil {
										c.Deposits = append(c.Deposits, deposit)
										deposit.Id, deposit.DepositNo, depositId = depositId, int64(n), depositId+1
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return &c
}

func emptySystem(point Point_t, numberOfStars int) *System_t {
	system := &System_t{Coordinates: point}
	for i := 0; i < numberOfStars; i++ {
		system.Stars = append(system.Stars, emptyStar(system, i+1))
	}
	return system
}

// emptyStar creates an empty star for a system.
func emptyStar(system *System_t, sequence int) *Star_t {
	star := &Star_t{System: system, Sequence: string(rune('A' + sequence - 1))}
	for j := 1; j <= 10; j++ {
		star.Orbits[j] = emptyOrbit(star)
	}
	return star
}

// emptyOrbit creates an empty orbit for a star.
func emptyOrbit(star *Star_t) *Orbit_t {
	orbit := &Orbit_t{System: star.System, Star: star}
	for j := 1; j <= 10; j++ {
		orbit.Planet = emptyPlanet(orbit)
	}
	return orbit
}

// emptyPlanet creates an empty planet for an orbit.
func emptyPlanet(orbit *Orbit_t) *Planet_t {
	planet := &Planet_t{System: orbit.System, Star: orbit.Star, Orbit: orbit}
	for k := 1; k <= 35; k++ {
		planet.Deposits[k] = emptyDeposit(planet, int64(k))
	}
	return planet
}

// emptyDeposit creates an empty deposit for a planet.
func emptyDeposit(p *Planet_t, no int64) *Deposit_t {
	return &Deposit_t{DepositNo: no, Planet: p}
}

// homeSystemTemplate updates a system using the home system template.
func homeSystemTemplate(system *System_t, r *rand.Rand) {
	for _, star := range system.Stars {
		if star != nil {
			if star.Id == 1 {
				homeStarTemplate(star, r)
			} else {
				coreStarTemplate(star, r)
			}
		}
	}
}

func coreSystemTemplate(system *System_t, r *rand.Rand) {
	for _, star := range system.Stars {
		if star != nil {
			coreStarTemplate(star, r)
		}
	}
}

// homeStarTemplate updates a star using the home star template.
func homeStarTemplate(star *Star_t, r *rand.Rand) {
	for _, entry := range []struct {
		No           int
		OrbitKind    Orbit_e
		PlanetKind   Planet_e
		Habitability int64
	}{
		{No: 1, OrbitKind: RockyPlanet, PlanetKind: TerrestrialPlanet},
		{No: 2, OrbitKind: EarthlikePlanet, PlanetKind: TerrestrialPlanet, Habitability: 3},
		{No: 3, OrbitKind: EarthlikePlanet, PlanetKind: TerrestrialPlanet, Habitability: 25},
		{No: 4, OrbitKind: EarthlikePlanet, PlanetKind: TerrestrialPlanet, Habitability: 8},
		{No: 5, OrbitKind: AsteroidBelt, PlanetKind: AsteroidBeltPlanet},
		{No: 6, OrbitKind: GasGiant, PlanetKind: GasGiantPlanet},
		{No: 7, OrbitKind: IceGiant, PlanetKind: GasGiantPlanet},
		{No: 8, OrbitKind: IceGiant, PlanetKind: GasGiantPlanet},
		{No: 10, OrbitKind: AsteroidBelt, PlanetKind: AsteroidBeltPlanet},
	} {
		orbit := star.Orbits[entry.No]
		orbit.Kind = entry.OrbitKind
		planet := orbit.Planet
		planet.Kind = entry.PlanetKind
		planet.Habitability = entry.Habitability
	}

	for _, orbit := range star.Orbits {
		if orbit != nil {
			homePlanetTemplate(orbit.Planet, r)
		}
	}
}

func coreStarTemplate(star *Star_t, r *rand.Rand) {
	rings := []Orbit_e{EarthlikePlanet, RockyPlanet, RockyPlanet, AsteroidBelt, AsteroidBelt, GasGiant, GasGiant, IceGiant, IceGiant}
	earthlikeCounter, gasGiantCounter := 0, 0
	for no := 1; no <= 10; no++ {
		orbit := star.Orbits[no]
		kind := EmptyOrbit
		// 30% chance of a planet in this orbit
		if r.IntN(100) < 30 {
			kind = rings[r.IntN(len(rings))]
			if kind == EarthlikePlanet {
				earthlikeCounter++
				if earthlikeCounter > 2 || no > 5 {
					// can't have more than two habitable planets or any past the 5th orbit
					kind = RockyPlanet
				}
			} else if kind == GasGiant {
				gasGiantCounter++
				if gasGiantCounter > 2 {
					// can't have more than two gas giants
					kind = RockyPlanet
				}
			}
		}
		orbit.Kind = kind
		switch kind {
		case EmptyOrbit:
			orbit.Planet.Kind = NoPlanet
		case AsteroidBelt:
			orbit.Planet.Kind = AsteroidBeltPlanet
		case EarthlikePlanet:
			orbit.Planet.Kind = TerrestrialPlanet
			orbit.Planet.Habitability = roll(r, 4, 6) - 3 // 4d6 - 3
		case GasGiant:
			orbit.Planet.Kind = GasGiantPlanet
		case IceGiant:
			orbit.Planet.Kind = GasGiantPlanet
		case RockyPlanet:
			orbit.Planet.Kind = TerrestrialPlanet
		}
	}

	for _, orbit := range star.Orbits {
		if orbit != nil {
			corePlanetTemplate(orbit.Planet, r)
		}
	}
}

func homePlanetTemplate(planet *Planet_t, r *rand.Rand) {
	goldDeposits, fuelDeposits, isHabitable := 0, 0, planet.Habitability != 0
	for _, deposit := range planet.Deposits {
		if deposit == nil {
			continue
		}
		isHomePlanet := isHabitable && planet.Orbit.OrbitNo == 3
		kind, qty, yield := deposit.Resource, deposit.Quantity, deposit.Yield
		switch planet.Kind {
		case NoPlanet:
			kind, qty, yield = NONE, 0, 0
		case AsteroidBeltPlanet:
			if roll(r, 1, 100) == 5 {
				kind, qty, yield = GOLD, normalRandInRange(r, 100_000, 5_000_000), roll(r, 1, 3)
			} else if roll(r, 1, 100) < 15 {
				kind, qty, yield = FUEL, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 3, 6)-2
			} else if roll(r, 1, 100) < 45 {
				kind, qty, yield = METALLICS, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 3, 10)-2
			} else if roll(r, 1, 100) < 45 {
				kind, qty, yield = NON_METALLICS, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 3, 10)-2
			}
		case GasGiantPlanet:
			if roll(r, 1, 100) < 15 && fuelDeposits < 2 {
				kind, qty, yield = FUEL, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 10, 4)-2
			} else if roll(r, 1, 100) < 35 {
				kind, qty, yield = METALLICS, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 10, 6)
			} else if roll(r, 1, 100) < 35 {
				kind, qty, yield = NON_METALLICS, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 10, 6)
			}
		case TerrestrialPlanet:
			if isHomePlanet {
				if deposit.DepositNo == 1 {
					kind, qty, yield = GOLD, normalRandInRange(r, 8_000_000, 20_000_000), roll(r, 3, 4)-3
				} else if deposit.DepositNo <= 4 {
					kind, qty, yield = FUEL, normalRandInRange(r, 35_000_000, 99_000_000), roll(r, 11, 6)
				} else if deposit.DepositNo <= 19 {
					kind, qty, yield = METALLICS, normalRandInRange(r, 10_000_000, 99_000_000), roll(r, 11, 9)
				} else {
					kind, qty, yield = NON_METALLICS, normalRandInRange(r, 10_000_000, 99_000_000), roll(r, 11, 8)
				}
			} else if roll(r, 1, 100) <= 2 && goldDeposits < 3 {
				kind, qty, yield = GOLD, normalRandInRange(r, 100_000, 1_000_000), roll(r, 1, 3)
				if isHabitable {
					yield = roll(r, 3, 4) - 3
				}
			} else if roll(r, 1, 100) < 15 && fuelDeposits < 5 {
				kind, qty, yield = FUEL, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 10, 4)-2
				if isHabitable {
					yield = roll(r, 10, 6)
				}
			} else if roll(r, 1, 100) < 25 {
				kind, qty, yield = METALLICS, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 10, 6)
				if isHabitable {
					yield = roll(r, 10, 8)
				}
			} else if roll(r, 1, 100) < 25 {
				kind, qty, yield = NON_METALLICS, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 10, 6)
				if isHabitable {
					yield = roll(r, 10, 8)
				}
			}
		}
		if kind == GOLD {
			goldDeposits++
		} else if kind == FUEL {
			fuelDeposits++
		}
		deposit.Resource, deposit.Quantity, deposit.Yield = kind, qty, yield
	}
	sortDeposits(planet)
}

func corePlanetTemplate(planet *Planet_t, r *rand.Rand) {
	goldDeposits, fuelDeposits, isHabitable := 0, 0, planet.Habitability != 0
	for _, deposit := range planet.Deposits {
		if deposit == nil {
			continue
		}
		kind, qty, yield := deposit.Resource, deposit.Quantity, deposit.Yield
		switch planet.Kind {
		case NoPlanet:
			kind, qty, yield = NONE, 0, 0
		case AsteroidBeltPlanet:
			if roll(r, 1, 100) == 3 {
				kind, qty, yield = GOLD, normalRandInRange(r, 100_000, 5_000_000), roll(r, 1, 3)
			} else if roll(r, 1, 100) < 10 {
				kind, qty, yield = FUEL, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 3, 6)-2
			} else if roll(r, 1, 100) < 35 {
				kind, qty, yield = METALLICS, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 3, 10)-2
			} else if roll(r, 1, 100) < 35 {
				kind, qty, yield = NON_METALLICS, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 3, 10)-2
			}
		case GasGiantPlanet:
			if roll(r, 1, 100) < 15 && fuelDeposits < 2 {
				kind, qty, yield = FUEL, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 10, 4)-2
			} else if roll(r, 1, 100) < 25 {
				kind, qty, yield = METALLICS, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 10, 6)
			} else if roll(r, 1, 100) < 25 {
				kind, qty, yield = NON_METALLICS, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 10, 6)
			}
		case TerrestrialPlanet:
			if roll(r, 1, 100) == 1 && goldDeposits == 0 {
				kind, qty, yield = GOLD, normalRandInRange(r, 100_000, 1_000_000), roll(r, 1, 3)
				if isHabitable {
					yield = roll(r, 3, 4) - 3
				}
			} else if roll(r, 1, 100) < 15 && fuelDeposits < 5 {
				kind, qty, yield = FUEL, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 10, 4)-2
				if isHabitable {
					yield = roll(r, 10, 6)
				}
			} else if roll(r, 1, 100) < 30 {
				kind, qty, yield = METALLICS, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 10, 6)
				if isHabitable {
					yield = roll(r, 10, 8)
				}
			} else if roll(r, 1, 100) < 30 {
				kind, qty, yield = NON_METALLICS, normalRandInRange(r, 1_000_000, 99_000_000), roll(r, 10, 6)
				if isHabitable {
					yield = roll(r, 10, 8)
				}
			}
		}
		if kind == GOLD {
			goldDeposits++
		} else if kind == FUEL {
			fuelDeposits++
		}
		deposit.Resource, deposit.Quantity, deposit.Yield = kind, qty, yield
	}
	sortDeposits(planet)
}

// sort deposits by resource type. puts GOLD first and NONE last.
func sortDeposits(planet *Planet_t) {
	// create a temporary slice to hold non-nil deposits (excluding index 0)
	var deposits []*Deposit_t
	for _, deposit := range planet.Deposits {
		if deposit != nil {
			deposits = append(deposits, deposit)
		}
	}
	// sort the slice by resource type.
	sort.Slice(deposits, func(i, j int) bool {
		if deposits[i].Resource == deposits[j].Resource {
			if deposits[i].Quantity == deposits[j].Quantity {
				return deposits[i].Yield > deposits[j].Yield
			}
			return deposits[i].Quantity > deposits[j].Quantity
		}
		// GOLD is first, NONE is last
		if deposits[i].Resource == NONE {
			return false
		} else if deposits[j].Resource == NONE {
			return true
		}
		return deposits[i].Resource < deposits[j].Resource
	})
	// number the deposits and put them back into the planet's deposits
	for n, deposit := range deposits {
		deposit.DepositNo = int64(n + 1)
		planet.Deposits[deposit.DepositNo] = deposit
	}
}

// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package oldgame implements state for the game.
package oldgame

import (
	"fmt"
	"math/rand"
	"sync"
)

func Create(id string, users map[string]*User) *Game {
	return &Game{id: id}
}

// Game implements the state for a single game.
type Game struct {
	sync.Mutex
	id      string             // unique identifier for the game.
	users   map[string]*User   // map of all users; to avoid having to call out? todo: bad.
	players map[string]*Player // map of all players currently in the game
	// map of all objects currently in the game.
	// you are strongly encouraged to verify the type of the object before using it!
	objects map[string]*OBJECT
	audit   []Order
	seqno   int
}

func (g *Game) uuidgen() string {
	g.seqno++
	return fmt.Sprintf("%09d", g.seqno)
}

// Order is the only way to change the state of the system
type Order interface {
	Eval(g *Game) error
}

// Eval executes an order.
// Input is always game state and an order.
// Output is updated game state.
func Eval(g *Game, o Order) error {
	g.Lock()
	defer g.Unlock()
	// save a copy of the order for auditing
	g.audit = append(g.audit, o)
	return o.Eval(g)
}

// this order initializes a game
type InitializeGameOrder struct {
	users map[string]*User
}

func randFloat(min, max float64) float64 {
	return rand.Float64()*(max-min+1) + min
}

func randInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func (o InitializeGameOrder) Eval(g *Game) error {
	g.players = make(map[string]*Player)
	g.objects = make(map[string]*OBJECT)
	g.users = o.users

	// create a new star cluster
	var cluster StarCluster
	cluster.id = g.uuidgen()
	// create solar systems within the cluster
	for i := 0; i < 100; i++ {
		system := &SolarSystem{id: g.uuidgen()}
		// coordinates are in the range (0..30)
		system.coordinates.x = randFloat(0, 30)
		system.coordinates.y = randFloat(0, 30)
		system.coordinates.z = randFloat(0, 30)
		// system will have from one to ten planets
		numberOfPlanets := randInt(1, 10)
		// randomly assign the planet to orbits
		var orbits []*Orbit
		for p := 0; p < 10; p++ {
			orbits = append(orbits, &Orbit{id: g.uuidgen()})
			if p < numberOfPlanets {
				planet := &Planet{
					id:           g.uuidgen(),
					habitability: randInt(0, 25),
				}
				switch rand.Intn(3) {
				case 0:
					planet.kind = "GAS GIANT"
				case 1:
					planet.kind = "TERRESTRIAL"
				case 2:
					planet.kind = "ASTEROIDS"
				default:
					panic("assert(rand.Intn(3) in [0,3))")
				}
				numberOfNaturalResources := randInt(0, 40)
				for r := 0; r <= numberOfNaturalResources; r++ {
					nr := &NaturalResource{
						id:            g.uuidgen(),
						initialAmount: randInt(1000000, 99000000),
					}
					nr.amountRemaining = nr.initialAmount
					switch rand.Intn(4) {
					case 0:
						nr.kind = "GOLD"
					case 1:
						nr.kind = "FUEL"
					case 2:
						nr.kind = "METALLIC"
					case 3:
						nr.kind = "NONMETALLIC"
					}
					planet.deposits = append(planet.deposits, nr)
				}
				orbits[p].planet = planet
			}
		}
		// shuffle the orbit of the planets in the system
		rand.Shuffle(len(orbits), func(i, j int) { orbits[i], orbits[j] = orbits[j], orbits[i] })
		cluster.systems[i] = system
	}

	return nil
}

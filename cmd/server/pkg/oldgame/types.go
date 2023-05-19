// empyr - a reimagining of Vern Holford's Empyrean Challenge
// Copyright (C) 2023 Michael D Henderson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package oldgame

import "fmt"

// OBJECT hahah
type OBJECT struct {
	kind    string
	id      string // the id of the contained object
	g       *Game
	factory *Factory
	nation  *Nation
	orbit   *Orbit
	//order *Order
	player      *Player
	starCluster *StarCluster
	system      *SolarSystem
	user        *User
}

type StarCluster struct {
	id      string
	systems [100]*SolarSystem
}

type Factory struct {
	id    string // unique identifier for the factory
	owner string // player who owns/controls this factory
}

type Nation struct {
	id        string
	homeWorld string
}

type NaturalResource struct {
	id              string
	kind            string // gold, fuel, metallics, non-metallics
	initialAmount   int
	amountRemaining int
}

// Orbit implements the state for an orbit within a system
type Orbit struct {
	id     string // unique identifier for the orbit
	name   string // name of the orbit in the game; must be unique in the game
	owner  string // player who owns/controls this orbit in the system
	planet *Planet
}

// Planet occupies an orbit
type Planet struct {
	id string // unique identifier for the planet
	// not all planets have nations.
	// the game starts with a ten solar systems having a single planet that is populated.
	// no other planets have a population.
	nations      map[string]*Nation
	kind         string // gas giant, terrestrial, asteroids
	habitability int    // range from 0 to 25
	deposits     []*NaturalResource
}

// Player implements the player in the game
type Player struct {
	id     string // unique identifier for the player
	userId string // unique identifier for the user
	name   string // name of the player in the game; must be unique in the game
}

// SolarSystem implements a single solar system.
type SolarSystem struct {
	id   string // unique identifier for the system
	name string // name of the system in the game; must be unique in the game
	// coordinates are the location of the system within the cluster
	coordinates struct {
		// coordinates should be integers ranging from 0 to 30
		x, y, z float64
	}
	// orbits in the system. nil means the orbit is empty.
	orbits [10]*Orbit
}

// User represents the person playing.
type User struct {
	id string // unique identifier for the user
}

type AddPlayer struct {
	userId string
	name   string // name of the player in the game; must be unique in the game
}

func (o AddPlayer) Eval(g *Game) error {
	val, ok := g.objects[o.userId]
	if ok {
		if val.kind != "USER" {
			return fmt.Errorf("invalid user id")
		}
		return nil
	}
	player := &Player{userId: o.userId, name: o.name}
	if _, ok := g.objects[player.name]; ok {
		return fmt.Errorf("duplicate player name")
	}
	g.players[player.id] = player
	g.players[player.name] = player
	//g.objects[player.id] = player
	//g.objects[player.name] = player
	return nil
}

type NamePlayer struct {
	id   string
	name string
}

func (o NamePlayer) Eval(g *Game) error {
	player, ok := g.players[o.id]
	if !ok {
		// player is not in game
		return nil
	}
	player.name = o.name
	return nil
}

type RemovePlayer struct {
	id string
}

func (o RemovePlayer) Eval(g *Game) error {
	player, ok := g.players[o.id]
	if !ok {
		// player is not in game
		return nil
	}
	delete(g.players, player.id)
	delete(g.players, player.name)
	delete(g.objects, player.id)
	delete(g.objects, player.name)
	for _, obj := range g.objects {
		switch obj.kind {
		case "FACTORY":
			val := obj.factory
			if val.owner == player.id {
				// remove player as owner of the object
				val.owner = ""
			}
		case "ORBIT":
			val := obj.orbit
			if val.owner == player.id {
				// remove player as owner of the object
				val.owner = ""
			}
		case "PLAYER":
			// players don't own other players, so ignore this entry
		case "SYSTEM":
			// players don't own systems, so ignore this entry
		case "USER":
			// players don't own users, so ignore this entry
		default:
			panic(fmt.Sprintf("assert(object.type != %q)", obj.kind))
		}
	}
	return nil
}

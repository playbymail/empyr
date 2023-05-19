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

package games

import (
	"fmt"
	"github.com/playbymail/empyr/cmd/server/pkg/planets"
	"github.com/playbymail/empyr/cmd/server/pkg/prng"
	"github.com/playbymail/empyr/cmd/server/pkg/systems"
	"github.com/playbymail/empyr/cmd/server/pkg/tribes"
	"math/rand"
	"time"
)

// Game defines the properties of a game.
type Game struct {
	ID        string    `json:"game_id"`
	Created   time.Time `json:"created_at"`             // when the game was created
	Completed time.Time `json:"completed_at,omitempty"` // when the game was completed

	// player name maps to a tribe in the game
	Players map[string]*tribes.Tribe `json:"players,omitempty"`

	Systems []*systems.System `json:"systems,omitempty"`
	Ships   struct {
		id   map[string]*Ship
		name map[string]*Ship
	}

	prng prng.Generator
}

type GameConfig struct {
	PRNG            *rand.Rand
	NumberOfSystems int
	Systems         systems.SystemConfig
}

// Roll returns the result of rolling N D-sided dice.
func (g *Game) roll(n, d int) (total int) {
	for i := 0; i < n; i++ {
		total = total + g.prng.Intn(d) + 1
	}
	return total
}

type Ship struct {
	ID    string  `json:"ship_id"`
	Name  string  `json:"name"`
	Dodge float64 `json:"dodge_pct,omitempty"` // percentage of speed allocated to dodging weapons fire
}

// AddPlayer adds a player as ruler of the next available system.
func (g *Game) AddPlayer(playerName, tribeName string) (err error) {
	nextSystem := len(g.Players)
	if nextSystem > len(g.Systems) {
		return fmt.Errorf("no available systems")
	}

	// force the third orbit in the system to contain a home world
	homeWorldGenerator := planets.GenerateHomeworld()
	if g.Systems[nextSystem].Stars[0].Orbits[3].Planet, err = homeWorldGenerator(g.prng); err != nil {
		return err
	}

	g.Players[playerName] = &tribes.Tribe{
		Name:       tribeName,
		HomeSystem: g.Systems[nextSystem],
		HomeWorld:  g.Systems[nextSystem].Stars[0].Orbits[3].Planet,
	}
	return nil
}

// ShuffleSystem randomizes the order of the systems in place.
func (g *Game) ShuffleSystems(ts prng.Generator) {
	ts.Shuffle(len(g.Systems), func(i, j int) {
		g.Systems[i], g.Systems[j] = g.Systems[j], g.Systems[i]
	})
}

// SortSystems sorts systems by system name.
func (g *Game) SortSystems() {
	// sort the systems for later reporting
	for i := range g.Systems {
		for j := i + 1; j < len(g.Systems); j++ {
			if g.Systems[i].Name > g.Systems[j].Name {
				g.Systems[i], g.Systems[j] = g.Systems[j], g.Systems[i]
			}
		}
	}
}

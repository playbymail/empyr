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
	"github.com/playbymail/empyr/cmd/server/pkg/prng"
	"github.com/playbymail/empyr/cmd/server/pkg/systems"
	"github.com/playbymail/empyr/cmd/server/pkg/tribes"
	"time"
)

type Generator func(ts prng.Generator) (*Game, error)

// DefaultGenerator returns a generator with these rules:
//
//	Create 10 Systems
//
// Someone must add players, tribes, and home systems after the systems are generated.
func DefaultGenerator() Generator {
	generateSystem := systems.DefaultGenerator()
	numberOfSystems := 3
	return func(ts prng.Generator) (*Game, error) {
		var game Game
		game.Created = time.Now()
		game.Players = make(map[string]*tribes.Tribe)
		game.Systems = make([]*systems.System, numberOfSystems, numberOfSystems)
		game.prng = ts

		for i := 0; i < numberOfSystems; i++ {
			system, err := generateSystem(ts)
			if err != nil {
				return nil, err
			}
			game.Systems[i] = system
		}

		// add home worlds to 10 random systems. (since the systems were
		// created randomly, we can just pick the first 10 in the set.)
		for i := 0; i < 10 && i < len(game.Systems); i++ {
			// create home world
			// assign 25 nations to the home world
		}

		return &game, nil
	}
}

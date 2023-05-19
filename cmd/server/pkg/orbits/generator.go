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

package orbits

import (
	"github.com/google/uuid"
	"github.com/playbymail/empyr/cmd/server/pkg/planets"
	"github.com/playbymail/empyr/cmd/server/pkg/prng"
)

type Generator func(generator prng.Generator) (*Orbit, error)

// DefaultGenerator returns a generator with the following rules:
//
//	75% chance orbit contains a planet
func DefaultGenerator() Generator {
	generatePlanet := planets.DefaultGenerator()
	return func(ts prng.Generator) (*Orbit, error) {
		var o Orbit
		o.ID = uuid.New().String()
		if ts.Intn(4) == 0 {
			// we're creating an empty orbit to give ships and colonies something to park in.
		} else {
			// we're creating a planet or asteroid belt in this orbit
			p, err := generatePlanet(ts)
			if err != nil {
				return nil, err
			}
			o.Planet = p
		}
		return &o, nil
	}
}

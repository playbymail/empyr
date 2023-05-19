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

package stars

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/playbymail/empyr/cmd/server/pkg/orbits"
	"github.com/playbymail/empyr/cmd/server/pkg/prng"
)

type Generator func(ts prng.Generator) (*Star, error)

// DefaultGenerator returns a generator with the following rules:
//
//	11 orbits
//	Orbit[0] is treated as the "11th Orbit" in the rulebook.
func DefaultGenerator() Generator {
	generateOrbit := orbits.DefaultGenerator()
	return func(ts prng.Generator) (*Star, error) {
		var s Star
		s.ID = uuid.New().String()
		s.Name = fmt.Sprintf("%02d-%02d-%02d", 0, 0, 0)
		for i := 1; i <= 10; i++ {
			orbit, err := generateOrbit(ts)
			if err != nil {
				return nil, err
			}
			s.Orbits[i] = orbit
		}
		return &s, nil
	}
}

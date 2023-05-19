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

package systems

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/playbymail/empyr/cmd/server/pkg/prng"
	"github.com/playbymail/empyr/cmd/server/pkg/stars"
)

type Generator func(ts prng.Generator) (*System, error)

// DefaultGenerator returns a generator with the following rules:
//
//	from 1 to 10 orbits
//	each orbit after the first has a 95% chance of containing a planet
//	if an orbit is empty, all remaining orbits are also empty
func DefaultGenerator() Generator {
	generateStar := stars.DefaultGenerator()
	return func(ts prng.Generator) (*System, error) {
		var s System
		s.ID = uuid.New().String()

		minXYZ, maxXYZ := 0, 30
		s.X = minXYZ + ts.Intn(maxXYZ-minXYZ)
		s.Y = minXYZ + ts.Intn(maxXYZ-minXYZ)
		s.Z = minXYZ + ts.Intn(maxXYZ-minXYZ)

		s.Name = fmt.Sprintf("%02d-%02d-%02d", s.X, s.Y, s.Z)
		star, err := generateStar(ts)
		if err != nil {
			return nil, err
		}
		star.Name = s.Name
		s.Stars = append(s.Stars, star)
		return &s, nil
	}
}

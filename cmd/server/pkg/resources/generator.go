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

package resources

import (
	"github.com/google/uuid"
	"github.com/playbymail/empyr/cmd/server/pkg/prng"
)

type Generator func(prng.Generator) (*Resource, error)

// DefaultGenerator returns a generator with the following distribution:
//
//	 Fuel: 33%, Yield 10..90%, Units: 10..99 million
//	 Gold:  1%, Yield  1.. 9%, Units: 10..99 million
//	Metal: 32%, Yield 10..90%, Units: 10..99 million
//
// NonMetal: 34%, Yield 10..90%, Units: 10..99 million
func DefaultGenerator() Generator {
	return func(ts prng.Generator) (*Resource, error) {
		var nr Resource
		nr.ID = uuid.New().String()
		switch roll := ts.Intn(100); {
		case roll < 33:
			nr.Type = FUEL
			nr.YieldPct = float64(10+ts.Intn(80)) / 100
			nr.InitialAmount = 1000000 * int64(1+ts.Intn(99))
		case roll < 34:
			nr.Type = GOLD
			nr.YieldPct = float64(1+ts.Intn(9)) / 100
			nr.InitialAmount = 1000000 * int64(1+ts.Intn(99))
		case roll < 66:
			nr.Type = METAL
			nr.YieldPct = float64(10+ts.Intn(80)) / 100
			nr.InitialAmount = 1_000_000 * int64(1+ts.Intn(99))
		default:
			nr.Type = NONMETAL
			nr.YieldPct = float64(10+ts.Intn(80)) / 100
			nr.InitialAmount = 1_000_000 * int64(1+ts.Intn(99))
		}
		nr.AmountRemaining = nr.InitialAmount
		return &nr, nil
	}
}

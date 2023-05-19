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

package orders

import (
	"errors"
	"github.com/playbymail/empyr/pkg/empyr"
)

// A JUMP command moves a ship to a new location.
type JUMP struct {
	Ship int            // ID of the ship that will jump
	To   empyr.Location // location the ship will jump to
}

// DoJump moves a ship to a new location.
// Actually, it doesn't move it?
func DoJump(ship empyr.Ship, to empyr.Location) (empyr.Location, float64, error) {
	// this is just to test the constant errors package.
	l, f, err := ship.Jump(to)
	if err != nil {
		switch {
		case errors.Is(err, empyr.ErrDistanceExceedsCapacity): // whoot
		case errors.Is(err, empyr.ErrInsufficientFuel): // whoot
		case errors.Is(err, empyr.ErrMassExceedsCapacity): // whoot
		}
	}
	return l, f, err
}

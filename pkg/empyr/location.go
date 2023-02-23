// empyr - a game engine for Empyrean Challenge
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

package empyr

import "math"

// Location for a ship, colony, planet, or system.
type Location struct {
	X, Y, Z int
}

func (l Location) dXYZ(l2 Location) (int, int, int) {
	return l.X - l2.X, l.Y - l2.Y, l.Z - l2.Z
}

func (l Location) DistanceFrom(l2 Location) float64 {
	return math.Sqrt(float64(l.DistanceSquared(l2)))
}

func (l Location) DistanceSquared(l2 Location) int {
	dx, dy, dz := l.dXYZ(l2)
	return dx*dx + dy*dy + dz*dz
}

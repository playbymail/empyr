// Copyright (c) 2025 Michael D Henderson. All rights reserved.

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

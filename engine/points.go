// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import "math"

// DistanceBetween calculates the Euclidean distance between two points in 3D space.
func DistanceBetween(a, b Point_t) float64 {
	dx, dy, dz := a.X-b.X, a.Y-b.Y, a.Z-b.Z
	return math.Sqrt(float64(dx*dx + dy*dy + dz*dz))
}

// tooClose returns true if the distance between a point and any
// point in the slice of points is less than 5 units.
func tooClose(a Point_t, points []Point_t) bool {
	const minDistance = 5.0
	for _, p := range points {
		if DistanceBetween(a, p) < minDistance {
			return true
		}
	}
	return false
}

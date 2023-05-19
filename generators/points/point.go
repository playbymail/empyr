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

package points

import (
	"math"
	"math/rand"
)

type Point struct { // location being set up
	X, Y, Z   float64
	Neighbors struct {
		avd float64     // average distance to neighbors
		nb  []*neighbor // neighbors sorted by distance
	}
	fromOrigin float64
}

type neighbor struct {
	point    *Point
	distance float64
}

func (p *Point) AvgDistance() float64 {
	return p.Neighbors.avd
}

func (p *Point) DistanceTo(b *Point) float64 {
	dx, dy, dz := p.X-b.X, p.Y-b.Y, p.Z-b.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (p *Point) Scale(n float64) *Point {
	return &Point{X: p.X * n, Y: p.Y * n, Z: p.Z * n}
}

func ClusteredPoint() *Point {
	var u = rand.Float64()
	var v = rand.Float64()
	var theta = u * 2.0 * math.Pi
	var phi = math.Acos(2.0*v - 1.0)
	var sinTheta = math.Sin(theta)
	var cosTheta = math.Cos(theta)
	var sinPhi = math.Sin(phi)
	var cosPhi = math.Cos(phi)
	var r = rand.Float64()
	return &Point{
		X: r * sinPhi * cosTheta,
		Y: r * sinPhi * sinTheta,
		Z: r * cosPhi,
	}
}

func SpherePoint() *Point {
	var u = rand.Float64()
	var v = rand.Float64()
	var theta = u * 2.0 * math.Pi
	var phi = math.Acos(2.0*v - 1.0)
	var sinTheta = math.Sin(theta)
	var cosTheta = math.Cos(theta)
	var sinPhi = math.Sin(phi)
	var cosPhi = math.Cos(phi)
	return &Point{
		X: sinPhi * cosTheta,
		Y: sinPhi * sinTheta,
		Z: cosPhi,
	}
}

func UniformPoint() *Point {
	var u = rand.Float64()
	var v = rand.Float64()
	var theta = u * 2.0 * math.Pi
	var phi = math.Acos(2.0*v - 1.0)
	var sinTheta = math.Sin(theta)
	var cosTheta = math.Cos(theta)
	var sinPhi = math.Sin(phi)
	var cosPhi = math.Cos(phi)
	var r = math.Cbrt(rand.Float64())
	return &Point{
		X: r * sinPhi * cosTheta,
		Y: r * sinPhi * sinTheta,
		Z: r * cosPhi,
	}
}

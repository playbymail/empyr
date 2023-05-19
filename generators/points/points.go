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
	"sort"
)

type Points struct {
	Points []*Point
}

func NewPoints(n int, pgen func() *Point) *Points {
	p := &Points{Points: make([]*Point, n, n)}
	for i := range p.Points {
		p.Points[i] = pgen()
	}
	p.SetNeighbors(0)
	return p
}

func (p *Points) clone() *Points {
	cp := &Points{Points: make([]*Point, p.Length(), p.Length())}
	for i, point := range p.Points {
		cp.Points[i] = &Point{
			X: point.X,
			Y: point.Y,
			Z: point.Z,
		}
	}
	return cp
}

// Length returns the number of points.
func (p *Points) Length() int {
	return len(p.Points)
}

// MinAvgMax returns the minimum, average, and maximum average distances between neighbors
func (p *Points) MinAvgMax() (min, avg, max float64) {
	min, max = math.MaxFloat64, -1.0
	for _, point := range p.Points {
		if point.Neighbors.nb[0].distance < min {
			min = point.Neighbors.nb[0].distance
		}
		avg += point.Neighbors.avd
		if point.Neighbors.nb[0].distance > max {
			max = point.Neighbors.nb[0].distance
		}
	}
	avg = avg / float64(len(p.Points))
	return min, avg, max
}

// CullByCompanions culls out the systems that are closest to each other.
func (p *Points) CullByCompanions(n int) *Points {
	cp := p.clone()
	cp.SetNeighbors(n)
	sort.Slice(cp.Points, func(i, j int) bool {
		return cp.Points[i].Neighbors.avd < cp.Points[j].Neighbors.avd
	})
	cp.Points = cp.Points[1:]
	if len(p.Points) == len(cp.Points) {
		panic("!")
	}
	return cp
}

// CullByDistanceFromOrigin culls out the systems that are farthest from the origin
func (p *Points) CullByDistanceFromOrigin() *Points {
	cp := p.clone()
	cp.SortByDistanceOrigin()
	cp.Points = cp.Points[:125]
	return cp
}

// CullByMinDistance culls out the systems that are closest to each other
func (p *Points) CullByMinDistance(min float64) *Points {
	cp := &Points{}
	for _, point := range p.Points {
		if point.Neighbors.nb[0].distance < min {
			continue
		}
		cp.Points = append(cp.Points, &Point{
			X: point.X,
			Y: point.Y,
			Z: point.Z,
		})
	}
	cp.SetNeighbors(0)
	return cp
}

func (p *Points) SetNeighbors(n int) {
	for _, origin := range p.Points {
		origin.Neighbors.nb = make([]*neighbor, len(p.Points), len(p.Points))
		for i, point := range p.Points {
			distance := origin.DistanceTo(point)
			origin.Neighbors.nb[i] = &neighbor{
				point:    point,
				distance: distance,
			}
		}
		sort.Slice(origin.Neighbors.nb, func(i, j int) bool {
			return origin.Neighbors.nb[i].distance < origin.Neighbors.nb[j].distance
		})
		// the first entry should be the origin itself
		origin.Neighbors.nb = origin.Neighbors.nb[1:]
		// if n is set, keep only that many neighbors
		if n > 0 && n < len(origin.Neighbors.nb) {
			origin.Neighbors.nb = origin.Neighbors.nb[:n]
		}
		// calculate the average distance of remaining neighbors
		var distance float64
		for _, nb := range origin.Neighbors.nb {
			distance += nb.distance
		}
		origin.Neighbors.avd = distance / float64(len(origin.Neighbors.nb))
	}
}

// SortByDistanceDesc sorts the points by average distance from neighbors, descending.
// N is the number of neighbors to use for distance.
func (p *Points) SortByDistanceDesc(n int) *Points {
	p.SetNeighbors(0)
	sort.Slice(p.Points, func(i, j int) bool {
		return p.Points[i].Neighbors.avd > p.Points[j].Neighbors.avd
	})
	return p
}

// SortByDistanceOrigin sorts the points by distance from origin, ascending.
func (p *Points) SortByDistanceOrigin() *Points {
	var origin Point
	for _, point := range p.Points {
		point.fromOrigin = origin.DistanceTo(point)
	}
	sort.Slice(p.Points, func(i, j int) bool {
		return p.Points[i].fromOrigin < p.Points[j].fromOrigin
	})
	return p
}

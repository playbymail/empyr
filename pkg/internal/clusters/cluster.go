// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package clusters implements the cluster generation algorithm.
// It returns a list of systems containing their coordinates
// and the number of stars in each system. The location of the
// 3- and 4-star systems are fixed; the rest are randomly located.
package clusters

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"math"
	"math/rand/v2"
)

func GenerateCluster() (systems []System) {
	// current list of points generated; used to determine distance to neighbors
	var points []Point

	// the location of the 4-star system is fixed
	system := System{Coordinates: Point{X: 0, Y: 0, Z: 0}, NumberOfStars: 4}
	points = append(points, system.Coordinates)
	systems = append(systems, system)

	// the location of the 3-star systems are fixed
	for _, point := range []Point{
		Point{X: 1, Y: 1, Z: 1},
		Point{X: -1, Y: -1, Z: -1},
		Point{X: -2, Y: 2, Z: 0},
	} {
		system = System{Coordinates: point, NumberOfStars: 3}
		points = append(points, system.Coordinates)
		systems = append(systems, system)
	}

	// the 2-star systems are always within 7 units of the origin
	for n := 1; n <= 9; n++ {
		point := randomPoint(7)
		for maxAttempts := 0; point.TooClose(points) && maxAttempts < 1_000; maxAttempts++ {
			point = randomPoint(7)
		}
		systems = append(systems, System{Coordinates: point, NumberOfStars: 2})
	}

	// the remaining single-star systems are always within the 31-unit cube centered at the origin
	for len(systems) < 100 {
		point := randomCubePoint()
		for maxAttempts := 0; point.TooClose(points) && maxAttempts < 1_000; maxAttempts++ {
			point = randomCubePoint()
		}
		systems = append(systems, System{Coordinates: point, NumberOfStars: 1})
	}

	return systems
}

type System struct {
	Coordinates   Point
	NumberOfStars int
}

type Point struct {
	X, Y, Z int
}

// DistanceBetween calculates the Euclidean distance between two points in 3D space.
func (a Point) DistanceBetween(b Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return math.Sqrt(float64(dx*dx + dy*dy + dz*dz))
}

// Scale scales the point by a factor of scale, which should be the radius of the cluster.
func (a Point) Scale(scale float64) Point {
	return Point{
		X: int(math.Round(float64(a.X) * scale)),
		Y: int(math.Round(float64(a.Y) * scale)),
		Z: int(math.Round(float64(a.Z) * scale)),
	}
}

// String implements the fmt.Stringer interface.
func (a Point) String() string {
	return fmt.Sprintf("(%3d %3d %3d)", a.X, a.Y, a.Z)
}

// TooClose returns true if the distance between a point and any
// point in the slice of points is less than 5 units.
func (a Point) TooClose(points []Point) bool {
	const minDistance = 5.0
	for _, p := range points {
		if a.DistanceBetween(p) < minDistance {
			return true
		}
	}
	return false
}

// randomPoint returns scaled coordinates with a uniform volume distribution
// within a sphere of the given radius.
func randomPoint(radius float64) Point {
	// generate a random distance to ensure uniform distribution within the sphere
	d := math.Cbrt(rand.Float64()) // Cube root to ensure uniform distribution
	d = radius * d                 // Scale to the desired radius

	// generate random angles for spherical coordinates
	theta := rand.Float64() * 2 * math.Pi  // 0 to 2π
	phi := math.Acos(2*rand.Float64() - 1) // 0 to π

	// convert spherical coordinates to Cartesian coordinates
	return Point{
		X: int(math.Round(d * math.Sin(phi) * math.Cos(theta))),
		Y: int(math.Round(d * math.Sin(phi) * math.Sin(theta))),
		Z: int(math.Round(d * math.Cos(phi))),
	}
}

// randomCubePoint returns coordinates within a cube with side length of 31,
// centered at the origin and uniformly distributed.
func randomCubePoint() Point {
	return Point{
		X: rand.IntN(31) - 15, // -15 to 15
		Y: rand.IntN(31) - 15, // -15 to 15
		Z: rand.IntN(31) - 15, // -15 to 15
	}
}

var (
	//go:embed templates
	templateFS embed.FS
)

func GenerateClusterHTML(input []System) (*bytes.Buffer, error) {
	buffer := &bytes.Buffer{}

	ts, err := template.ParseFS(templateFS, "templates/cluster-map.gohtml")
	if err != nil {
		return nil, err
	}

	type ClusterSystem struct {
		Id      int
		X, Y, Z int
		Size    float64     // for rendering cluster
		Color   template.JS // for rendering cluster
		Warps   []Point
	}
	var systems []ClusterSystem
	for i := range input {
		system := ClusterSystem{
			Id: i + 1,
			X:  input[i].Coordinates.X,
			Y:  input[i].Coordinates.Y,
			Z:  input[i].Coordinates.Z,
		}
		switch input[i].NumberOfStars {
		case 1:
			system.Color = "Gray"
			system.Size = 0.33
		case 2:
			system.Color = "White"
			system.Size = 0.45
		case 3:
			system.Color = "Teal"
			system.Size = 0.5
		case 4:
			system.Color = "Yellow"
			system.Size = 0.75
		default:
			system.Color = "Random"
			system.Size = 0.3
		}
		systems = append(systems, system)
	}

	if err = ts.Execute(buffer, systems); err != nil {
		return nil, err
	}
	return buffer, nil
}

// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package empyr

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/playbymail/empyr/pkg/internal/clusters"
	"html/template"
	"log"
	"math"
	"math/rand/v2"
	"os"
)

// Game implements the entire game state.
type Game struct {
	Id      int
	Code    string
	Name    string
	Turn    int
	Systems []System
}

func NewGame(code string, name string) (Game, error) {
	systems, err := generateSystems()
	if err != nil {
		return Game{}, err
	}
	for _, system := range systems {
		fmt.Printf("System %3d at (%3d, %3d, %3d) has %d stars\n", system.Id, system.X, system.Y, system.Z, len(system.Stars))
	}
	cc := clusters.GenerateCluster()
	log.Printf("cc.gch: len(cc) %d\n", len(cc))
	if buf, err := clusters.GenerateClusterHTML(cc); err != nil {
		log.Printf("cc.gch: %v\n", err)
	} else if err = os.WriteFile("cluster-map.html", buf.Bytes(), 0644); err != nil {
		log.Printf("cc.gch: write: %v\n", err)
	} else {
		log.Printf("cc.ghc: created cluster-map.html\n")
	}
	return Game{
		Code:    code,
		Name:    name,
		Systems: systems,
	}, nil
}

// ReadGame loads a game's data file.
func ReadGame(filename string) (Game, error) {
	var g Game
	if data, err := os.ReadFile(filename); err != nil {
		return g, fmt.Errorf("game: open: %w", err)
	} else if err = json.Unmarshal(data, &g); err != nil {
		return g, fmt.Errorf("game: parse: %w", err)
	}
	return g, nil
}

type Point struct {
	X, Y, Z float64
}

// randomPoint generates a random point in 3D space, with x, y, and z coordinates rangins from -1.0 to 1.0.
func randomPoint(radius, minDistToNeighbor float64, points []Point) (Point, error) {
	origin := Point{X: 0, Y: 0, Z: 0}
	for maxAttempts := 1_000; maxAttempts > 0; maxAttempts-- {
		point := genXYZ(radius)
		if distance(point, origin) < 0.001 {
			continue
		}
		// the system must not be too close to any other system
		if tooClose(point, points, minDistToNeighbor) {
			continue
		}
		return point, nil
	}
	return Point{}, fmt.Errorf("could not generate a random point")
}

func tooClose(point Point, points []Point, minDistToNeighbor float64) bool {
	for _, p := range points {
		if distance(point, p) < minDistToNeighbor {
			return true
		}
	}
	return false
}

// genXYZ returns un-scaled coordinates with a uniform distribution within a 1 unit sphere
func genXYZ(radius float64) Point {
	// generate a random distance to ensure uniform distribution within the sphere
	d := math.Cbrt(rand.Float64()) // Cube root to ensure uniform distribution
	d = radius * d                 // Scale to the desired radius

	// generate random angles for spherical coordinates
	theta := rand.Float64() * 2 * math.Pi  // 0 to 2π
	phi := math.Acos(2*rand.Float64() - 1) // 0 to π

	// convert spherical coordinates to Cartesian coordinates
	return Point{
		X: math.Round(d * math.Sin(phi) * math.Cos(theta)),
		Y: math.Round(d * math.Sin(phi) * math.Sin(theta)),
		Z: math.Round(d * math.Cos(phi)),
	}
}

// distance calculates the Euclidean distance between two points in 3D space.
func distance(a, b Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

type System struct {
	Id      int
	X, Y, Z int
	Stars   []Star
	Warps   []any
	Size    float64     // for rendering cluster
	Color   template.JS // for rendering cluster
}

type Star struct {
	Id int
}

// generateSystems generates a slice of systems, each containing a number of stars
// according to a predefined distribution. Systems are positioned in 3D space
// with constraints to ensure they are not placed too close to one another or too far
// from the origin based on their star count:
// - Systems with 4 stars are within 2 units of the origin.
// - Systems with 3 stars are within 4 units of the origin.
// - Systems with 2 stars are within 10 units of the origin.
// - Systems with 1 star can be anywhere within 16 units of the origin.
//
// Each system is placed at a randomly assigned position that meets these constraints,
// with a maximum number of attempts (1,000) to ensure placement feasibility.
// If a valid placement cannot be found, the function returns an error.
func generateSystems() (systems []System, err error) {
	distribution := []struct {
		starsPerSystem    int
		numSystems        int
		minDistToNeighbor float64
		radius            float64
	}{
		{4, 1, 0.0, 0.0},  //  1 system  with 4 stars ==  4 stars
		{3, 3, 0.1, 0.2},  //  3 systems with 3 stars ==  9 stars
		{2, 9, 0.2, 0.75}, //  9 systems with 2 stars == 18 stars
		{1, 69, 0.1, 1.0}, // 69 systems with 1 star  == 69 stars
	}

	var points []Point
	for _, entry := range distribution {
		for i := 0; i < entry.numSystems; i++ {
			if len(points) == 0 {
				points = append(points, Point{})
				continue
			}
			var point Point
			switch len(points) {
			case 0:
				point = Point{} // 4-star system is always at the origin
			case 1:
				point = Point{X: 1, Y: 1, Z: 1} // first 3-star system is at (1, 1, 1)
			case 2:
				point = Point{X: -1, Y: -1, Z: -1} // second 3-star system is at (-1, -1, -1)
			case 3:
				point = Point{X: -2, Y: 2, Z: 0} // third 3-star system is at (-2, 2, 0)
			case 4, 5, 6, 7, 8, 9, 10, 11, 12:
				// 2-star systems are always within 7 units of the origin
				point, err = randomPoint(7, 1.9, points)
				if err != nil {
					return nil, err
				}
			default:
				point, err = randomPoint(15, 2.9, points)
				if err != nil {
					return nil, err
				}
			}
			points = append(points, point)
		}
	}

	systemID, starID := 1, 1
	for _, entry := range distribution {
		for i := 0; i < entry.numSystems; i++ {
			// create a new system with the given number of stars
			system := System{
				Id: systemID,
				X:  int(points[0].X),
				Y:  int(points[0].Y),
				Z:  int(points[0].Z),
			}

			systemID++
			points = points[1:]

			for j := 0; j < entry.starsPerSystem; j++ {
				system.Stars = append(system.Stars, Star{Id: starID})
				starID++
			}

			systems = append(systems, system)
		}
	}

	const templatesPath = "../templates/cluster.gohtml"

	//// move the origin from 0,0,0 to 16,16,16
	//for i := range systems {
	//	systems[i].X += 17
	//	systems[i].Y += 17
	//	systems[i].Z += 17
	//}

	return systems, nil
}

var (
	//go:embed templates
	templatesFS embed.FS
)

func GenerateClusterHTML(path string, systems []System) error {
	ts, err := template.ParseFS(templatesFS, "templates/cluster.gohtml")
	if err != nil {
		return err
	}
	w, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	for i := range systems {
		switch len(systems[i].Stars) {
		case 1:
			systems[i].Color = "Gray"
			systems[i].Size = 0.33
		case 2:
			systems[i].Color = "White"
			systems[i].Size = 0.45
		case 3:
			systems[i].Color = "Teal"
			systems[i].Size = 0.5
		case 4:
			systems[i].Color = "Yellow"
			systems[i].Size = 0.75
		default:
			systems[i].Color = "Random"
			systems[i].Size = 0.3
		}
	}
	defer func(fp *os.File) {
		_ = fp.Close()
	}(w)
	return ts.Execute(w, systems)
}

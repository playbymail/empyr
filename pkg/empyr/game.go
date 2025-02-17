// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package empyr

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/playbymail/empyr/models/coordinates"
	"github.com/playbymail/empyr/pkg/internal/clusters"
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
	// generate the systems within the cluster
	var systems []System
	for no, cc := range clusters.GenerateCluster() {
		system := System{
			Id: no + 1,
			Coordinates: coordinates.Coordinates{
				X: cc.Coordinates.X + 15,
				Y: cc.Coordinates.Y + 15,
				Z: cc.Coordinates.Z + 15,
			},
		}
		for i := 1; i <= cc.NumberOfStars; i++ {
			star := Star{Id: i}
			system.Stars = append(system.Stars, star)
		}
		systems = append(systems, system)
	}

	//for _, system := range systems {
	//	fmt.Printf("System %3d at (%3d%3d%3d) has %d stars\n", system.Id, system.Coordinates.X-15, system.Coordinates.Y-15, system.Coordinates.Z-15, len(system.Stars))
	//}

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

type System struct {
	Id          int
	Coordinates coordinates.Coordinates
	Stars       []Star
}

type Star struct {
	Id int
}

var (
	//go:embed templates
	templatesFS embed.FS
)

func (g *Game) ClusterHTML() (*bytes.Buffer, error) {
	// adapt systems back to clusters.System
	var cc []clusters.System
	for _, system := range g.Systems {
		cc = append(cc, clusters.System{
			Coordinates: clusters.Point{
				X: system.Coordinates.X,
				Y: system.Coordinates.Y,
				Z: system.Coordinates.Z,
			},
			NumberOfStars: len(system.Stars),
		})
	}
	return clusters.GenerateClusterHTML(cc)
}

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

package cli

import (
	"encoding/json"
	"fmt"
	"github.com/playbymail/empyr/generators/clusters"
	"github.com/playbymail/empyr/models/games"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

// cmdGenerateCluster runs the cluster generator command
var cmdGenerateCluster = &cobra.Command{
	Use:   "cluster",
	Short: "generate a new cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		if argsGenerate.path != "" {
			argsGenerate.path = filepath.Clean(argsGenerate.path)
			if sb, err := os.Stat(argsGenerate.path); err != nil {
				return err
			} else if !sb.IsDir() {
				return fmt.Errorf("path must be a valid directory")
			}
		}
		gamePath := filepath.Join(argsGenerate.path, argsGenerateCluster.game)
		if sb, err := os.Stat(gamePath); err != nil {
			log.Fatalf("invalid game path %s\n", gamePath)
		} else if !sb.IsDir() {
			log.Fatalf("game path %s is not a directory\n", gamePath)
		}
		log.Printf("generate: cluster: loading %s\n", gamePath)
		var g games.Game
		buffer, err := os.ReadFile(filepath.Join(gamePath, "game.json"))
		if err != nil {
			log.Fatal(err)
		} else if err = json.Unmarshal(buffer, &g); err != nil {
			log.Fatal(err)
		}
		log.Printf("generate: cluster: loaded  %s\n", filepath.Join(gamePath, "game.json"))

		optCluster := []clusters.Option{}
		if opt, err := clusters.CreateHtmlMap(filepath.Join(gamePath, argsGenerateCluster.mapFile)); err != nil {
			log.Fatal(err)
		} else {
			optCluster = append(optCluster, opt)
		}

		if argsGenerateCluster.kind.cluster {
			if opt, err := clusters.SetKind("clustered"); err != nil {
				log.Fatal(err)
			} else {
				optCluster = append(optCluster, opt)
			}
		} else if argsGenerateCluster.kind.surface {
			if opt, err := clusters.SetKind("surface"); err != nil {
				log.Fatal(err)
			} else {
				optCluster = append(optCluster, opt)
			}
		} else {
			if opt, err := clusters.SetKind("uniform"); err != nil {
				log.Fatal(err)
			} else {
				optCluster = append(optCluster, opt)
			}
		}

		if opt, err := clusters.SetSystems(128); err != nil {
			log.Fatal(err)
		} else {
			optCluster = append(optCluster, opt)
		}

		if opt, err := clusters.SetRadius(argsGenerateCluster.radius); err != nil {
			log.Fatal(err)
		} else {
			optCluster = append(optCluster, opt)
		}

		c, sy, st, err := clusters.Generate(optCluster...)
		if err != nil {
			log.Fatal(err)
		}
		// adapt c to json
		file := filepath.Join(gamePath, "cluster.json")
		if data, err := json.MarshalIndent(c, "", "  "); err != nil {
			log.Fatal(err)
		} else if err = os.WriteFile(file, data, 0660); err != nil {
			log.Fatal(err)
		}
		log.Printf("cluster: created %s\n", file)
		// adapt sy to json
		file = filepath.Join(gamePath, "systems.json")
		if data, err := json.MarshalIndent(sy, "", "  "); err != nil {
			log.Fatal(err)
		} else if err = os.WriteFile(file, data, 0660); err != nil {
			log.Fatal(err)
		}
		log.Printf("cluster: created %s\n", file)
		// adapt st to json
		file = filepath.Join(gamePath, "stars.json")
		if data, err := json.MarshalIndent(st, "", "  "); err != nil {
			log.Fatal(err)
		} else if err = os.WriteFile(file, data, 0660); err != nil {
			log.Fatal(err)
		}
		log.Printf("cluster: created %s\n", file)

		return nil
	},
}

var argsGenerateCluster struct {
	game string // code for game to generate into
	kind struct {
		uniform, cluster, surface bool
	}
	mapFile string
	radius  float64
}

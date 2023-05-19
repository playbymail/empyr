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

package cli

import (
	"encoding/json"
	"github.com/playbymail/empyr/generators/clusters"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// cmdGenerateCluster runs the cluster generator command
var cmdGenerateCluster = &cobra.Command{
	Use:   "cluster",
	Short: "generate a new cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		optCluster := []clusters.Option{}
		if opt, err := clusters.CreateHtmlMap(argsGenerateCluster.mapFile); err != nil {
			log.Fatal(err)
		} else {
			optCluster = append(optCluster, opt)
		}
		if opt, err := clusters.SetKind(argsGenerateCluster.kind); err != nil {
			log.Fatal(err)
		} else {
			optCluster = append(optCluster, opt)
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
		if data, err := json.MarshalIndent(c, "", "  "); err != nil {
			log.Fatal(err)
		} else if err = os.WriteFile("g1/out/cluster.json", data, 0660); err != nil {
			log.Fatal(err)
		}
		log.Printf("cluster: created g1/out/cluster.json")
		// adapt sy to json
		if data, err := json.MarshalIndent(sy, "", "  "); err != nil {
			log.Fatal(err)
		} else if err = os.WriteFile("g1/out/systems.json", data, 0660); err != nil {
			log.Fatal(err)
		}
		log.Printf("cluster: created g1/out/stars.json")
		// adapt st to json
		if data, err := json.MarshalIndent(st, "", "  "); err != nil {
			log.Fatal(err)
		} else if err = os.WriteFile("g1/out/stars.json", data, 0660); err != nil {
			log.Fatal(err)
		}
		log.Printf("cluster: created g1/out/stars.json")

		return nil
	},
}

var argsGenerateCluster struct {
	kind    string // uniform, cluster, surface
	mapFile string
	radius  float64
}

func init() {
	cmdGenerate.AddCommand(cmdGenerateCluster)

	// inputs
	cmdGenerateCluster.Flags().StringVar(&argsGenerateCluster.kind, "kind", "uniform", "point distribution (uniform, clustered, sphere)")
	cmdGenerateCluster.Flags().StringVar(&argsGenerateCluster.mapFile, "html-map", "", "name of map file to create (optional)")
	cmdGenerateCluster.Flags().Float64Var(&argsGenerateCluster.radius, "radius", 15.0, "cluster radius")

	// outputs
}

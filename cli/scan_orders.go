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
	"github.com/playbymail/empyr/adapters"
	"github.com/playbymail/empyr/ec"
	"github.com/playbymail/empyr/parsers/orders"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var cmdScanOrders = &cobra.Command{
	Use:   "orders",
	Short: "Scan an orders file",
	Long:  `Load an orders file, scan it, and report on all errors.`,
	Run: func(cmd *cobra.Command, args []string) {
		argsScanOrders.ordersPath = filepath.Clean(argsScanOrders.ordersPath)
		log.Printf("scanning %q\n", argsScanOrders.ordersPath)

		e, err := ec.LoadGame(argsScanOrders.ordersPath)
		if err != nil {
			log.Fatal(err)
		}

		// load all the files
		files, err := filepath.Glob(filepath.Join(argsScanOrders.ordersPath, "orders.*.txt"))
		if err != nil {
			log.Fatal(err)
		}
		for _, name := range files {
			log.Printf("scanning %q\n", name)
			input, err := os.ReadFile(name)
			if err != nil {
				log.Fatal(err)
			}
			lexemes, err := orders.Scan(input)
			if err != nil {
				log.Fatal(err)
			}
			ods := orders.Parse(lexemes)
			err = e.AddOrders(adapters.OrdersToEngineOrders(ods))
			if err != nil {
				log.Printf("%s: %v\n", name, err)
			}
		}

		err = e.Process()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var argsScanOrders struct {
	ordersPath string
}

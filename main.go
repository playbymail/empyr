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

package main

import (
	"fmt"
	"github.com/playbymail/empyr/adapters"
	"github.com/playbymail/empyr/cli"
	"github.com/playbymail/empyr/ec"
	"github.com/playbymail/empyr/parsers/orders"
	"log"
	"os"
	"time"
)

func main() {
	started := time.Now()

	log.SetFlags(log.LstdFlags | log.LUTC)

	if err := dotfiles("EMPYR"); err != nil {
		log.Fatalf("main: %+v\n", err)
	}

	rv := 0
	if err := cli.Execute(); err != nil {
		log.Printf("\n%+v\n", err)
		rv = 2
	}

	log.Printf("\n")
	log.Printf("completed in %v\n", time.Now().Sub(started))

	os.Exit(rv)
}

func run(path string, debug bool) error {
	e, err := ec.LoadGame(path)
	if err != nil {
		return err
	}

	// load all the files
	for _, name := range []string{"orders.txt"} {
		input, err := os.ReadFile(name)
		if err != nil {
			return err
		}
		lexemes, err := orders.Scan(input)
		if err != nil {
			return err
		}
		ods := orders.Parse(lexemes)
		if debug {
			for _, od := range ods {
				fmt.Println(od)
			}
		}
		err = e.AddOrders(adapters.OrdersToEngineOrders(ods))
		if err != nil {
			log.Printf("%s: %v\n", name, err)
		}
	}

	err = e.Process()
	if err != nil {
		return err
	}

	return e.SaveGame(path)
}

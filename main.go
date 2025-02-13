// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"github.com/playbymail/empyr/adapters"
	"github.com/playbymail/empyr/cmd"
	"github.com/playbymail/empyr/ec"
	"github.com/playbymail/empyr/parsers/orders"
	"github.com/playbymail/empyr/pkg/dotenv"
	"log"
	"os"
	"time"
)

func main() {
	started := time.Now()

	log.SetFlags(log.LstdFlags | log.LUTC)

	if err := dotenv.Load("EMPYR"); err != nil {
		log.Fatalf("main: %+v\n", err)
	}

	rv := 0
	if err := cmd.Execute(); err != nil {
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

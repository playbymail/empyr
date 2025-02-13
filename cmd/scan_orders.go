// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cmd

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
	Short: "Scan orders file",
	Long:  `Load all orders file, scan them, and report on all errors.`,
	Run: func(cmd *cobra.Command, args []string) {
		argsScanOrders.ordersPath = filepath.Clean(argsScanOrders.ordersPath)
		log.Printf("scanning %q\n", argsScanOrders.ordersPath)

		e, err := ec.LoadGame(argsScanOrders.ordersPath)
		if err != nil {
			log.Fatal(err)
		}

		// find all orders files
		files, err := filepath.Glob(filepath.Join(argsScanOrders.ordersPath, "orders.*.txt"))
		if err != nil {
			log.Fatal(err)
		}

		// scan the order files
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
